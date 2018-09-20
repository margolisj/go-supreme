package main

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindItem(t *testing.T) {
	taskItem := taskItem{
		[]string{"hanes", "boxer"},
		"accessories",
		"Medium",
		"white",
	}

	supremeItems := SupremeItems{SupremeItem{
		"Supreme®/Hanes® Boxer Briefs (4 Pack)",
		"White",
		"shop/accessories/nckme38ul/iimyp2ogd",
	}}

	_, err := FindItem(taskItem, supremeItems)

	if err != nil {
		t.Fail("Unable to find item")
	}
}

// func TestFailure

// func TestPost(t *testing.T) {
// 	localRo := &grequests.RequestOptions{
// 		UserAgent: sharedUserAgent,
// 		Headers: map[string]string{
// 			"accept":           "application/json",
// 			"accept-encoding":  "gzip, deflate, br",
// 			"accept-language":  "en-US,en;q=0.9",
// 			"origin":           "https://www.whatever.com",
// 			"referer":          "https://www.whatever.com/checkout",
// 			"x-requested-with": "XMLHttpRequest",
// 			"content-type":     "application/x-www-form-urlencoded; charset=UTF-8",
// 			"dnt":              "1",
// 		},
// 		Data: map[string]string{
// 			"utf8":                    "✓",
// 			"asec":                    "Rmasn",
// 			"same_as_billing_address": "1",
// 			"store_credit_id":         "",
// 			"credit_card[nlb]":        "1285 4827 5948 2017",
// 			"order[terms]":            "1",
// 			"credit_card[vval]":       "",
// 		},
// 	}

// 	resp, err := grequests.Post("http://httpbin.org/post", localRo)

// 	if err != nil {
// 		t.Log("Cannot post: ", err)
// 	}

// 	if resp.Ok != true {
// 		t.Log("Request did not return OK")
// 	}

// 	values := resp.RawResponse.Request.PostForm
// 	t.Log(values)
// 	// for k, v := range values {
// 	// 	t.Log()
// 	// }

// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(resp.RawResponse.Body)
// 	newStr := buf.String()
// 	t.Log(newStr)

// 	// var respJSON interface{}
// 	// resp.JSON(&respJSON)
// 	// t.Log(respJSON)
// }

func encodePostValues(postValues map[string]string) string {
	urlValues := &url.Values{}

	for key, value := range postValues {
		urlValues.Set(key, value)
	}

	return urlValues.Encode() // This will sort all of the string values
}

func TestEncoding2(t *testing.T) {
	data := map[string]string{
		"utf8":                    "✓",
		"asec":                    "Rmasn",
		"same_as_billing_address": "1",
		"store_credit_id":         "",
		"credit_card[nlb]":        "1234 5678 9012 3456",
		"order[terms]":            "1",
		"credit_card[vval]":       "",
	}

	encodedVals := encodePostValues(data)

	assert.Equal(t, "asec=Rmasn&credit_card%5Bnlb%5D=1234+5678+9012+3456&credit_card%5Bvval%5D=&order%5Bterms%5D=1&same_as_billing_address=1&store_credit_id=&utf8=%E2%9C%93", encodedVals)
}

// func TestProxy(t *testing.T) {

// 	proxyURL, err := url.Parse("http://US-30m.geosurf.io:10000")
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	resp, err := grequests.Get("http://httpbin.org/ip",
// 		&grequests.RequestOptions{Proxies: map[string]*url.URL{proxyURL.Scheme: proxyURL}})

// 	if err != nil {
// 		log.Println(err)
// 	}

// 	if resp.Ok != true {
// 		log.Println("Request did not return OK")
// 	}

// 	log.Println(resp)
// }

func TestCheckoutResponses(t *testing.T) {
	queuedStatus := []byte(`{"status":"queued","slug":"q7j84cuad93wnyrg0"}`)
	var dat map[string]interface{}
	var err error

	if err = json.Unmarshal(queuedStatus, &dat); err != nil {
		panic(err)
	}
	t.Log(dat)

	if val, ok := dat["status"]; ok {
		assert.Equal(t, "queued", val)

		if val, ok := dat["slug"]; ok {
			assert.Equal(t, "q7j84cuad93wnyrg0", val)
		} else {
			t.Fail()
		}

	} else {
		t.Fail()
	}

	failedCreditCard := []byte(`{"status":"failed","cart":[{"size_id":"59765","in_stock":true}],"errors":{"order":"","credit_card":"number is not a valid credit card number"}}`)
	if err = json.Unmarshal(failedCreditCard, &dat); err != nil {
		panic(err)
	}
	if val, ok := dat["status"]; ok {
		assert.Equal(t, "failed", val)

		if val, ok := dat["errors"]; ok {
			if str, ok := val.(string); ok {
				assert.True(t, strings.Contains(str, "not a valid credit card "))
			} else {
				t.Fail()
			}
		} else {
			t.Fail()
		}

	} else {
		t.Fail()
	}

}
