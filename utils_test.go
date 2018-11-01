package main

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"testing"
	"time"

	"4d63.com/tz"
	"golang.org/x/net/publicsuffix"

	"github.com/levigross/grequests"
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

func TestRetryWithError(t *testing.T) {
	var attemptVal int
	retry(10, 50*time.Millisecond, func(attempt int) error {
		attemptVal = attempt
		return errors.New("")
	})
	assert.Equal(t, 1, attemptVal)
}

func TestReadTimeFromString(t *testing.T) {
	str := "2018-10-10T14:59:30.000Z"
	rTime, err := time.Parse(time.RFC3339, str)
	loc, err := tz.LoadLocation("America/New_York")
	rTime.In(loc)
	if err != nil {
		t.Error(err)
	}
	t.Log(rTime)
}

func TestJarChange(t *testing.T) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Error(err)
	}

	localRo := &grequests.RequestOptions{
		UserAgent: mobileUserAgent,
		CookieJar: jar,
	}
	localRo2 := &grequests.RequestOptions{
		UserAgent: mobileUserAgent,
	}
	session := *grequests.NewSession(localRo)
	resp, _ := session.Get("https://httpbin.org/cookies/set/cookie/hungry", localRo2)
	t.Log(resp)
	//httpbin.org

	httpbinURL, _ := url.Parse("https://httpbin.org")

	t.Log(jar.Cookies(httpbinURL))

	jar.SetCookies(httpbinURL, []*http.Cookie{
		&http.Cookie{
			Domain: "httpbin.org",
			Name:   "monster",
			Path:   "/",
			Value:  "good",
		},
	})

	resp, _ = session.Get("https://httpbin.org/cookies", localRo2)
	t.Log(resp)

}

func TestJarChangeSupreme(t *testing.T) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Error(err)
	}

	session := *grequests.NewSession(&grequests.RequestOptions{
		CookieJar: jar,
	})

	task := &Task{
		Item: taskItem{
			Keywords: []string{
				"brieFs",
				"BoXeR",
			},
			Size:     "medium",
			Color:    "white",
			Category: "accessories",
		},
	}

	AddToCartMobile(&session, task, 171745, 21347, 59765)
	supURL, _ := url.Parse("http://www.supremenewyork.com")
	t.Log(jar.Cookies(supURL))
}

// func TestATCSkipSupreme(t *testing.T) {
// 	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	session := *grequests.NewSession(&grequests.RequestOptions{
// 		UserAgent: mobileUserAgent,
// 		CookieJar: jar,
// 	})

// 	task := testTask()
// 	task.Item = taskItem{
// 		Keywords: []string{
// 			"brieFs",
// 			"BoXeR",
// 		},
// 		Size:     "medium",
// 		Color:    "white",
// 		Category: "accessories",
// 	}

// 	session.Get("https://www.supremenewyork.com/mobile/", nil)

// 	supURLHTTP, _ := url.Parse("http://www.supremenewyork.com")
// 	supURLHTTPS, _ := url.Parse("https://www.supremenewyork.com")
// 	t.Log(jar.Cookies(supURLHTTP))
// 	t.Log(jar.Cookies(supURLHTTPS))

// 	// cart
// 	// 1+item--59765%2C21347 => 1+item--59765,21347
// 	cartValue := "1+item--" + url.QueryEscape(fmt.Sprintf("%d,%d", 59765, 21347))

// 	//pure_cart
// 	// %7B%2259765%22%3A1%7D => {"59765":1}
// 	pureCartValue := url.QueryEscape(fmt.Sprintf("{\"%d\":1}", 59765))

// 	// cookies := jar.Cookies(supURL)
// 	jar.SetCookies(supURLHTTP, []*http.Cookie{
// 		&http.Cookie{
// 			Domain: "www.supremenewyork.com",
// 			Name:   "cart",
// 			Path:   "/",
// 			Value:  cartValue,
// 		},
// 		&http.Cookie{
// 			Domain: "www.supremenewyork.com",
// 			Name:   "pure_cart",
// 			Path:   "/",
// 			Value:  pureCartValue,
// 		},
// 	})

// 	t.Log(jar.Cookies(supURLHTTP))
// 	t.Log(jar.Cookies(supURLHTTPS))
// 	success, checkoutResp, err := CheckoutMobile(&session, &task, &pureCartValue)

// 	t.Log(success)
// 	t.Log(err)
// }

func TestUKKeywords(t *testing.T) {
	items := []taskItem{
		taskItem{
			Keywords: []string{
				"trust",
				"killer",
			},
			Size:     "medium",
			Color:    "white",
			Category: "t-shirts",
		},
		taskItem{
			Keywords: []string{
				"trust",
				"killer",
			},
			Size:     "medium",
			Color:    "heather grey",
			Category: "t-shirts",
		},
		taskItem{
			Keywords: []string{
				"L/S",
				"killer",
			},
			Size:     "medium",
			Color:    "white",
			Category: "t-shirts",
		},
	}

	task := testTask()
	proxyString := "81.130.135.142:48057"
	proxyURL, err := url.Parse("http://" + proxyString) // Proxy URL
	if err != nil {
		t.Error(err)
	}

	localRo := &grequests.RequestOptions{
		UserAgent: mobileUserAgent,
		Proxies: map[string]*url.URL{
			"http":  proxyURL,
			"https": proxyURL,
		},
	}

	session := *grequests.NewSession(localRo)
	tresp, _ := session.Get("https://api.ipify.org?format=json", nil)
	t.Logf("Current IP: %s", tresp.String())

	for k, item := range items {
		task.Item = item
		itemMobile, err := waitForItemMatchMobile(&session, &task)
		if err != nil {
			t.Logf("Item %d keyword couldn't be found", k)
		} else {
			t.Logf("Item %d keyword was found: %+v", k, itemMobile)
		}

		style, err := waitForStyleMatchMobile(&session, &task, itemMobile)
		if err != nil {
			t.Logf("Item %d couldn't be found style / color", k)
		} else {
			t.Logf("Item %d was found style / color: %+v", k, style)
		}

		ID, _, err := PickSizeMobile(&task.Item, style)
		if err != nil {
			t.Logf("Item %d couldn't be found size", k)
		} else {
			t.Logf("Item %d was found size: %d", k, ID)
		}

	}
}
