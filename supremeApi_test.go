package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

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

func TestParseSizesSingleSize(t *testing.T) {
	// dat, err := ioutil.ReadFile("/testData/supremeSite/9-24-18-lucettaLight.html")
	f, err := os.Open("./testData/supremeSite/9-24-18-lucettaLight.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Error(err)
	}

	sizeResponse := parseSizes(doc)
	assert.Equal(t, "60885", sizeResponse.singleSizeID)
}

func TestParseSizesMultipleSizes(t *testing.T) {
	// dat, err := ioutil.ReadFile("/testData/supremeSite/9-24-18-lucettaLight.html")
	f, err := os.Open("./testData/supremeSite/9-24-18-blackTagless.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Error(err)
	}

	sizeResponse := parseSizes(doc)
	assert.Equal(t, "", sizeResponse.singleSizeID)
	assert.Equal(t, map[string]string{
		"Small":  "59758",
		"Medium": "59759",
		"Large":  "59760",
	}, *sizeResponse.multipleSizes)
}

func TestPickSizeNoSize(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"", // No size
		"blue",
	}

	sizeResponse := SizeResponse{"60885", nil}

	itemID, err := PickSize(item, sizeResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "60885", itemID)
}

func TestPickSizeMultipleSizes(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"Medium", // No size
		"blue",
	}

	sizeResponse := SizeResponse{"", &map[string]string{
		"Small":  "59758",
		"Medium": "59759",
		"Large":  "59760",
	}}

	itemID, err := PickSize(item, sizeResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "59759", itemID)
}

func TestDesktopCheckoutResponsesUnmarshall(t *testing.T) {
	queuedStatus := []byte(`{"status":"queued","slug":"q7j84cuad93wnyrg0"}`)
	var res checkoutJSON
	var err error

	if err = json.Unmarshal(queuedStatus, &res); err != nil {
		t.Error(err)
	}
	assert.Equal(t, "queued", res.Status)
	assert.Equal(t, "q7j84cuad93wnyrg0", res.Slug)

	// failedCreditCard := []byte(`{"status":"failed","cart":[{"size_id":"59765","in_stock":true}],"errors":{"order":"","credit_card":"number is not a valid credit card number"}}`)
	// if err = json.Unmarshal(failedCreditCard, &res); err != nil {
	// 	t.Error(err)
	// }
	// assert.Equal(t, "failed", res.Status)
	// assert.Equal(t, `{"order":"","credit_card":"number is not a valid credit card number"}`, res.Errors)
}
