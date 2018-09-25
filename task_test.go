package main

import (
	"encoding/json"
	"testing"
)

// TestFullTaskUnmarhsal tests basic unmarhsalling of the
func TestFullTaskUnmarhsal(t *testing.T) {
	s := []byte(`{"task": {"taskName": "task1", "item": {"keywords": ["group"], "size": "Medium", "color": "white", "category": "t-shirts"}}, "account": {"person": {"firstname": "Jack", "lastname": "Black", "email": "bob@bob.net", "phoneNumber": "610-659-8745"}, "address": {"address1": "654 E Hockey Rd", "address2": "", "zipcode": "25148", "city": "Jersey City", "state": "NJ", "country": "USA"}, "card": {"cardtype": "visa", "number": "1254 0987 6541 8375", "month": "04", "year": "2020", "cvv": "845"}}}`)
	var tas FullTask
	if err := json.Unmarshal(s, &tas); err != nil {
		panic(err)
	}
}

// TestReadTasks tests reading
func TestReadTasks(t *testing.T) {
	ImportTasksFromJSON("testdata/testFullTaskCorrect.json")
}

// func TestReadTasksBrokenJSON(t *testing.T) {
// 	ImportTasksFromJSON("testFile.txt")
// }
