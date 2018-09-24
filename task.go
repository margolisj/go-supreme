package main

import (
	"encoding/json"
	"io/ioutil"
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
	Cardtype string `json:"cardtype"`
	Number   string `json:"number"` // Card numbers must have spaces "XXXX XXXX XXXX XXXX"
	Month    string `json:"month"`  // Two digit month, ex. 03
	Year     string `json:"year"`   // 4 digit year, ex. 2019
	Cvv      string `json:"cvv"`
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

// type proxy struct {
// 	ip       string
// 	port     string
// 	username string
// 	password string
// }

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
