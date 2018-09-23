package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/levigross/grequests"
)

func main() {
	tasks, err := ImportTasksFromJSON("testFile.txt")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded %d tasks. Waiting to run.\n", len(tasks))

	// tasks = []FullTask{tasks[0]}
	// tasks = []FullTask{testFullTask()}

	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for i, fullTask := range tasks {

		go func(i int, t FullTask) {
			defer wg.Done()
			fmt.Printf("%d Running task\n", &i)
			success, err := supremeCheckout(i, t.Task, t.Account)
			if err != nil {
				log.Println(err)
			}

			if !success {
				fmt.Printf("%d Checkout was unsuccessful", i)
			} else {
				fmt.Printf("%d Checkout was successful", i)
			}

		}(i, fullTask)

	}

	wg.Wait()

	// task := realTask()
	// success, err := supremeCheckout(task)

}

func supremeCheckout(i int, task Task, acc account) (bool, error) {
	taskItem := task.Item

	var matchedItem SupremeItem
	var err error
	for {
		supremeItems := GetCollectionItems(taskItem, true)
		matchedItem, err = FindItem(taskItem, *supremeItems)
		if err != nil {
			fmt.Printf("%d Error matching item, sleeping: %s\n", i, err.Error())
			time.Sleep(500 * time.Millisecond)
		} else {
			break
		}
	}

	fmt.Printf("%d Found item %s", i, matchedItem)

	session := grequests.NewSession(nil)
	st, sizes, addURL, xcsrf := GetSizeInfo(session, matchedItem.url)
	fmt.Printf("%d %s %s %s\n", i, st, addURL, xcsrf)

	time.Sleep(1000 * time.Millisecond)

	atdSuccess := AddToCart(session, addURL, xcsrf, st, (*sizes)["Medium"])
	fmt.Printf("%d ATC: %t\n", i, atdSuccess)

	time.Sleep(1300 * time.Millisecond)

	fmt.Printf("%d Using data %s \n", i, acc)
	checkoutSuccess := Checkout(session, xcsrf, &acc)
	fmt.Printf("%d Checkout: %t\n", i, checkoutSuccess)

	return true, nil
}
