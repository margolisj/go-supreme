package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
)

func setupLogging() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})

	// Create file and set output to both if possible
	filename := fmt.Sprintf("logs/logfile-%d.log", time.Now().Unix())
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	} else {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}

}

func main() {
	setupLogging()
	rand.Seed(time.Now().UnixNano())

	tasks, err := ImportTasksFromJSON("taskFiles/testFile.json")
	if err != nil {
		log.Panic("Unable to correctly parse tasks.") // Will call panic
	}

	valid, errs := VerifyTasks(&tasks)
	if !valid {
		log.Panic(errs)
	}

	log.Infof("Loaded %d tasks. Waiting to run.", len(tasks))

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
			log.Infof("%d Running task.", i)

			success, err := supremeCheckout(i, innerTask)
			if err != nil {
				log.Errorf("%d Error checkout: %s", i, err)
			}

			log.WithFields(log.Fields{
				"thread":  i,
				"success": success,
			}).Info("Checkout has compeleted")

		}(i, task)

	}

	wg.Wait()

}

func supremeCheckout(i int, task Task) (bool, error) {
	var matchedItem SupremeItem // The item on the supreme site we will buy
	var err error
	session := grequests.NewSession(nil)

	// Try to find the item provided in keywords etc
	for {
		supremeItems, err := GetCollectionItems(task.Item, true)

		if err != nil {
			log.Errorf("%d Error getting collection", 1)
		} else {
			if len(*supremeItems) > 0 {
				matchedItem, err = findItem(task.Item, *supremeItems)
			}
			if err != nil {
				log.Warnf("%d Error matching item, sleeping: %s", i, err.Error())
			} else {
				break
			}
		}

		time.Sleep(300 * time.Millisecond)
	}
	log.Debugf("%d Found item %s", i, matchedItem)

	// Get the ATC info from the item page
	var st string
	var sizeResponse SizeResponse
	var addURL string
	var xcsrf string
	err = retry(10, 50*time.Millisecond, func(attempt int) error {
		log.Debugf("%d Checkout attempt: %d", i, attempt)
		var err error
		st, sizeResponse, addURL, xcsrf, err = GetSizeInfo(session, matchedItem.url)
		return err
	})
	if err != nil {
		log.Error(err)
		return false, err
	}
	log.Debugf("%d %s %v %s %s", i, st, sizeResponse, addURL, xcsrf)
	time.Sleep(600 * time.Millisecond)

	// Add the item to cart
	pickedSizeID, err := PickSize(task.Item, sizeResponse)
	if err != nil {
		log.Errorf("%d Unable to find size", i)
		return false, err
	}
	var atcSuccess bool
	err = retry(10, 50*time.Millisecond, func(attempt int) error {
		log.Debugf("%d ATC attempt: %d", i, attempt)
		var err error
		atcSuccess, err = AddToCart(session, addURL, xcsrf, st, pickedSizeID)
		return err
	})
	log.Debugf("%d ATC: %t", i, atcSuccess)
	time.Sleep(600 * time.Millisecond)

	// Checkout
	log.Debugf("%d Checking out using data %s", i, task.Account)
	var checkoutSuccess bool
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		log.Debugf("%d Checkout attempt: %d", i, attempt)
		var err error
		checkoutSuccess, err = Checkout(i, session, xcsrf, &task.Account)
		return err
	})
	log.Debugf("%d Checkout: %t", i, checkoutSuccess)

	return checkoutSuccess, nil
}
