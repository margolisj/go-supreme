// +build integration

package main

import (
	"fmt"
	"testing"
)

func TestGetCollectionItemsMobile(t *testing.T) {
	session := buildTestSession(t)

	item := TaskItem{
		[]string{"temp"},
		"accessories",
		"",
		"blue",
	}
	task := testTask()
	task.Item = item
	items, err := GetCollectionItemsMobile(session, &task)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", items)

}
