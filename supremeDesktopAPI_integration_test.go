// +build integration

package main

import (
	"testing"
)

func TestGetCollectionItems(t *testing.T) {
	session := buildTestSession(t)

	item := taskItem{
		[]string{"temp"},
		"accessories",
		"",
		"blue",
	}
	task := testTask()
	task.Item = item
	GetCollectionItems(session, &task, false)

}
