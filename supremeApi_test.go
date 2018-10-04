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

func TestParseSizesMultipleSizesPickSizes(t *testing.T) {
	f, err := os.Open("./testData/supremeSite/10-3-18-whiteTagless.html")
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
		"Small":  "59754",
		"Medium": "59755",
		"Large":  "59756",
		"XLarge": "59757",
	}, *sizeResponse.multipleSizes)

	pickTests := []struct {
		name string
		in   taskItem
		out  string
	}{
		{"small", taskItem{Size: "small"}, "59754"},
		{"medium", taskItem{Size: "Medium"}, "59755"},
		{"large", taskItem{Size: "LaRGe"}, "59756"},
		{"xlarge", taskItem{Size: "xLarge"}, "59757"},
	}

	for _, tt := range pickTests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := PickSize(&tt.in, sizeResponse)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.out, size)
		})
	}
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

	// Currently the checkout JSON is pretty useless because supreme doesn't follow a real protocal
	failedCreditCard := []byte(`{"status":"failed","cart":[{"size_id":"59765","in_stock":true}],"errors":{"order":"","credit_card":"number is not a valid credit card number"}}`)
	if err = json.Unmarshal(failedCreditCard, &res); err != nil {
		assert.Error(t, err)
	}
}
