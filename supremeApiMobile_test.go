package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindSingleItemMobile(t *testing.T) {
	taskItem := taskItem{
		[]string{"hanes", "boxer"},
		"accessories",
		"Medium",
		"white",
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

func TestPickSizeMobile(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"medium",
		"blue",
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

	itemID, err := PickSizeMobile(&item, &style)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 59765, itemID)
}

func TestPickSizeMobileNoSize(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"",
		"blue",
	}

	style := Style{ID: 21945, Name: "Navy", Sizes: []SizeMobile{SizeMobile{Name: "N/A", ID: 61615, StockLevel: 0}}}

	itemID, err := PickSizeMobile(&item, &style)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 61615, itemID)
}

func TestPickSizeMobileNoTaskSizeIntoSizes(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"",
		"blue",
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

	itemID, err := PickSizeMobile(&item, &style)
	assert.Equal(t, errors.New("Unable to pick size, no task size specificed and style not N/A"), err)
	assert.Equal(t, itemID, 0)
}

func TestPickSizeMobileTaskSizeIntoNoSize(t *testing.T) {
	item := taskItem{
		[]string{"temp"},
		"accessories",
		"medium",
		"blue",
	}

	style := Style{ID: 21945, Name: "Navy", Sizes: []SizeMobile{SizeMobile{Name: "N/A", ID: 61615, StockLevel: 0}}}

	itemID, err := PickSizeMobile(&item, &style)
	assert.Equal(t, errors.New("Unable to pick size, unable to find size in multiple styles"), err)
	assert.Equal(t, itemID, 0)
}

func TestMultipleSizesPickSizesMobile(t *testing.T) {
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

	pickTests := []struct {
		name string
		in   taskItem
		out  int
	}{
		{"small", taskItem{Size: "small"}, 59764},
		{"medium", taskItem{Size: "Medium"}, 59765},
		{"large", taskItem{Size: "LaRGe"}, 59766},
		{"xlarge", taskItem{Size: "xLarge"}, 59767},
	}

	for _, tt := range pickTests {
		t.Run(tt.name, func(t *testing.T) {
			size, err := PickSizeMobile(&tt.in, &style)
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, tt.out, size)
		})
	}
}
