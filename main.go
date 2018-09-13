package main

import (
	"fmt"
	"time"

	"github.com/levigross/grequests"
)

func testData() account {
	p := person{
		"Jax",
		"Blax",
		"none@none.com",
		"215-834-1857", // Phone numbers have dashes
	}

	a := address{
		"102 Broad Street",
		"",
		"12345",
		"Philadeliphia",
		"PA",
		"USA",
	}

	c := card{
		"visa",                // Don't think this matters for desktop
		"1285 4827 5948 2017", //TODO: These cards must have spaces
		"02",
		"2019", // 4 digit dates
		"847",  //
	}

	return account{p, a, c}
}

func testTask() task {
	i := taskItem{
		[]string{"white", "boxers"},
		"new",
		"M",
	}

	return task{
		"Task1",
		[]taskItem{i},
		false,
		"Waiting",
	}
}

func realData() account {
	p := person{
		"Josh",
		"Margolis",
		"quickslash85@gmail.com",
		"610-212-1488", // Phone numbers have dashes
	}

	a := address{
		"ASDF 111 E Levering Mill Rd",
		"",
		"19004",
		"Bala Cynywd",
		"PA",
		"USA",
	}

	c := card{
		"visa",                // Don't think this matters for desktop
		"5466 2803 1652 6908", //TODO: These cards must have spaces
		"04",
		"2019", // 4 digit dates
		"761",  //
	}

	return account{p, a, c}
}

func realTask() task {
	i := taskItem{
		[]string{"white", "boxers"},
		"accessories",
		"Medium",
	}

	return task{
		"Checkout Task",
		[]taskItem{i},
		false,
		"Waiting",
	}
}

func main() {
	task := realTask()
	success, err := supremeCheckout(task)

	if err != nil {
		log.Println(err)
	}

	if !success {
		log.PrintLn("Checkout was not successful")
	}
}

func supremeCheckout(task task) (bool, err) {
	session := grequests.NewSession(nil)

	GetCollectionItems(task.items[0], true)

	st, sizes, addURL, xcsrf := GetSizeInfo(session, "/shop/accessories/mfb130dig/evr2ecnmx")
	fmt.Printf("%s %s %s\n", st, addURL, xcsrf)

	time.Sleep(2 * time.Second)

	atdSuccess := AddToCart(session, addURL, xcsrf, st, (*sizes)["Medium"])
	fmt.Printf("ATC: %t\n", atdSuccess)

	time.Sleep(2 * time.Second)

	checkoutSuccess := Checkout(session, xcsrf, testData())
	fmt.Printf("Checkout: %t\n", checkoutSuccess)
}
