package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/levigross/grequests"
	"github.com/rs/zerolog"
)

// Person is a struct modeling personal information
type Person struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"` // Phone numbers must have dashes
}

// Card is a struct modeling a credit card
type Card struct {
	Number   string `json:"number"`   // Card numbers must have spaces "XXXX XXXX XXXX XXXX"
	Month    string `json:"month"`    // Two digit month, ex. 03
	Year     string `json:"year"`     // 4 digit year, ex. 2019
	Cvv      string `json:"cvv"`      // 3 digit or 4 digit
	Cardtype string `json:"cardtype"` // Don't think this matters for desktop, also should be able to figure this out without user entry
}

// Address is a struct modeling an address
type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Zipcode  string `json:"zipcode"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

// Account is a checkout account, a person, address and card
type Account struct {
	Person  Person  `json:"person"`
	Address Address `json:"address"`
	Card    Card    `json:"card"`
}

type taskItem struct {
	Keywords []string `json:"keywords"`
	Category string   `json:"category"`
	Size     string   `json:"size"`
	Color    string   `json:"color"`
}

// Task is a checkout account and an item(s) to checkout
type Task struct {
	TaskName string   `json:"taskName"`
	Item     taskItem `json:"item"`
	Account  Account  `json:"account"`
	API      string   `json:"api"`
	status   string
	id       string
	log      *zerolog.Logger
}

// ImportTasksFromJSON imports a list of tasks from a json file
func ImportTasksFromJSON(filename string) ([]Task, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err := json.Unmarshal(fileBytes, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// UpdateStatus sets the task status to the status provided
func (task *Task) UpdateStatus(status string) {
	task.status = status
}

// Log returns the logger associated with the task
func (task *Task) Log() *zerolog.Logger {
	if task.log == nil {
		// Logger for testing
		if log == nil {
			testLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
			task.log = &testLogger
		} else {
			// If task wasn't provided a logger during runtime
			tempLogger := log.With().Str("taskID", task.id).Logger()
			task.log = &tempLogger
		}
	}
	return task.log
}

// SetLog sets the task's logger
func (task *Task) SetLog(newLogger *zerolog.Logger) {
	task.log = newLogger
}

// VerifyTask verifies the information provided in the task to make sure it is
// what the rest of the application expects
func (task *Task) VerifyTask() (bool, error) {
	// Task category
	_, ok := supremeCategoriesDesktop[task.Item.Category]
	_, ok2 := supremeCategoriesMobile[task.Item.Category]
	if !ok && !ok2 {
		return false, errors.New("Task category not found")
	}
	// Task keywords
	if len(task.Item.Keywords) == 0 {
		return false, errors.New("Task keywords were not provided")
	}

	// Email
	if task.Account.Person.Email == "" {
		return false, errors.New("Email address field was empty")
	}

	// Phone number
	phoneMatch, _ := regexp.MatchString(`\d{3}-\d{3}-\d{4}`, task.Account.Person.PhoneNumber)
	if !phoneMatch || len(task.Account.Person.PhoneNumber) != 12 {
		return false, errors.New("Phone number was not correct")
	}

	// Credit card numbers
	ccFour, _ := regexp.MatchString(`\d{4} \d{4} \d{4} \d{4}`, task.Account.Card.Number)
	ccAmex, _ := regexp.MatchString(`\d{4} \d{6} \d{5}`, task.Account.Card.Number)
	if !(ccFour || ccAmex) {
		return false, errors.New("Credit card number was not correct")

	}

	// CVV
	ccvMatch := false
	if ccFour && len(task.Account.Card.Cvv) == 3 {
		ccvMatch, _ = regexp.MatchString(`\d{3}`, task.Account.Card.Cvv)
	} else if ccAmex && len(task.Account.Card.Cvv) == 4 {
		ccvMatch, _ = regexp.MatchString(`\d{4}`, task.Account.Card.Cvv)
	}
	if !ccvMatch {
		return false, errors.New("CVV was not correct")
	}

	// Month
	monthMatch, _ := regexp.MatchString(`\d{2}`, task.Account.Card.Month)
	if !monthMatch || len(task.Account.Card.Month) != 2 {
		return false, errors.New("Month was not correct")
	}

	// Year
	yearMatch, _ := regexp.MatchString(`\d{4}`, task.Account.Card.Year)
	if !yearMatch || len(task.Account.Card.Year) != 4 {
		return false, errors.New("Year was not correct")
	}

	// API
	if task.API != "" && strings.ToLower(task.API) != "desktop" && strings.ToLower(task.API) != "mobile" {
		return false, fmt.Errorf("API value %s was incorrect", task.API)
	}

	return true, nil
}

// VerifyTasks is a helper function that verifies multiple tasks and returns
// a slice containing the task number and the error
func VerifyTasks(tasks *[]Task) (bool, map[int]error) {
	allValid := true
	taskErrors := make(map[int]error)
	for i, task := range *tasks {
		valid, err := task.VerifyTask()
		if !valid {
			allValid = false
			taskErrors[i] = err
		}
	}
	if !allValid {
		return false, taskErrors
	}
	return true, nil
}

// SupremeCheckoutDesktop attempts to add to cart, waiting until it is available, and item and then check it out, on the desktop API
func (task *Task) SupremeCheckoutDesktop() (bool, error) {
	var matchedItem SupremeItem // The item on the supreme site we will buy
	var err error
	session := grequests.NewSession(nil)
	task.Log().Debug().
		Str("item", fmt.Sprintf("%+v", task.Item)).
		Msg("Checking out item")

	// Try to find the item provided in keywords etc
	task.UpdateStatus("Looking for item")
	for {
		// Get items in category
		matchedItem, err = waitForItemMatchDesktop(session, task)
		if err != nil {
			task.Log().Error().Err(err).
				Msgf("Error getting collection, sleeping.")
		} else {
			break
		}
		time.Sleep(time.Duration(appSettings.RefreshWait) * time.Millisecond)
	}
	task.UpdateStatus("Found item")
	task.Log().Debug().Msgf("Found item %+v %s %s %s", matchedItem, matchedItem.color, matchedItem.name, matchedItem.url)
	startTime := time.Now()

	// Get the ATC info from the item page
	var st string
	var sizeResponse SizeResponse
	var addURL string
	var xcsrf string
	task.UpdateStatus("Going to item page")
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Getting item info attempt: %d", attempt)
		var err error
		st, sizeResponse, addURL, xcsrf, err = GetSizeInfo(session, task, matchedItem.url)
		return err
	})
	if err != nil {
		task.Log().Error().Err(err)
		return false, err
	}
	task.Log().Debug().Msgf("%s (%s:%+v) %s %s", st, sizeResponse.singleSizeID, sizeResponse.multipleSizes, addURL, xcsrf)
	time.Sleep(time.Duration(appSettings.AtcWait) * time.Millisecond)

	// Add the item to cart
	task.UpdateStatus("Adding item to cart")
	pickedSizeID, err := PickSize(&task.Item, sizeResponse)
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
	time.Sleep(time.Duration(appSettings.CheckoutWait) * time.Millisecond)

	// Checkout
	task.UpdateStatus("Checking out")
	task.Log().Debug().Msgf("Checking out task: %s %s", task.Account.Person, task.Account.Address)
	var checkoutSuccess bool
	err = retry(10, 10*time.Millisecond, func(attempt int) error {
		task.Log().Debug().Msgf("Checkout attempt: %d", attempt)
		var err error
		checkoutSuccess, err = Checkout(session, task, xcsrf)
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

	return checkoutSuccess, nil
}

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
		task.Log().Error().Msg("Unable to find style")
		return false, errors.New("Unable to find style")
	}
	task.UpdateStatus("Matched style")
	task.Log().Debug().Msgf("Matched Style: %+v", matchedStyle)

	pickedSizeID, err := PickSizeMobile(&task.Item, matchedStyle)
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

	// Purecart implementation?
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
