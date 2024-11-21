package handler

import (
	"container/list"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (Result, bool)
	Set(key string, result Result)
}

type cacheEntry[T any] struct {
	Key         string
	Data        T
	CreatedAt   time.Time
	TTL         time.Duration
	LastAccess  time.Time
	AccessCount int
}

type BaseCache[T any] struct {
	mutex      sync.Mutex
	lruList    *list.List
	maxEntries int
}

func NewBaseCache[T any](maxEntries int) BaseCache[T] {
	return BaseCache[T]{
		lruList:    list.New(),
		maxEntries: maxEntries,
	}
}

func (c *BaseCache[T]) addToLRU(key string, entry cacheEntry[T]) {
	for e := c.lruList.Front(); e != nil; e = e.Next() {
		if e.Value.(cacheEntry[T]).Key == key {
			c.lruList.Remove(e)
			break
		}
	}
	c.lruList.PushFront(entry)
}

func (c *BaseCache[T]) updateLRU(element *list.Element) {
	c.lruList.MoveToFront(element)
}

func (c *BaseCache[T]) evictIfNeeded() {
	if c.lruList.Len() > c.maxEntries {
		element := c.lruList.Back()
		if element != nil {
			c.lruList.Remove(element)
		}
	}
}

type InMemoryCache[T any] struct {
	BaseCache[T]
}

func NewInMemoryCache[T any](maxEntries int) *InMemoryCache[T] {
	return &InMemoryCache[T]{
		BaseCache: NewBaseCache[T](maxEntries),
	}
}

func (c *InMemoryCache[T]) Get(key string) (T, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for e := c.lruList.Front(); e != nil; e = e.Next() {
		entry := e.Value.(cacheEntry[T])
		if entry.Key == key {
			if entry.TTL == -1 || time.Since(entry.CreatedAt) <= entry.TTL {
				entry.LastAccess = time.Now()
				entry.AccessCount++
				e.Value = entry
				c.updateLRU(e)
				return entry.Data, true
			}
			c.lruList.Remove(e)
			break
		}
	}
	var zero T
	return zero, false
}

func (c *InMemoryCache[T]) Set(key string, result T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry := cacheEntry[T]{
		Key:         key,
		Data:        result,
		CreatedAt:   time.Now(),
		TTL:         24 * time.Hour,
		LastAccess:  time.Now(),
		AccessCount: 0,
	}
	c.addToLRU(key, entry)
	c.evictIfNeeded()
}

type FileCache[T any] struct {
	BaseCache[T]
	dir          string
	ttl          time.Duration
	asyncGetChan chan fileCacheGetRequest[T]
	asyncSetChan chan fileCacheSetRequest[T]
	batchSize    int
	batchEntries map[string]T
	batchMutex   sync.Mutex
}

type fileCacheGetRequest[T any] struct {
	key      string
	response chan fileCacheGetResponse[T]
}

type fileCacheGetResponse[T any] struct {
	data T
	ok   bool
}

type fileCacheSetRequest[T any] struct {
	key    string
	result T
}

func NewFileCache[T any](dir string, ttl time.Duration, maxEntries int, batchSize int) (*FileCache[T], error) {
	if len(dir) == 0 {
		dir = os.TempDir()
	}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	cache := &FileCache[T]{
		BaseCache:    NewBaseCache[T](maxEntries),
		dir:          dir,
		ttl:          ttl,
		asyncGetChan: make(chan fileCacheGetRequest[T]),
		asyncSetChan: make(chan fileCacheSetRequest[T]),
		batchSize:    batchSize,
		batchEntries: make(map[string]T),
	}
	go cache.asyncWorker()
	go cache.batchWriter()
	return cache, nil
}

func (c *FileCache[T]) asyncWorker() {
	for {
		select {
		case req := <-c.asyncGetChan:
			data, ok := c.Get(req.key)
			req.response <- fileCacheGetResponse[T]{data: data, ok: ok}
		case req := <-c.asyncSetChan:
			c.Set(req.key, req.result)
		}
	}
}

func (c *FileCache[T]) AsyncGet(key string) (T, bool) {
	response := make(chan fileCacheGetResponse[T])
	c.asyncGetChan <- fileCacheGetRequest[T]{key: key, response: response}
	res := <-response
	return res.data, res.ok
}

func (c *FileCache[T]) AsyncSet(key string, result T) {
	c.asyncSetChan <- fileCacheSetRequest[T]{key: key, result: result}
}

func (c *FileCache[T]) Get(key string) (T, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// Check in-memory cache first
	for e := c.lruList.Front(); e != nil; e = e.Next() {
		entry := e.Value.(cacheEntry[T])
		if entry.Key == key {
			if entry.TTL == -1 || time.Since(entry.CreatedAt) <= entry.TTL {
				entry.LastAccess = time.Now()
				entry.AccessCount++
				e.Value = entry
				c.updateLRU(e)
				return entry.Data, true
			}
			// Remove expired in-memory cache entry
			c.lruList.Remove(e)
			break
		}
	}
	return c.getFromFile(key)
}

func (c *FileCache[T]) getFromFile(key string) (T, bool) {
	filePath := filepath.Join(c.dir, key+".cache.json")
	file, err := os.Open(filePath)
	if err != nil {
		var zero T
		return zero, false
	}
	defer file.Close()
	var result T
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		var zero T
		return zero, false
	}
	info, _ := file.Stat()
	if time.Since(info.ModTime()) > c.ttl && c.ttl != -1 {
		os.Remove(filePath)
		var zero T
		return zero, false
	}
	// Store the result in the in-memory cache
	c.addToLRU(key, cacheEntry[T]{
		Data:        result,
		CreatedAt:   info.ModTime(),
		TTL:         c.ttl,
		LastAccess:  time.Now(),
		AccessCount: 1,
	})
	return result, true
}

func (c *FileCache[T]) Set(key string, result T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// Store in in-memory cache
	entry := cacheEntry[T]{
		Key:         key,
		Data:        result,
		CreatedAt:   time.Now(),
		TTL:         c.ttl,
		LastAccess:  time.Now(),
		AccessCount: 0,
	}
	c.addToLRU(key, entry)
	c.evictIfNeeded()
	c.addToBatch(key, result)
}

func (c *FileCache[T]) addToBatch(key string, result T) {
	c.batchMutex.Lock()
	defer c.batchMutex.Unlock()
	c.batchEntries[key] = result
	if len(c.batchEntries) >= c.batchSize {
		c.flushBatch()
	}
}

func (c *FileCache[T]) flushBatch() {
	for key, result := range c.batchEntries {
		c.saveToFile(key, result)
	}
	c.batchEntries = make(map[string]T)
}

func (c *FileCache[T]) batchWriter() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		c.batchMutex.Lock()
		if len(c.batchEntries) > 0 {
			c.flushBatch()
		}
		c.batchMutex.Unlock()
	}
}

func (c *FileCache[T]) saveToFile(key string, result T) {
	filePath := filepath.Join(c.dir, key+".cache.json")
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	json.NewEncoder(file).Encode(result)
}

func (c *FileCache[T]) CleanUp(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.mutex.Lock()
			filepath.Walk(c.dir, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() && filepath.Ext(path) == ".cache.json" && time.Since(info.ModTime()) > c.ttl && c.ttl != -1 {
					os.Remove(path)
				}
				return nil
			})
			c.mutex.Unlock()
		}
	}()
}
