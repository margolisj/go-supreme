package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFindItem tests finding a single item
func TestFindItem(t *testing.T) {
	taskItem := taskItem{
		[]string{"hanes", "boxer"},
		"accessories",
		"Medium",
		"white",
	}

	supremeItems := SupremeItems{SupremeItem{
		"Supreme®/Hanes® Boxer Briefs (4 Pack)",
		"White",
		"shop/accessories/nckme38ul/iimyp2ogd",
	}}

	_, err := FindItem(taskItem, supremeItems)

	if err != nil {
		t.Fail()
	}
}

//TODO: HTML source

// TestDesktopCheckoutResponsesUnmarshall tests the json unmarshalling of the standard checkout responses from desktop
func TestDesktopCheckoutResponsesUnmarshall(t *testing.T) {
	queuedStatus := []byte(`{"status":"queued","slug":"q7j84cuad93wnyrg0"}`)
	var res checkoutJSON
	var err error

	if err = json.Unmarshal(queuedStatus, &res); err != nil {
		panic(err)
	}
	assert.Equal(t, "queued", res.Status)
	assert.Equal(t, "q7j84cuad93wnyrg0", res.Slug)

	// failedCreditCard := []byte(`{"status":"failed","cart":[{"size_id":"59765","in_stock":true}],"errors":{"order":"","credit_card":"number is not a valid credit card number"}}`)
	// if err = json.Unmarshal(failedCreditCard, &res); err != nil {
	// 	panic(err)
	// }
	// assert.Equal(t, "failed", res.Status)
	// assert.Equal(t, `{"order":"","credit_card":"number is not a valid credit card number"}`, res.Slug)
}
