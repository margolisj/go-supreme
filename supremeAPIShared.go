package main

import "strings"

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
	"t-shirts":      "T-shirts", // Currently unknown
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
