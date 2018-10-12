package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskUnmarhsal(t *testing.T) {
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
		"api": "desktop"
	}`)
	var tas Task
	if err := json.Unmarshal(s, &tas); err != nil {
		t.Error(err)
	}
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
