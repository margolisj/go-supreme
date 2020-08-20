// +build unit

package supreme

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/stretchr/testify/assert"
)

func TestFindSingleItem(t *testing.T) {
	taskItem := TaskItem{
		Keywords: []string{"hanes", "boxer"},
		Category: "accessories",
		Size:     "Medium",
		Color:    "white",
	}
	targetItem := SupremeItem{
		"Supreme®/Hanes® Boxer Briefs (4 Pack)",
		"White",
		"shop/accessories/nckme38ul/iimyp2ogd",
	}
	supremeItems := []SupremeItem{targetItem}

	foundItem, err := findItem(taskItem, supremeItems)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, targetItem, foundItem)
}

func TestFindSingleItemFromPageSource(t *testing.T) {
	f, err := os.Open("./testData/supremeSite/3-3-2020-desktop-accessories.html")
	if err != nil {
		t.Error(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Error(err)
	}

	supremeItems := parseCategoryPage(doc, false)
	assert.Equal(t, 17, len(*supremeItems))

	taskItem := TaskItem{
		Keywords: []string{"Plate"},
		Category: "accessories",
		Size:     "",
		Color:    "Gold",
	}

	foundItem, err := findItem(taskItem, *supremeItems)

	if err != nil {
		t.Fail()
	}
	assert.Equal(t, SupremeItem{
		name:  "Name Plate 14K Gold Pendant",
		color: "Gold",
		url:   "/shop/accessories/zzxnomh1w/ue59f7gvh",
	}, foundItem)
}

// TODO: Replace this with an updated empty desktop category
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
		{"small", TaskItem{Size: "small"}, "59754"},
		{"medium", TaskItem{Size: "Medium"}, "59755"},
		{"large", TaskItem{Size: "LaRGe"}, "59756"},
		{"xlarge", TaskItem{Size: "xLarge"}, "59757"},
	}

	for _, tt := range pickTests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := PickSize(&tt.in, &sizeResponse)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.out, size)
		})
	}
}

func TestPickSizeNoSize(t *testing.T) {
	item := TaskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "",
		Color:    "blue",
	}

	sizeResponse := SizeResponse{"60885", nil}

	itemID, err := PickSize(&item, &sizeResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "60885", itemID)
}

func TestPickSizeMultipleSizes(t *testing.T) {
	item := TaskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "Medium",
		Color:    "blue",
	}

	sizeResponse := SizeResponse{"", &map[string]string{
		"Small":  "59758",
		"Medium": "59759",
		"Large":  "59760",
	}}

	itemID, err := PickSize(&item, &sizeResponse)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "59759", itemID)
}

func TestPickSizeSingleSizeReturned(t *testing.T) {
	item := TaskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "Medium",
		Color:    "blue",
	}

	sizeResponse := SizeResponse{"60885", nil}

	itemID, err := PickSize(&item, &sizeResponse)
	assert.Error(t, err, "Unable to pick size")
	assert.Equal(t, "", itemID)
}

func TestPickSizeMultipleSizesReturnedButNoneInTask(t *testing.T) {
	item := TaskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "",
		Color:    "blue",
	}

	sizeResponse := SizeResponse{"", &map[string]string{
		"Small":  "59758",
		"Medium": "59759",
		"Large":  "59760",
	}}

	itemID, err := PickSize(&item, &sizeResponse)
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
