// +build integration

package main

import (
	"testing"

	"net/http/cookiejar"

	"github.com/levigross/grequests"
	"golang.org/x/net/publicsuffix"
)

func buildTestSession(t *testing.T) *grequests.Session {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Error(err)
	}

	session := grequests.NewSession(&grequests.RequestOptions{
		CookieJar: jar,
	})
	return session
}

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
