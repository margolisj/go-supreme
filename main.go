package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
)

func main() {
	tasks, err := ImportTasksFromJSON("testFile.txt")

	if err != nil {
		log.Fatal("Unable to correctly parse tasks.")
		panic(err)
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
			log.Infof("%d Running task.", &i)

			success, err := supremeCheckout(i, innerTask)
			if err != nil {
				log.Fatal("Error after checkout", err)
			}

			log.WithFields(log.Fields{
				"thread":  i,
				"success": success,
			}).Info("Checkout has compelted")

		}(i, task)

	}

	wg.Wait()

}

func supremeCheckout(i int, task Task) (bool, error) {
	taskItem := task.Item

	var matchedItem SupremeItem
	var err error
	for {
		supremeItems := GetCollectionItems(taskItem, true)
		matchedItem, err = FindItem(taskItem, *supremeItems)
		if err != nil {
			log.Warnf("%d Error matching item, sleeping: %s", i, err.Error())
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}

	log.Debugf("%d Found item %s", i, matchedItem)

	session := grequests.NewSession(nil)
	st, sizes, addURL, xcsrf := GetSizeInfo(session, matchedItem.url)
	log.Debugf("%d %s %s %s", i, st, addURL, xcsrf)

	time.Sleep(1000 * time.Millisecond)

	atdSuccess := AddToCart(session, addURL, xcsrf, st, (*sizes)["Medium"])
	log.Debugf("%d ATC: %t", i, atdSuccess)

	time.Sleep(1300 * time.Millisecond)

	log.Debugf("%d Using data %s", i, task.Account)
	checkoutSuccess := Checkout(session, xcsrf, &task.Account)
	log.Debugf("%d Checkout: %t", i, checkoutSuccess)

	return true, nil
}
