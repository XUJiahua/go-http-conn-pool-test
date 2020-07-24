package main

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var _httpCli = &http.Client{
	Timeout: time.Duration(15) * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:          1,
		MaxIdleConnsPerHost:   1,
		MaxConnsPerHost:       1,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func get(url string) {
	resp, err := _httpCli.Get(url)
	if err != nil {
		// do nothing
		return
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// do nothing
		return
	}
}

func TestLongShort(t *testing.T) {
	go func() {
		for i := 0; i < 1000; i++ {
			go get("http://192.168.33.10:8087")
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			go get("http://192.168.33.10:8088")
		}
	}()

	time.Sleep(time.Second * 10)
}

func TestLongLong(t *testing.T) {
	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(time.Second)
			}
			go get("http://192.168.33.10:8087")
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			if i%10 == 0 {
				time.Sleep(time.Second)
			}
			go get("http://192.168.33.10:8089")
		}
	}()

	time.Sleep(time.Second * 10)
}

func TestLong(t *testing.T) {
	go func() {
		for i := 0; i < 1000; i++ {
			go get("http://192.168.33.10:8087")
		}
	}()

	time.Sleep(time.Second * 10)
}
