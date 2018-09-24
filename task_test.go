package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testAccount() Account {
	p := Person{
		"Jax",
		"Blax",
		"none@none.com",
		"215-834-1857", // Phone numbers have dashes
	}

	a := Address{
		"102 Broad Street",
		"",
		"12345",
		"Philadeliphia",
		"PA",
		"USA",
	}

	c := Card{
		"visa",                // Don't think this matters for desktop
		"1285 4827 5948 2017", //TODO: These cards must have spaces
		"02",
		"2019", // 4 digit dates
		"847",  //
	}

	return Account{p, a, c}
}

func testTask() Task {
	item := taskItem{
		[]string{"shaolin"},
		"hats",
		"",
		"orange",
	}

	return Task{
		"Task1",
		item,
		testAccount(),
	}
}

func TestTaskMarshal(t *testing.T) {
	_, err := json.Marshal(testTask())
	if err != nil {
		t.Error("Unable to marshall task")
	}
}

// TestFullTaskUnmarhsal tests basic unmarhsalling of the
func TestFullTaskUnmarhsal(t *testing.T) {
	s := []byte(`{"taskName":"Task1","item":{"keywords":["shaolin"],"category":"hats","size":"","color":"orange"},"account":{"person":{"firstname":"Jax","lastname":"Blax","email":"none@none.com","phoneNumber":"215-834-1857"},"address":{"address1":"102 Broad Street","address2":"","zipcode":"12345","city":"Philadeliphia","state":"PA","country":"USA"},"card":{"cardtype":"visa","number":"1285 4827 5948 2017","month":"02","year":"2019","cvv":"847"}}}`)
	var tas Task
	if err := json.Unmarshal(s, &tas); err != nil {
		panic(err)
	}
}

// TestReadTasks tests reading
func TestReadTasks(t *testing.T) {
	tasks, err := ImportTasksFromJSON("testdata/validSingleTask.json")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(tasks))
	task := tasks[0]
	assert.Equal(t, "shopsafe 0 0RU3", task.TaskName)
}
