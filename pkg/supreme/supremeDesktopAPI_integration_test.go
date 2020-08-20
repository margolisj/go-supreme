// +build integration

package supreme

import (
	"testing"
)

func TestGetCollectionItems(t *testing.T) {
	session := buildTestSession(t)

	item := TaskItem{
		[]string{"temp"},
		"accessories",
		"",
		"blue",
	}
	task := testTask()
	task.Item = item
	GetCollectionItems(session, &task, false)

}
