package main

// import (
// 	"fmt"
// 	"net/url"
// 	"testing"

// 	"github.com/levigross/grequests"
// 	"github.com/stretchr/testify/assert"
// )

// func TestPost(t *testing.T) {
// 	localRo := grequests.RequestOptions{
// 		UserAgent: mobileUserAgent,
// 		Headers: map[string]string{
// 			"accept-language":  "en-US,en;q=0.9",
// 			"accept":           "application/json",
// 			"reffer":           "http://www.supremenewyork.com/mobile",
// 			"x-requested-with": "XMLHttpRequest",
// 		},
// 		Data: map[string]string{"One": "Two"},
// 	}
// 	resp, _ := grequests.Post("http://httpbin.org/post", &localRo)

// 	t.Log(resp.String())
// }

// func TestGetCollectionItemsMobile(t *testing.T) {
// 	session := grequests.NewSession(nil)
// 	task := testTask()

// 	items, err := GetCollectionItemsMobile(session, &task)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(items)

// }

// func TestGetSizeInfoMobile(t *testing.T) {
// 	session := grequests.NewSession(nil)
// 	task := testTask()
// 	item := SupremeItemMobile{"Corduroy Shirt", 171868, ""}

// 	styles, err := GetSizeInfoMobile(session, &task, item)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(styles)

// }

// func TestAddToCartMobile(t *testing.T) {
// 	session := grequests.NewSession(nil)
// 	task := testTask()
// 	success, err := AddToCartMobile(session, &task, 171870, 21889, 61459)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	t.Log(success)
// }

// func TestCheckoutMobile(t *testing.T) {
// 	session := grequests.NewSession(nil)
// 	task := testTask()
// 	success, err := AddToCartMobile(session, &task, 171870, 21889, 61459)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	assert.True(t, success)
// 	success, err = CheckoutMobile(session, &task, url.QueryEscape(fmt.Sprintf("{\"%d\":1}", 61459)))
// }
