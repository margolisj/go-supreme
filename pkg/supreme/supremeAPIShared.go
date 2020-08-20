package supreme

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/levigross/grequests"
)

var supremeCategoriesDesktop = map[string]string{
	"jackets":       "jackets",
	"shirts":        "shirts",
	"tops/sweaters": "tops_sweaters",
	"sweatshirts":   "sweatshirts",
	"pants":         "pants",
	"t-shirts":      "t-shirts",
	"hats":          "hats",
	"bags":          "bags",
	"shorts":        "shorts",
	"accessories":   "accessories",
	"skate":         "skate",
	"shoes":         "shoes",
}

var supremeCategoriesMobile = map[string]string{
	"jackets":       "Jackets",
	"shirts":        "Shirts",
	"tops/sweaters": "Tops/Sweaters",
	"sweatshirts":   "Sweatshirts",
	"pants":         "Pants",
	"t-shirts":      "T-Shirts",
	"hats":          "Hats",
	"bags":          "Bags",
	"shorts":        "Shorts",
	"accessories":   "Accessories",
	"skate":         "Skate",
	"shoes":         "Shoes",
	"new":           "new",
}

func checkKeywords(keywords []string, supremeItemName string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(strings.ToLower(supremeItemName), strings.ToLower(keyword)) {
			// fmt.Printf("%s doesn't contain %s\n", supremeItemName, keyword)
			return false
		}
	}
	return true
}

func checkColor(taskItemColor string, supremeItemColor string) bool {
	if taskItemColor == "" {
		return true
	}
	return strings.Contains(strings.ToLower(strings.TrimSpace(supremeItemColor)), strings.ToLower(taskItemColor))
}

func handleCheckoutResponse(task *Task, respString *string) bool {
	if strings.Contains(*respString, "queued") {
		task.Log().Info().
			Str("response", *respString).
			Msg("Checkout successful.")
		return true
	} else if strings.Contains(*respString, "failed") {
		task.Log().Error().
			Str("reason", "failed").
			Str("response", *respString).
			Msg("Checkout failed.")
		return false
	} else if strings.Contains(*respString, "outOfStock") {
		task.Log().Error().
			Str("reason", "outOfStock").
			Str("response", *respString).
			Msg("Checkout failed.")
		return false
	} else if strings.Contains(*respString, "declined") {
		task.Log().Error().
			Str("reason", "declined").
			Str("response", *respString).
			Msg("Checkout failed.")
		return false
	} else if strings.Contains(*respString, "status\":\"dup") {
		task.Log().Error().
			Str("reason", "dup").
			Str("response", *respString).
			Msg("Checkout failed.")
		return false
	} else if strings.Contains(*respString, "status\":\"dup") {
		task.Log().Error().
			Str("reason", "dup").
			Str("response", *respString).
			Msg("Checkout failed.")
		return true
	}

	task.Log().Error().
		Str("response", *respString).
		Msg("Error processing checkout response.")
	return false
}

// checkoutJSON the json response provided after check out.
// This does not capture all the possible checkout response only
// the response if we need to queue
type checkoutJSON struct {
	Status string `json:"status"`
	Slug   string `json:"slug"`
	Errors string `json:"errors"`
}

// Queue handles the queue response and finished the checkout
func Queue(session *grequests.Session, task *Task, originalRespString string) (bool, error) {
	var queueJSON checkoutJSON
	if err := json.Unmarshal([]byte(originalRespString), &queueJSON); err != nil {
		task.Log().Error().
			Str("originalResponse", originalRespString).
			Msg("Unable to marshall json in queue")
		return false, nil
	}

	task.Log().Debug().Msgf("%+v", queueJSON)
	task.Log().Info().Msg("Sleeping 10 seconds in queue")
	time.Sleep(10 * time.Second)

	localRo := &grequests.RequestOptions{
		UserAgent: sharedUserAgent,
		Headers: map[string]string{
			"accept":           "application/json",
			"accept-encoding":  "gzip, deflate, br",
			"accept-language":  "en-US,en;q=0.9",
			"origin":           "https://www.supremenewyork.com",
			"referer":          "https://www.supremenewyork.com/checkout",
			"x-requested-with": "XMLHttpRequest",
		},
	}

	queueURL := fmt.Sprintf("https://www.supremenewyork.com/checkout/%s/status.json", queueJSON.Slug)
	task.Log().Debug().Str("queueURL", queueURL)
	resp, err := session.Get(queueURL, localRo)
	if err != nil {
		task.Log().Error().Err(err).Msg("Queue error on update")
		return false, err
	}

	task.Log().Debug().Msgf("%+v", resp.RawResponse.Request)
	task.Log().Debug().Msgf("%+v", resp.RawResponse)
	if resp.Ok != true {
		task.Log().Warn().Msg("Queue response did not return OK")
		return false, errors.New("Queue response did not return OK")
	}

	respString := resp.String()
	task.Log().Debug().Msgf("Queue Response: %s", respString)

	// Process queue response
	isStillInQueue, queueSuccess := handleQueueResponse(task, &respString)
	if queueSuccess {
		return true, nil
	}
	if isStillInQueue {
		return Queue(session, task, originalRespString)
	}

	return queueSuccess, err
}

// handleQueueResponse returns two bools, the first being if we are still in the queue, the second
// is if the queue was successful
func handleQueueResponse(task *Task, respString *string) (bool, bool) {
	if strings.Contains(*respString, "queued") {
		task.Log().Debug().
			Str("response", *respString).
			Msg("Still in queue.")
		return true, false
	} else if strings.Contains(*respString, "failed") {
		task.Log().Error().
			Str("reason", "failed").
			Str("response", *respString).
			Msg("Queue failed.")
		return false, false
	} else if strings.Contains(*respString, "outOfStock") {
		task.Log().Error().
			Str("reason", "outOfStock").
			Str("response", *respString).
			Msg("Queue failed.")
		return false, false
	} else if strings.Contains(*respString, "paid") {
		task.Log().Info().
			Str("status", "paid").
			Str("response", *respString).
			Msg("Queue successful.")
		return false, true
	}

	task.Log().Error().
		Str("response", *respString).
		Msg("Error processing queue response.")
	return false, false
}
