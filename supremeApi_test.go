package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/stretchr/testify/assert"
)

func TestFindSingleItem(t *testing.T) {
	taskItem := taskItem{
		[]string{"hanes", "boxer"},
		"accessories",
		"Medium",
		"white",
	}
	targetItem := SupremeItem{
		"Supreme®/Hanes® Boxer Briefs (4 Pack)",
		"White",
		"shop/accessories/nckme38ul/iimyp2ogd",
	}
	supremeItems := SupremeItems{targetItem}

	foundItem, err := findItem(taskItem, supremeItems)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, targetItem, foundItem)
}

func TestFindSingleItemFromPageSource(t *testing.T) {
	f, err := os.Open("./testData/supremeSite/9-25-18-jackets.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Error(err)
	}

	supremeItems := parseCategoryPage(doc, true)

	taskItem := taskItem{
		[]string{"bone"},
		"jackets",
		"Medium",
		"Black",
	}

	foundItem, err := findItem(taskItem, *supremeItems)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t, SupremeItem{
		name:  "Bone Varsity Jacket",
		color: "Black",
		url:   "/shop/jackets/m2ihxzpus/wq798ar2h",
	}, foundItem)
}

func TestParseCategoryPage(t *testing.T) {
	f, err := os.Open("./testData/supremeSite/9-25-18-jackets.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Error(err)
	}

	items := parseCategoryPage(doc, true)
	assert.Equal(t, 10, len(*items))

	item := (*items)[0]
	assert.Equal(t, SupremeItem{
		"Bone Varsity Jacket",
		"Black",
		"/shop/jackets/m2ihxzpus/wq798ar2h",
	}, item)
}
func TestParseEmptyCategoryPage(t *testing.T) {
	f, err := os.Open("./testData/supremeSite/9-29-18-jackets-empty.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Error(err)
	}

	items := parseCategoryPage(doc, true)
	assert.Equal(t, 0, len(*items))
}

func TestParseSizesSingleSize(t *testing.T) {
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

	itemID, err := PickSize(&item, sizeResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "60885", itemID)
}

func TestPickSizeMultipleSizes(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"Medium",
		"blue",
	}

	sizeResponse := SizeResponse{"", &map[string]string{
		"Small":  "59758",
		"Medium": "59759",
		"Large":  "59760",
	}}

	itemID, err := PickSize(&item, sizeResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "59759", itemID)
}

func TestPickSizeSingleSizeReturned(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"Medium",
		"blue",
	}

	sizeResponse := SizeResponse{"60885", nil}

	itemID, err := PickSize(&item, sizeResponse)
	assert.Error(t, err, "Unable to pick size")
	assert.Equal(t, "", itemID)
}

func TestPickSizeMultipleSizesReturnedButNoneInTask(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"",
		"blue",
	}

	sizeResponse := SizeResponse{"", &map[string]string{
		"Small":  "59758",
		"Medium": "59759",
		"Large":  "59760",
	}}

	itemID, err := PickSize(&item, sizeResponse)
	assert.Error(t, err, "Unable to pick size")
	assert.Equal(t, "", itemID)
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
