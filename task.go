package main

import (
	"encoding/json"
	"io/ioutil"
)

type person struct {
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"` // Phone numbers must have dashes
}

type card struct {
	Cardtype string `json:"cardtype"`
	Number   string `json:"number"` // Card numbers must have spaces "XXXX XXXX XXXX XXXX"
	Month    string `json:"month"`  // Two digit month, ex. 03
	Year     string `json:"year"`   // 4 digit year, ex. 2019
	Cvv      string `json:"cvv"`
}

type address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Zipcode  string `json:"zipcode"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

type account struct {
	Person  person  `json:"person"`
	Address address `json:"address"`
	Card    card    `json:"card"`
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

// Task is a combination of items, status
type Task struct {
	TaskName string `json:"taskName"`
	// proxy    proxy
	Item taskItem `json:"item"`
	// Success bool
	// status  string
}

// FullTask is an account and task item
type FullTask struct {
	Task    Task    `json:"task"`
	Account account `json:"account"`
}

// ImportTasksFromJSON imports a list of tasks from a json file
func ImportTasksFromJSON(filename string) ([]FullTask, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tasks []FullTask
	if err := json.Unmarshal(fileBytes, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
