package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"regexp"
)

// Person is a struct modeling personal information
type Person struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"` // Phone numbers must have dashes
}

// Card is a struct moding a credit card
type Card struct {
	Cardtype string `json:"cardtype"` // Don't think this matters for desktop, also should be able to figure this out without user entry
	Number   string `json:"number"`   // Card numbers must have spaces "XXXX XXXX XXXX XXXX"
	Month    string `json:"month"`    // Two digit month, ex. 03
	Year     string `json:"year"`     // 4 digit year, ex. 2019
	Cvv      string `json:"cvv"`      // 3 digit or 4 digit
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

type proxy struct {
	ip       string
	port     string
	username string
	password string
}

// Task is a checkout acount and an item(s) to checkout
type Task struct {
	TaskName string `json:"taskName"`
	// proxy    proxy
	Item taskItem `json:"item"`
	// Success bool
	// status  string
	Account Account `json:"account"`
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

// VerifyTask verifies the information provided in the task to make sure it is
// what the rest of the applciation expects
func VerifyTask(task *Task) (bool, error) {
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

	return true, nil
}

// VerifyTasks is a helper function that verifies multiple tasks and returns
// a slice containing the task number and the error
func VerifyTasks(tasks *[]Task) (bool, map[int]error) {
	allValid := true
	taskErrors := make(map[int]error)
	for i, task := range *tasks {
		valid, err := VerifyTask(&task)
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
