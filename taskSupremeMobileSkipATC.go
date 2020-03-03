package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/levigross/grequests"
	"golang.org/x/net/publicsuffix"
)

// SupremeCheckoutMobileSkipATC Completes a checkout on supreme using the mobile API
func (task *Task) SupremeCheckoutMobileSkipATC() (bool, error) {
	var matchedItem *SupremeItemMobile // The item on the supreme site we will buy
	var err error
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		task.Log().Error().Err(err).Msg("Unable to create cookie jar")
		return false, err
	}

	session := grequests.NewSession(&grequests.RequestOptions{
		CookieJar: jar,
	})
	task.Log().Debug().
		Str("item", fmt.Sprintf("%+v", task.Item)).
		Str("waitTimes", fmt.Sprintf("%d %d %d", task.GetTaskRefreshRate(), task.GetTaskAtcWait(), task.GetTaskCheckoutWait())).
		Msg("Checking out item")

	// LOOK FOR ITEM
	task.UpdateStatus("Looking for item")
	for {
		// Get items in category
		matchedItem, err = waitForItemMatchMobile(session, task)
		if err != nil {
			task.Log().Error().Err(err).
				Msgf("Error getting collection, sleeping %dms.", task.GetTaskRefreshRate())
		} else {
			break
		}
		time.Sleep(time.Duration(task.GetTaskRefreshRate()) * time.Millisecond)
	}
	task.UpdateStatus("Found item")
	task.Log().Debug().Msgf("Found item %+v", matchedItem)

	// FIND STYLE
	task.UpdateStatus("Getting styles")
	var matchedStyle *Style
	for {
		// Get items in category
		matchedStyle, err = waitForStyleMatchMobile(session, task, matchedItem)
		if err != nil {
			task.Log().Error().Err(err).
				Msgf("Error matching style, sleeping %dms.", task.GetTaskRefreshRate())
		} else {
			break
		}
		time.Sleep(time.Duration(task.GetTaskRefreshRate()) * time.Millisecond)
	}
	task.UpdateStatus("Matched color / style")
	task.Log().Debug().Msgf("Matched Style: %+v", matchedStyle)

	// PICK SIZE
	pickedSizeID, isInStock, err := PickSizeMobile(&task.Item, matchedStyle)
	if err != nil {
		task.Log().Error().Err(err).Msg("Error picking size")
	}
	// WAIT FOR RESTOCK IF NOT IN STOCK
	if !isInStock {
		task.Log().Info().Msg("Item is not in stock, waiting for restock")
		task.status = "Waiting for restock"
		pickedSizeID = waitForRestock(session, task, matchedItem)
	}

	// ADD TO CART
	startTime := time.Now()
	task.Log().Info().
		Msgf("ATC Adding item to cart")
	task.Log().Debug().Msgf("item ID: %d st: %d s: %d", matchedItem.id, matchedStyle.ID, pickedSizeID)
	task.UpdateStatus("Adding item to cart")

	// cart
	// 1+item--59765%2C21347 => 1+item--59765,21347
	cartValue := "1+item--" + url.QueryEscape(fmt.Sprintf("%d,%d", pickedSizeID, matchedStyle.ID))
	// pure_cart
	// %7B%2259765%22%3A1%7D => {"59765":1}
	pureCartValue := url.QueryEscape(fmt.Sprintf("{\"%d\":1}", pickedSizeID))

	supURLHTTP, _ := url.Parse("http://www.supremenewyork.com")
	jar.SetCookies(supURLHTTP, []*http.Cookie{
		&http.Cookie{
			Domain: "www.supremenewyork.com",
			Name:   "cart",
			Path:   "/",
			Value:  cartValue,
		},
		&http.Cookie{
			Domain: "www.supremenewyork.com",
			Name:   "pure_cart",
			Path:   "/",
			Value:  pureCartValue,
		},
	})
	task.Log().Debug().Msgf("%+v", jar.Cookies(supURLHTTP))

	// CHECKOUT
	task.Log().Info().
		Msgf("Checkout Wait, sleeping: %dms", task.GetTaskCheckoutWait())
	time.Sleep(time.Duration(task.GetTaskCheckoutWait()) * time.Millisecond)
	task.UpdateStatus("Checking out")
	task.Log().Debug().
		Msgf("Checking out task: %s %s", task.Account.Person, task.Account.Address)
	var checkoutSuccess bool
	var queueResponse string
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Checkout attempt: %d", attempt)
		var err error
		checkoutSuccess, queueResponse, err = CheckoutMobile(session, task, &pureCartValue)
		return err
	})
	elapsed := time.Since(startTime)
	task.Log().Debug().
		Float64("timeElapsed", elapsed.Seconds()).
		Bool("success", checkoutSuccess).
		Str("respString", queueResponse).
		Msg("Supreme checkout completed")
	if checkoutSuccess {
		task.UpdateStatus("Checked out successfully")
	} else {
		task.UpdateStatus("Checkout failed")
		return false, err // TODO: Maybe return nil for error
	}

	// QUEUE
	task.UpdateStatus("Waiting for queue")
	var queueSuccess bool
	queueErr := retry(2, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Queue attempt: %d", attempt)
		var err error
		queueSuccess, err = Queue(session, task, queueResponse)
		return err
	})
	if queueErr != nil {
		return queueSuccess, queueErr
	}

	task.UpdateStatus("Completed")
	return queueSuccess, nil
}
