//go:build wasm

package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/syumai/workers/cloudflare/fetch"
)

var (
	cli       = fetch.NewClient()
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
)

func httpGetWithToken(url, token string) (*http.Response, error) {
	r, err := fetch.NewRequest(context.TODO(), http.MethodGet, url, nil)
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
	resp, err := cli.Do(r, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s: %s", url, err)
	}
	return resp, nil
}

func httpGet(url string) (*http.Response, error) {
	req, err := fetch.NewRequest(context.TODO(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	//I'm a browser... :)
	req.Header.Set("User-Agent", userAgent)
	resp, err := cli.Do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}
	return resp, nil
}
