package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetryJitterSpeed(t *testing.T) {
	pastFirstRetry := false
	var wg sync.WaitGroup
	wg.Add(10)
	lastAttempt := time.Now()

	retry(10, 50*time.Millisecond, func(attempt int) error {
		defer wg.Done()
		if !pastFirstRetry {
			pastFirstRetry = true
			lastAttempt = time.Now()
			return errors.New("Getting to retry")
		}
		elapsed := time.Now().Sub(lastAttempt)
		assert.True(t, elapsed < 100*time.Millisecond, "elaspsed time is less than 50 ms")
		t.Logf("%s %d %f", time.Now().UTC(), attempt, elapsed.Seconds())
		lastAttempt = time.Now()
		return errors.New("Error so this repeats")
	})

}

func TestReadTimeFromSTring(t *testing.T) {
	str := "2018-10-10T14:59:30.000Z"
	rTime, err := time.Parse(time.RFC3339, str)

	if err != nil {
		fmt.Println(err)
	}
	t.Log(rTime)
}

// func TestProxy(t *testing.T) {
// 	proxyHTTP, err := url.Parse("http://45.79.136.67:3375") // Proxy URL
// 	if err != nil {
// 		tempLog.Panicln(err)
// 	}
// 	proxyHTTPS, err := url.Parse("https://45.79.136.67:3375") // Proxy URL
// 	if err != nil {
// 		tempLog.Panicln(err)
// 	}

// 	resp, err := grequests.Get("https://api.ipify.org?format=json",
// 		&grequests.RequestOptions{
// 			Proxies: map[string]*url.URL{
// 				proxyHTTP.Scheme:  proxyHTTP,
// 				proxyHTTPS.Scheme: proxyHTTP,
// 			},
// 		},
// 	)

// 	if err != nil {
// 		tempLog.Println(err)
// 	}

// 	if resp.Ok != true {
// 		tempLog.Println("Request did not return OK")
// 	}

// 	tempLog.Println(resp)
// }

// func TestOtherProxy(t *testing.T) {

// 	u, err := url.Parse("http://45.79.159.35:8034")
// 	if err != nil {
// 		panic(err)
// 	}
// 	tr := &http.Transport{
// 		Proxy: http.ProxyURL(u),
// 	}
// 	client := &http.Client{Transport: tr}
// 	resp, err := client.Get("https://api.ipify.org?format=json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// 	dump, err := httputil.DumpResponse(resp, true)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("%q", dump)
// }
