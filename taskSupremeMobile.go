package main

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/levigross/grequests"
)

// waitForItemMatchMobile is a helper function for checkout. It waits until we find an item in the collection.
func waitForItemMatchMobile(session *grequests.Session, task *Task) (SupremeItemMobile, error) {
	itemsMobile, err := GetCollectionItemsMobile(session, task)
	if err != nil {
		return SupremeItemMobile{}, errors.New("Error getting collection items")
	}

	if len(*itemsMobile) > 0 {
		matchedItem, err := findItemMobile(task.Item, itemsMobile)
		if err != nil {
			return SupremeItemMobile{}, errors.New("Items in collection but unable to find items")
		}
		return matchedItem, nil
	}

	return SupremeItemMobile{}, errors.New("No matches found in collection")
}

// SupremeCheckoutMobile Completes a checkout on supreme using the mobile API
func (task *Task) SupremeCheckoutMobile() (bool, error) {
	var matchedItem SupremeItemMobile // The item on the supreme site we will buy
	var err error
	session := grequests.NewSession(nil)
	task.Log().Debug().
		Str("item", fmt.Sprintf("%+v", task.Item)).
		Msg("Checking out item")

		// Try to find the item provided in keywords etc
	task.UpdateStatus("Looking for item")
	for {
		// Get items in category
		matchedItem, err = waitForItemMatchMobile(session, task)
		if err != nil {
			task.Log().Error().Err(err).
				Msgf("Error getting collection, sleeping.")
		} else {
			break
		}
		time.Sleep(time.Duration(appSettings.RefreshWait) * time.Millisecond)
	}
	task.UpdateStatus("Found item")
	task.Log().Debug().Msgf("Found item %+v", matchedItem)

	startTime := time.Now()
	task.UpdateStatus("Getting styles")
	var styles []Style
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Getting item info attempt: %d", attempt)
		var err error
		styles, err = GetSizeInfoMobile(session, task, matchedItem)
		return err
	})
	if err != nil {
		task.Log().Error().Err(err)
		return false, err
	}

	var matchedStyle Style
	foundMatchedStyle := false
	for _, style := range styles {
		if checkColor(task.Item.Color, style.Name) {
			matchedStyle = style
			foundMatchedStyle = true
			break
		}
	}

	if !foundMatchedStyle {
		task.Log().Error().Msg("Unable to find color")
		return false, errors.New("Unable to find color")
	}
	task.UpdateStatus("Matched color / style")
	task.Log().Debug().Msgf("Matched Style: %+v", matchedStyle)

	pickedSizeID, err := PickSizeMobile(&task.Item, &matchedStyle)
	if err != nil {
		task.Log().Error().Err(err).Msg("Error picking size")
	}

	time.Sleep(time.Duration(appSettings.AtcWait) * time.Millisecond)
	task.Log().Debug().Msgf("item ID: %d st: %d s: %d", matchedItem.id, matchedStyle.ID, pickedSizeID)
	task.UpdateStatus("Adding item to cart")
	var atcSuccess bool
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("ATC attempt: %d", attempt)
		var err error
		atcSuccess, err = AddToCartMobile(session, task, matchedItem.id, matchedStyle.ID, pickedSizeID)
		return err
	})
	task.Log().Debug().Msgf("ATC Results: %t", atcSuccess)
	if !atcSuccess {
		return false, nil
	}

	// Purecart implementation instead of building cookie sub our selves
	// supremeURL, _ := url.Parse("http://www.supremenewyork.com")
	// task.Log().Debug().Msgf("%+v", session.HTTPClient.Jar.Cookies(supremeURL))

	time.Sleep(time.Duration(appSettings.CheckoutWait) * time.Millisecond)
	task.UpdateStatus("Checking out")
	cookieSub := url.QueryEscape(fmt.Sprintf("{\"%d\":1}", pickedSizeID))
	var checkoutSuccess bool
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Checkout attempt: %d", attempt)
		var err error
		checkoutSuccess, err = CheckoutMobile(session, task, cookieSub)
		return err
	})
	elapsed := time.Since(startTime)

	task.UpdateStatus("Completed")
	// Status and send info
	task.Log().Debug().
		Float64("timeElapsed", elapsed.Seconds()).
		Bool("success", checkoutSuccess).
		Msgf("Supreme checkout completed")
	if checkoutSuccess {
		task.UpdateStatus("Checked out successfully")
	} else {
		task.UpdateStatus("Checkout failed")
	}

	return checkoutSuccess, nil // TODO: Replace with real value
}
