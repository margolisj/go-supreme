package main

import (
	"fmt"
	"log"
	"time"

	"github.com/levigross/grequests"
)

func main() {
	task := realTask()
	success, err := supremeCheckout(task)

	if err != nil {
		log.Println(err)
	}

	if !success {
		log.Println("Checkout was not successful")
	}

}

func supremeCheckout(task task) (bool, error) {
	taskItem := task.items[0]

	supremeItems := GetCollectionItems(taskItem, true)

	matchedItem, err := FindItem(taskItem, *supremeItems)
	if err != nil {
		log.Fatal("Error matching item", err)
		return false, err
	}
	fmt.Printf("Found item %s", matchedItem)

	session := grequests.NewSession(nil)
	st, sizes, addURL, xcsrf := GetSizeInfo(session, matchedItem.url)
	fmt.Printf("%s %s %s\n", st, addURL, xcsrf)

	time.Sleep(2 * time.Second)

	atdSuccess := AddToCart(session, addURL, xcsrf, st, (*sizes)["Medium"])
	fmt.Printf("ATC: %t\n", atdSuccess)

	time.Sleep(2 * time.Second)

	acc := realData()
	fmt.Printf("Using data %s \n", acc)
	checkoutSuccess := Checkout(session, xcsrf, &acc)
	fmt.Printf("Checkout: %t\n", checkoutSuccess)

	return true, nil
}

// {"status":"queued","slug":"q7j84cuad93wnyrg0"}
