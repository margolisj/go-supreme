package main

type person struct {
	firstname   string
	lastname    string
	email       string
	phoneNumber string
}

type card struct {
	cardtype string
	number   string
	month    string
	year     string
	cvv      string
}

type address struct {
	address1 string
	address2 string
	zipcode  string
	city     string
	state    string
	country  string
}

type account struct {
	person  person
	address address
	card    card
}

type taskItem struct {
	keywords []string
	category string
	size     string
	color    string
}

// type proxy struct {
// 	ip       string
// 	port     string
// 	username string
// 	password string
// }

type task struct {
	taskName string
	// proxy    proxy
	items   []taskItem
	success bool
	status  string
}

// ImportTasksFromJSON imports a list of tasks from a json file
func ImportTasksFromJSON(filename string) {

}
