package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/levigross/grequests"
)

// waitForItemMatchDesktop is a helper function for checkout. It waits until we find an item in the collection.
func waitForItemMatchDesktop(session *grequests.Session, task *Task) (SupremeItem, error) {
	supremeItems, err := GetCollectionItems(session, task, true)
	if err != nil {
		return SupremeItem{}, errors.New("Error getting collection items")
	}

	if len(*supremeItems) > 0 {
		task.Log().Debug().Msgf("Found %d items", len(*supremeItems))
		matchedItem, err := findItem(task.Item, *supremeItems)
		if err != nil {
			return SupremeItem{}, errors.New("Items in collection but unable to find items")
		}
		return matchedItem, nil
	}

	return SupremeItem{}, errors.New("No matches found in collection")
}

// SupremeCheckoutDesktop attempts to add to cart, waiting until it is available, and item and then check it out, on the desktop API
func (task *Task) SupremeCheckoutDesktop() (bool, error) {
	var matchedItem SupremeItem // The item on the supreme site we will buy
	var err error
	session := grequests.NewSession(nil)
	task.Log().Debug().
		Str("item", fmt.Sprintf("%+v", task.Item)).
		Str("waitTimes", fmt.Sprintf("%d %d %d", task.GetTaskRefreshRate(), task.GetTaskAtcWait(), task.GetTaskCheckoutWait())).
		Msg("Checking out item")

	// Try to find the item provided in keywords etc
	task.UpdateStatus("Looking for item")
	for {
		// Get items in category
		matchedItem, err = waitForItemMatchDesktop(session, task)
		if err != nil {
			task.Log().Error().Err(err).
				Msgf("Error getting collection, sleeping: %dms", task.GetTaskRefreshRate())
		} else {
			break
		}
		time.Sleep(time.Duration(task.GetTaskRefreshRate()) * time.Millisecond)
	}
	task.Log().Debug().Msgf("Found item %+v %s %s %s", matchedItem, matchedItem.color, matchedItem.name, matchedItem.url)
	startTime := time.Now()

	// Get the ATC info from the item page
	var st string
	var sizeResponse SizeResponse
	var addURL string
	var xcsrf string
	task.UpdateStatus("Found. Getting item details")
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Getting item info attempt: %d", attempt)
		var err error
		st, sizeResponse, addURL, xcsrf, err = GetSizeInfo(session, task, &matchedItem.url)
		return err
	})
	if err != nil {
		task.Log().Error().Err(err)
		return false, err
	}
	task.Log().Debug().Msgf("%s (%s:%+v) %s %s", st, sizeResponse.singleSizeID, sizeResponse.multipleSizes, addURL, xcsrf)

	// ADD TO CART
	task.Log().Info().
		Msgf("ATC Wait, sleeping: %dms", task.GetTaskAtcWait())
	time.Sleep(time.Duration(task.GetTaskAtcWait()) * time.Millisecond)
	task.UpdateStatus("Adding item to cart")
	pickedSizeID, err := PickSize(&task.Item, &sizeResponse)
	if err != nil {
		task.Log().Error().Err(err).Msgf("Unable to pick size")
		return false, err
	}
	// TODO: Figure out if we want ATC to continue to try or fail
	var atcSuccess bool
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("ATC attempt: %d", attempt)
		var err error
		atcSuccess, err = AddToCart(session, task, addURL, xcsrf, st, pickedSizeID)
		return err
	})
	task.Log().Debug().Msgf("ATC Results: %t", atcSuccess)
	if !atcSuccess {
		return false, nil
	}
	task.UpdateStatus("Item Added")

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
		checkoutSuccess, queueResponse, err = Checkout(session, task, xcsrf)
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
