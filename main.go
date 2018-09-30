package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/levigross/grequests"
)

// log is the main logging instance used in this application
var log = setupLogging()

func main() {
	rand.Seed(time.Now().UnixNano())

	tasks, err := ImportTasksFromJSON("taskFiles/testFileV2.json")
	if err != nil {
		log.Fatal().Msg("Unable to correctly parse tasks.") // Will call panic
	}
	log.Info().Msg("Parsed task files.")

	valid, errs := VerifyTasks(&tasks)
	if !valid {
		log.Fatal().Msgf("%+v", errs)
	}

	log.Info().Msgf("Loaded %d tasks. Waiting to run.", len(tasks))

	// Wait for the comand to start
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// Create wait group and run
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for i, task := range tasks {
		go func(i int, innerTask Task) {
			// Use wait group to hold application open
			defer wg.Done()
			log.Info().Msgf("%d Running task.", i)

			success, err := supremeCheckout(i, innerTask)
			if err != nil {
				log.Error().Msgf("%d Error checkout: %s", i, err)
			}

			log.Info().
				Int("task", i).
				Bool("success", success).
				Msg("Checkout compeleted")

		}(i, task)
	}

	wg.Wait()
}

func waitForItemMatch(i int, task Task) (SupremeItem, error) {
	supremeItems, err := GetCollectionItems(&task, true)
	if err != nil {
		return SupremeItem{}, errors.New("Error getting collection items")
	}

	if len(*supremeItems) > 0 {
		log.Debug().Msgf("Found %d items", len(*supremeItems))
		matchedItem, err := findItem(task.Item, *supremeItems)
		if err != nil {
			return SupremeItem{}, errors.New("Items in collection but unable to find items")
		}
		return matchedItem, nil
	}

	return SupremeItem{}, errors.New("No matches found in collection")
}

func supremeCheckout(i int, task Task) (bool, error) {
	var matchedItem SupremeItem // The item on the supreme site we will buy
	var err error
	session := grequests.NewSession(nil)

	// Try to find the item provided in keywords etc
	for {
		// Get items in category
		matchedItem, err = waitForItemMatch(i, task)
		if err != nil {
			log.Error().Err(err).Msgf("%d Error getting collection, sleeping.", 1)
		} else {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	log.Debug().Msgf("%d Found item %+v %s %s %s", i, matchedItem, matchedItem.color, matchedItem.name, matchedItem.url)

	// Get the ATC info from the item page
	var st string
	var sizeResponse SizeResponse
	var addURL string
	var xcsrf string
	err = retry(10, 50*time.Millisecond, func(attempt int) error {
		log.Debug().Msgf("%d Checkout attempt: %d", i, attempt)
		var err error
		st, sizeResponse, addURL, xcsrf, err = GetSizeInfo(session, matchedItem.url)
		return err
	})
	if err != nil {
		log.Error().Err(err)
		return false, err
	}
	log.Debug().Msgf("%d %s %+v %s %s", i, st, sizeResponse, addURL, xcsrf)
	time.Sleep(800 * time.Millisecond)

	// Add the item to cart
	pickedSizeID, err := PickSize(task.Item, sizeResponse)
	if err != nil {
		log.Error().Err(err).Msgf("%d Unable to pick size", i)
		return false, err
	}
	var atcSuccess bool
	err = retry(10, 50*time.Millisecond, func(attempt int) error {
		log.Debug().Msgf("%d ATC attempt: %d", i, attempt)
		var err error
		atcSuccess, err = AddToCart(session, addURL, xcsrf, st, pickedSizeID)
		return err
	})
	log.Debug().Msgf("%d ATC: %t", i, atcSuccess)
	time.Sleep(800 * time.Millisecond)

	// Checkout
	log.Debug().Msgf("%d Checking out using data %s", i, task.Account)
	var checkoutSuccess bool
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		log.Debug().Msgf("%d Checkout attempt: %d", i, attempt)
		var err error
		checkoutSuccess, err = Checkout(i, session, xcsrf, &task.Account)
		return err
	})
	log.Debug().Msgf("%d Checkout: %t", i, checkoutSuccess)

	return checkoutSuccess, nil
}
