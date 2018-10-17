package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskUnmarshal(t *testing.T) {
	s := []byte(`{
		"taskName": "Task1",
		"item": {
			"keywords": [
				"shaolin"
			],
			"category": "hats",
			"size": "",
			"color": "orange"
		},
		"account": {
			"person": {
				"firstname": "Jax",
				"lastname": "Blax",
				"email": "none@none.com",
				"phoneNumber": "215-834-1857"
			},
			"address": {
				"address1": "102 Broad Street",
				"address2": "",
				"zipcode": "12345",
				"city": "Philadeliphia",
				"state": "PA",
				"country": "USA"
			},
			"card": {
				"cardtype": "visa",
				"number": "1285 4827 5948 2017",
				"month": "02",
				"year": "2019",
				"cvv": "847"
			}
		},
		"api": "desktop",
		"waitSettings": {
			"RefreshWait": 1000,
			"AtcWait": 900,
			"CheckoutWait": 800
		}
	}`)
	var task Task
	if err := json.Unmarshal(s, &task); err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Task1", task.TaskName)
	assert.Equal(t, taskItem{
		Keywords: []string{"shaolin"},
		Category: "hats",
		Size:     "",
		Color:    "orange",
	}, task.Item)
	assert.Equal(t, Person{
		Firstname:   "Jax",
		Lastname:    "Blax",
		Email:       "none@none.com",
		PhoneNumber: "215-834-1857",
	}, task.Account.Person)
	assert.Equal(t, Address{
		Address1: "102 Broad Street",
		Address2: "",
		Zipcode:  "12345",
		City:     "Philadeliphia",
		State:    "PA",
		Country:  "USA",
	}, task.Account.Address)
	assert.Equal(t, Card{
		Cardtype: "visa",
		Number:   "1285 4827 5948 2017",
		Month:    "02",
		Year:     "2019",
		Cvv:      "847",
	}, task.Account.Card)
	assert.Equal(t, "desktop", task.API)
	assert.Equal(t, 1000, task.WaitSettings.RefreshWait)
	assert.Equal(t, 900, task.WaitSettings.AtcWait)
	assert.Equal(t, 800, task.WaitSettings.CheckoutWait)
}

func TestReadImportTasksFromJSONFile(t *testing.T) {
	tasks, err := ImportTasksFromJSON("testdata/validSingleTask.json")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(tasks))
	task := tasks[0]
	assert.Equal(t, "shopsafe 0 0RU3", task.TaskName)
}

func TestVerifyTaskValid(t *testing.T) {
	task := testTask()
	valid, err := task.VerifyTask()
	if err != nil {
		t.Log(err)
	}
	assert.True(t, valid)
}

func TestVerifyTaskValidAmex(t *testing.T) {
	task := testTask()
	task.Account.Card.Number = "1234 567890 12345"
	task.Account.Card.Cvv = "1234"
	valid, err := task.VerifyTask()
	if err != nil {
		t.Error(err)
	}
	assert.True(t, valid)
}

func TestVeriifyTaskMissingItem(t *testing.T) {
	task := testTask()
	task.Item = taskItem{}
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("Task category not found"), err)
	assert.False(t, valid)
}

func TestVeriifyTaskIncorrectCategory(t *testing.T) {
	task := testTask()
	task.Item = taskItem{}
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("Task category not found"), err)
	assert.False(t, valid)
}

func TestVerifyTaskMissingKeywords(t *testing.T) {
	task := testTask()
	task.Item.Keywords = []string{}
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("Task keywords were not provided"), err)
	assert.False(t, valid)
}

func TestVerifyTaskBadPhoneNumber(t *testing.T) {
	task := testTask()
	// No dashes
	task.Account.Person.PhoneNumber = "954 876 6543"
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("Phone number was not correct"), err)
	assert.False(t, valid)

	// Missing number
	task.Account.Person.PhoneNumber = "954-876-154"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("Phone number was not correct"), err)
	assert.False(t, valid)

	// Additional number
	task.Account.Person.PhoneNumber = "954-876-15434"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("Phone number was not correct"), err)
	assert.False(t, valid)
}

func TestVerifyTaskBadCCNumber(t *testing.T) {
	task := testTask()
	// Missing number
	task.Account.Card.Number = "1234 1548 1548 125"
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("Credit card number was not correct"), err)
	assert.False(t, valid)

	// Dashes
	task.Account.Card.Number = "1234-1548-1548-1548"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("Credit card number was not correct"), err)
	assert.False(t, valid)
}

func TestVerifyTaskBadCvv(t *testing.T) {
	task := testTask()
	// ccFour missing number
	task.Account.Card.Cvv = "12"
	valid, err := task.VerifyTask()
	assert.Error(t, errors.New("CVV was not correct"), err)
	assert.False(t, valid)

	// ccFour additional number
	task.Account.Card.Cvv = "1234"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("CVV was not correct"), err)
	assert.False(t, valid)

	task.Account.Card.Number = "1234 567890 12345"
	// amex missing number
	task.Account.Card.Cvv = "123"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("CVV was not correct"), err)
	assert.False(t, valid)

	// amex additional number
	task.Account.Card.Cvv = "12345"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("CVV was not correct"), err)
	assert.False(t, valid)
}

func TestVerifyTaskBadMonth(t *testing.T) {
	task := testTask()
	// Missing number
	task.Account.Card.Month = "4"
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("Month was not correct"), err)
	assert.False(t, valid)

	// Additional number
	task.Account.Card.Month = "122"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("Month was not correct"), err)
	assert.False(t, valid)
}

func TestVerifyTaskBadYear(t *testing.T) {
	task := testTask()
	// Missing number
	task.Account.Card.Year = "21"
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("Year was not correct"), err)
	assert.False(t, valid)

	// Additional number
	task.Account.Card.Year = "20199"
	valid, err = task.VerifyTask()
	assert.Equal(t, errors.New("Year was not correct"), err)
	assert.False(t, valid)
}

func TestVerifyTaskBadAPI(t *testing.T) {
	task := testTask()
	task.API = "asdf"
	valid, err := task.VerifyTask()
	assert.Equal(t, fmt.Errorf("API value %s was incorrect", task.API), err)
	assert.False(t, valid)
}

func TestVerifyTaskValidNewCategory(t *testing.T) {
	task := testTask()
	task.API = "mobile"
	task.Item.Category = "new"
	valid, err := task.VerifyTask()
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestVerifyTaskBadNewCategory(t *testing.T) {
	task := testTask()
	task.API = "desktop"
	task.Item.Category = "new"
	valid, err := task.VerifyTask()
	assert.Equal(t, errors.New("new category can only be used with mobile API"), err)
	assert.False(t, valid)
}

func TestVerifyTasks(t *testing.T) {
	tasks := []Task{testTask(), testTask(), testTask()}
	valid, errs := VerifyTasks(&tasks)
	assert.True(t, valid)
	assert.Nil(t, errs)
}

func TestVertifyTasksBad(t *testing.T) {
	tasks := []Task{testTask(), testTask(), testTask()}
	tasks[2].Account.Card.Number = "1234-5849-2894-6753"
	valid, errs := VerifyTasks(&tasks)
	assert.False(t, valid)
	assert.Equal(t, map[int]error{
		2: errors.New("Credit card number was not correct"),
	}, errs)
}

func TestGetRates(t *testing.T) {
	task := testTask()
	// These should be equal because the test task is missing waitSettings values
	assert.Equal(t, appSettings.RefreshWait, task.GetTaskRefreshRate())
	assert.Equal(t, appSettings.AtcWait, task.GetTaskAtcWait())
	assert.Equal(t, appSettings.CheckoutWait, task.GetTaskCheckoutWait())

	task.WaitSettings = WaitSettings{
		RefreshWait: 1000,
		AtcWait:     900,
	}
	assert.Equal(t, 1000, task.GetTaskRefreshRate())
	assert.Equal(t, 900, task.GetTaskAtcWait())

	task.WaitSettings = WaitSettings{
		RefreshWait:  343,
		AtcWait:      0,
		CheckoutWait: 800,
	}
	assert.Equal(t, 343, task.GetTaskRefreshRate())
	// This should be appSettings.AtcWait because the value of AtcWait is 0
	assert.Equal(t, appSettings.AtcWait, task.GetTaskAtcWait())
	assert.Equal(t, 800, task.GetTaskCheckoutWait())
}

// func TestTaskSupremeCheckoutMobile(t *testing.T) {
// 	task := Task{
// 		TaskName: "Task1",
// 		Item: taskItem{
// 			[]string{"Briefs"},
// 			"accessories",
// 			"medium",
// 			"white",
// 		},
// 		Account: testAccount(),
// 	}
// 	// task := Task{
// 	// 	TaskName: "Task1",
// 	// 	Item: taskItem{
// 	// 		[]string{"gold"},
// 	// 		"accessories",
// 	// 		"",
// 	// 		"gold",
// 	// 	},
// 	// 	Account: testAccount(),
// 	// }
// 	success, err := task.SupremeCheckoutMobile()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(success)
// }
