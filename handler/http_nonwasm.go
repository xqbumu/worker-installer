//go:build !wasm

package handler

import (
	"fmt"
	"net/http"
)

var (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
)

func httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	//I'm a browser... :)
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}

	return resp, nil
}

func httpGetWithToken(url, token string) (*http.Response, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	//I'm a browser... :)
	r.Header.Set("User-Agent", userAgent)
	if len(token) > 0 {
		r.Header.Set("Authorization", "token "+token)
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s: %s", url, err)
	}
	return resp, nil
}
