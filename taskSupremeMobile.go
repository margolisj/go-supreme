package main

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/levigross/grequests"
)

// waitForItemMatchMobile is a helper function for checkout. It waits until we find an item in the collection.
func waitForItemMatchMobile(session *grequests.Session, task *Task) (*SupremeItemMobile, error) {
	itemsMobile, err := GetCollectionItemsMobile(session, task)
	if err != nil {
		return &SupremeItemMobile{}, errors.New("Error getting collection items")
	}
	// task.Log().Info().Msgf("%+v", itemsMobile)

	if len(*itemsMobile) > 0 {
		matchedItem, err := findItemMobile(task.Item, itemsMobile)
		if err != nil {
			return &SupremeItemMobile{}, errors.New("Items in collection but unable to find items")
		}
		return &matchedItem, nil
	}

	return &SupremeItemMobile{}, errors.New("No matches found in collection")
}

func waitForStyleMatchMobile(session *grequests.Session, task *Task, matchedItem *SupremeItemMobile) (*Style, error) {
	styles, err := GetSizeInfoMobile(session, task, matchedItem)
	if err != nil {
		return &Style{}, errors.New("Error getting styles")
	}
	// task.Log().Info().Msgf("%+v", styles)

	if len(*styles) > 0 {
		for _, style := range *styles {
			if checkColor(task.Item.Color, style.Name) {
				return &style, nil
			}
		}
	}

	return &Style{}, errors.New("No matches found for style")
}

func waitForRestock(session *grequests.Session, task *Task, matchedItem *SupremeItemMobile) int {
	for {
		task.Log().Debug().
			Msgf("Waiting for restock, sleeping %dms.", appSettings.RefreshWait)
		time.Sleep(time.Duration(appSettings.RefreshWait) * time.Millisecond)

		matchedStyle, err := waitForStyleMatchMobile(session, task, matchedItem)
		if err != nil {
			task.Log().Error().Err(err).Msg("Error matching style")
			continue
		}

		pickedSizeID, isInStock, err := PickSizeMobile(&task.Item, matchedStyle)
		if err != nil {
			task.Log().Error().Err(err).Msg("Error picking size")
		}

		if !isInStock {
			task.Log().Debug().Msg("Item not in stock")
			continue
		}

		return pickedSizeID
	}
}

// SupremeCheckoutMobile Completes a checkout on supreme using the mobile API
func (task *Task) SupremeCheckoutMobile() (bool, error) {
	var matchedItem *SupremeItemMobile // The item on the supreme site we will buy
	var err error
	session := grequests.NewSession(nil)
	task.Log().Debug().
		Str("item", fmt.Sprintf("%+v", task.Item)).
		Msg("Checking out item")

	// LOOK FOR ITEM
	task.UpdateStatus("Looking for item")
	for {
		// Get items in category
		matchedItem, err = waitForItemMatchMobile(session, task)
		if err != nil {
			task.Log().Error().Err(err).
				Msgf("Error getting collection, sleeping %dms.", appSettings.RefreshWait)
		} else {
			break
		}
		time.Sleep(time.Duration(appSettings.RefreshWait) * time.Millisecond)
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
				Msgf("Error matching style, sleeping %dms.", appSettings.RefreshWait)
		} else {
			break
		}
		time.Sleep(time.Duration(appSettings.RefreshWait) * time.Millisecond)
	}
	task.UpdateStatus("Matched color / style")
	task.Log().Debug().Msgf("Matched Style: %+v", matchedStyle)

	// PICK SIZE
	pickedSizeID, isInStock, err := PickSizeMobile(&task.Item, matchedStyle)
	if err != nil {
		task.Log().Error().Err(err).Msg("Error picking size")
	}
	// WAIT FOR RESTOCK IF NOT INSTOCK
	if !isInStock {
		task.Log().Info().Msg("Item is not in stock, waiting for restock")
		task.status = "Waiting for restock"
		pickedSizeID = waitForRestock(session, task, matchedItem)
	}

	// ADD TO CART
	startTime := time.Now()
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

	// CHECKOUT
	time.Sleep(time.Duration(appSettings.CheckoutWait) * time.Millisecond)
	task.UpdateStatus("Checking out")
	cookieSub := url.QueryEscape(fmt.Sprintf("{\"%d\":1}", pickedSizeID))
	var checkoutSuccess bool
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Checkout attempt: %d", attempt)
		var err error
		checkoutSuccess, err = CheckoutMobile(session, task, &cookieSub)
		return err
	})
	elapsed := time.Since(startTime)

	// Status and send info
	task.UpdateStatus("Completed")
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
