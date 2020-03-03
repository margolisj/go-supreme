// +build unit

package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSingleItemMobile(t *testing.T) {
	taskItem := taskItem{
		Keywords: []string{"hanes", "boxer"},
		Category: "accessories",
		Size:     "Medium",
		Color:    "white",
	}
	targetItem := SupremeItemMobile{
		"Supreme®/Hanes® Boxer Briefs (4 Pack)",
		171745,
	}
	supremeItems := []SupremeItemMobile{targetItem}

	foundItem, err := findItemMobile(taskItem, &supremeItems)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, targetItem, foundItem)
}

func TestFindSingleItemMobileHeadband(t *testing.T) {
	taskItem := taskItem{
		Keywords: []string{"headband"},
		Category: "accessories",
		Size:     "",
		Color:    "red",
	}
	targetItem := SupremeItemMobile{
		"New Era® Big Logo Headband",
		303674,
	}
	supremeItems := []SupremeItemMobile{targetItem}

	foundItem, err := findItemMobile(taskItem, &supremeItems)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, targetItem, foundItem)
}

func TestPickSizeMobile(t *testing.T) {
	item := taskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "",
		Color:    "blue",
	}

	sizes := []SizeMobile{
		SizeMobile{Name: "Small", ID: 59764, StockLevel: 1},
		SizeMobile{Name: "Medium", ID: 59765, StockLevel: 1},
		SizeMobile{Name: "Large", ID: 59766, StockLevel: 1},
		SizeMobile{Name: "XLarge", ID: 59767, StockLevel: 0},
	}

	style := Style{
		ID:    21347,
		Name:  "White Sizes",
		Sizes: sizes,
	}

	itemID, isInStock, err := PickSizeMobile(&item, &style)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 59765, itemID)
	assert.True(t, isInStock)
}

func TestPickSizeMobileNoSize(t *testing.T) {
	item := taskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "",
		Color:    "blue",
	}

	style := Style{ID: 21945, Name: "Navy", Sizes: []SizeMobile{SizeMobile{Name: "N/A", ID: 61615, StockLevel: 0}}}

	itemID, isInStock, err := PickSizeMobile(&item, &style)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 61615, itemID)
	assert.False(t, isInStock)
}

func TestPickSizeMobileNoSizeHeadband(t *testing.T) {
	item := taskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "",
		Color:    "blue",
	}

	style := Style{ID: 23439, Name: "Dark Green", Sizes: []SizeMobile{SizeMobile{Name: "N/A", ID: 50557, StockLevel: 0}}}

	itemID, isInStock, err := PickSizeMobile(&item, &style)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 50557, itemID)
	assert.False(t, isInStock)
}

func TestPickSizeMobileNoTaskSizeIntoSizes(t *testing.T) {
	item := taskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "",
		Color:    "blue",
	}

	sizes := []SizeMobile{
		SizeMobile{Name: "Small", ID: 59764, StockLevel: 1},
		SizeMobile{Name: "Medium", ID: 59765, StockLevel: 1},
		SizeMobile{Name: "Large", ID: 59766, StockLevel: 1},
		SizeMobile{Name: "XLarge", ID: 59767, StockLevel: 0},
	}

	style := Style{
		ID:    21347,
		Name:  "White Sizes",
		Sizes: sizes,
	}

	itemID, isInstock, err := PickSizeMobile(&item, &style)
	assert.Equal(t, errors.New("Unable to pick size, no task size specificed and style not N/A"), err)
	assert.Equal(t, itemID, 0)
	assert.False(t, isInstock)
}

func TestPickSizeMobileTaskSizeIntoNoSize(t *testing.T) {
	item := taskItem{
		Keywords: []string{"temp"},
		Category: "accessories",
		Size:     "",
		Color:    "blue",
	}

	style := Style{ID: 21945, Name: "Navy", Sizes: []SizeMobile{SizeMobile{Name: "N/A", ID: 61615, StockLevel: 0}}}

	itemID, isInStock, err := PickSizeMobile(&item, &style)
	assert.Equal(t, errors.New("Unable to pick size, unable to find size in multiple styles"), err)
	assert.Equal(t, itemID, 0)
	assert.False(t, isInStock)
}

func TestMultipleSizesPickSizesMobile(t *testing.T) {
	sizes := []SizeMobile{
		SizeMobile{Name: "Small", ID: 59764, StockLevel: 1},
		SizeMobile{Name: "Medium", ID: 59765, StockLevel: 0},
		SizeMobile{Name: "Large", ID: 59766, StockLevel: 1},
		SizeMobile{Name: "XLarge", ID: 59767, StockLevel: 0},
	}

	style := Style{
		ID:    21347,
		Name:  "White Sizes",
		Sizes: sizes,
	}

	pickTests := []struct {
		name      string
		in        taskItem
		styleID   int
		isInStock bool
	}{
		{"small", taskItem{Size: "small"}, 59764, true},
		{"medium", taskItem{Size: "Medium"}, 59765, false},
		{"large", taskItem{Size: "LaRGe"}, 59766, true},
		{"xlarge", taskItem{Size: "xLarge"}, 59767, false},
	}

	for _, tt := range pickTests {
		t.Run(tt.name, func(t *testing.T) {
			size, isInStock, err := PickSizeMobile(&tt.in, &style)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.styleID, size)
			assert.Equal(t, tt.isInStock, isInStock)
		})
	}
}
