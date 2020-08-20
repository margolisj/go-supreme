// +build keywordsUK

package main

import (
	"net/url"
	"testing"

	"github.com/levigross/grequests"
	"github.com/margolisj/go-supreme/task"
)

func TestUKKeywords(t *testing.T) {
	items := []task.TaskItem{
		task.TaskItem{
			Keywords: []string{
				"sanders",
			},
			Size:     "medium",
			Color:    "white",
			Category: "t-shirts",
		},
		// task.TaskItem{
		// 	Keywords: []string{
		// 		"trust",
		// 		"killer",
		// 	},
		// 	Size:     "medium",
		// 	Color:    "heather grey",
		// 	Category: "t-shirts",
		// },
		// task.TaskItem{
		// 	Keywords: []string{
		// 		"L/S",
		// 		"killer",
		// 	},
		// 	Size:     "medium",
		// 	Color:    "white",
		// 	Category: "t-shirts",
		// },
	}

	task := testTask()
	proxyString := "81.130.135.142:48057"
	proxyURL, err := url.Parse("http://" + proxyString) // Proxy URL
	if err != nil {
		t.Error(err)
	}

	localRo := &grequests.RequestOptions{
		UserAgent: supremeAPIMobile.mobileUserAgent,
		Proxies: map[string]*url.URL{
			"http":  proxyURL,
			"https": proxyURL,
		},
	}

	session := *grequests.NewSession(localRo)
	tresp, _ := session.Get("https://api.ipify.org?format=json", nil)
	t.Logf("Current IP: %s", tresp.String())

	for k, item := range items {
		task.Item = item
		itemMobile, err := waitForItemMatchMobile(&session, &task)
		if err != nil {
			t.Logf("Item %d keyword couldn't be found", k)
		} else {
			t.Logf("Item %d keyword was found: %+v", k, itemMobile)
		}

		style, err := waitForStyleMatchMobile(&session, &task, itemMobile)
		if err != nil {
			t.Logf("Item %d couldn't be found style / color", k)
		} else {
			t.Logf("Item %d was found style / color: %+v", k, style)
		}

		ID, _, err := PickSizeMobile(&task.Item, style)
		if err != nil {
			t.Logf("Item %d couldn't be found size", k)
		} else {
			t.Logf("Item %d was found size: %d", k, ID)
		}

	}
}
