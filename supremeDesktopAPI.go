package main

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/levigross/grequests"
)

const sharedUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"

//SupremeItem an item found on the supreme webpage
type SupremeItem struct {
	name  string
	color string
	url   string
}

/* These are the collection "names" -> actual urls
jackets -> https://www.supremenewyork.com/shop/all/jackets
shirts -> https://www.supremenewyork.com/shop/all/shirts
tops/sweaters -> https://www.supremenewyork.com/shop/all/tops_sweaters
sweatshirts -> https://www.supremenewyork.com/shop/all/sweatshirts
pants -> https://www.supremenewyork.com/shop/all/pants
t-shirts -> https://www.supremenewyork.com/shop/all/t-shirts
hats -> https://www.supremenewyork.com/shop/all/hats
bags -> https://www.supremenewyork.com/shop/all/bags
accessories -> https://www.supremenewyork.com/shop/all/accessories
skate -> https://www.supremenewyork.com/shop/all/skate
*/

// GetCollectionItems Gets the items from a specific category. If inStockOnly is true then
// the function will only return in stock items.
func GetCollectionItems(session *grequests.Session, task *Task, inStockOnly bool) (*[]SupremeItem, error) {
	localRo := grequests.RequestOptions{
		UserAgent: sharedUserAgent,
		Headers: map[string]string{
			"accept-language": "en-US,en;q=0.9",
			"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
			// "accept-encoding": "gzip, deflate, br",
			"dnt": "1",
		},
	}
	taskItem := task.Item
	collectionURL := "https://www.supremenewyork.com/shop/all/" + supremeCategoriesDesktop[taskItem.Category]
	resp, err := session.Get(collectionURL, &localRo)
	if err != nil {
		task.Log().Error().Err(err)
		return nil, err
	}
	if resp.Ok != true {
		task.Log().Warn().Msgf("GetCollectionItems Request did not return OK: %d", resp.StatusCode)
		return nil, errors.New("GetCollectionItems request did not return OK")
	}

	// Build goquery doc and find each article
	doc, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	if err != nil {
		task.Log().Error().Err(err)
		return nil, err
	}
	items := parseCategoryPage(doc, inStockOnly)
	task.Log().Debug().Msgf("Items Found: %d in category %s", len(*items), task.Item.Category)
	task.Log().Debug().Msgf("%+v", *items)

	return items, nil
}

func parseCategoryPage(doc *goquery.Document, inStockOnly bool) *[]SupremeItem {
	var items []SupremeItem
	doc.Find(".inner-article").Each(func(i int, s *goquery.Selection) {

		// First check sold out status
		soldOut := s.Find("a .sold_out_tag").Size() != 0
		if inStockOnly && soldOut { // Ignore soldout items
			return
		}
		nameSelector := s.Find(".product-name")
		name := nameSelector.Text()

		variantSelector := s.Find(".product-style .name-link")
		color := variantSelector.Text()
		url, _ := variantSelector.Attr("href")

		items = append(items, SupremeItem{name, color, url})
	})

	return &items
}

// GetSizeInfo gets st and size options for an item by going to the item page
// and retrieving the options from it.
// itemURLStuffix is in the format "/shop/accessories/jdbpyos48/iimyp2ogd"
func GetSizeInfo(session *grequests.Session, task *Task, itemURLSuffix *string) (string, SizeResponse, string, string, error) {
	localRo := grequests.RequestOptions{
		UserAgent: sharedUserAgent,
		Headers: map[string]string{
			"accept-language": "en-US,en;q=0.9",
			"accept":          "accept: text/html, application/xhtml+xml, application/xml",
			// "accept-encoding": "gzip, deflate, br",
			"dnt": "1",
		},
	}
	// Ex. itemURLStuffix = "/shop/accessories/jdbpyos48/iimyp2ogd"
	itemURL := "https://www.supremenewyork.com" + *itemURLSuffix
	resp, err := session.Get(itemURL, &localRo)
	if err != nil {
		return "", SizeResponse{}, "", "", err
	}
	if resp.Ok != true {
		task.Log().Warn().Msgf("GetSizeInfo Request did not return OK: %d", resp.StatusCode)
		return "", SizeResponse{}, "", "", errors.New("GetSizeInfo request did not return OK")
	}

	// Build goquery doc and find each size and style codes
	doc, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	if err != nil {
		return "", SizeResponse{}, "", "", err
	}

	st, _ := doc.Find("#st").Last().Attr("value") //TODO: Fix error / eturn value if it doesnt exist
	sizeResponse := parseSizes(doc)
	addURL, _ := doc.Find("#cart-addf").Attr("action") //TODO: Fix error / eturn value if it doesnt exist

	// Get xcsrf code
	var xcsrf string
	doc.Find("[name=\"csrf-token\"]").Each(func(i int, s *goquery.Selection) {
		xcsrf, _ = s.Attr("content")
	})

	return st, sizeResponse, addURL, xcsrf, nil
}

// SizeResponse holds either a single size or a pointer to a map
// of multiple sizes
type SizeResponse struct {
	singleSizeID  string
	multipleSizes *map[string]string
}

func parseSizes(doc *goquery.Document) SizeResponse {
	// Check for single size and return if found
	singleVal, exists := doc.Find("#s").Attr("value")
	if exists {
		return SizeResponse{
			singleSizeID:  singleVal,
			multipleSizes: nil,
		}
	}

	// Find the multiple sizes
	sizesToID := make(map[string]string)
	doc.Find("#s > option").Each(func(i int, s *goquery.Selection) {
		size := s.Text()
		value, _ := s.Attr("value")
		sizesToID[size] = value
	})
	return SizeResponse{
		singleSizeID:  "",
		multipleSizes: &sizesToID,
	}
}

// PickSize picks a size out of the size map
func PickSize(taskItem *taskItem, sizes *SizeResponse) (string, error) {
	// If the task item is an empty string, task was set up to target no-size item
	if taskItem.Size == "" {
		if sizes.singleSizeID == "" {
			return "", errors.New("Unable to pick size, task size and singleSizeId both empty")
		}
		return sizes.singleSizeID, nil
	}

	// Make sure we found sizes on the page before we check them
	if sizes.multipleSizes != nil {
		for size, sizeID := range *sizes.multipleSizes {
			if strings.EqualFold(taskItem.Size, size) {
				return sizeID, nil
			}
		}
	}

	return "", errors.New("Unable to pick size, unable to find size in multipleSizes")
}

// AddToCart adds an item to the cart
func AddToCart(session *grequests.Session, task *Task, addURL string, xcsrf string, st string, s string) (bool, error) {
	localRo := &grequests.RequestOptions{
		UserAgent: sharedUserAgent,
		Headers: map[string]string{
			"accept":           "*/*;q=0.5, text/javascript, application/javascript, application/ecmascript, application/x-ecmascript",
			"accept-encoding":  "gzip, deflate, br", //TODO: Test if this works without it
			"accept-language":  "en-US,en;q=0.9",
			"referer":          addURL,
			"x-csrf-token":     xcsrf,
			"x-requested-with": "XMLHttpRequest",
			"dnt":              "1",
			"origin":           "https://www.supremenewyork.com",
		},
		Data: map[string]string{
			"utf8":   "✓",
			"st":     st, // Style
			"s":      s,  // Size
			"commit": "add to cart",
		},
	}

	resp, err := session.Post(
		"https://www.supremenewyork.com"+addURL,
		localRo,
	)

	if err != nil {
		task.Log().Error().Err(err).Msgf("Error addding to cart")
		return false, err
	}

	if resp.Ok != true {
		task.Log().Warn().Msgf("%v", resp.RawResponse.Request)
		task.Log().Warn().Msgf("%v", resp.RawResponse)
		task.Log().Warn().Msgf("AddToCart Request did not return OK: %d", resp.StatusCode)
		return false, errors.New("ATC Req did not return OK")
	}
	respString := resp.String()
	task.Log().Info().Msg(resp.String())

	if strings.Contains(respString, "cart-controls-sold-out") {
		log.Info().Msg("Item sold out on ATc attempt")
		return false, nil
	}

	return true, nil
}

// FindItem finds a task item in the slice of supreme items
func findItem(taskItem taskItem, supremeItems []SupremeItem) (SupremeItem, error) {
	for _, supItem := range supremeItems {
		if checkKeywords(taskItem.Keywords, supItem.name) && checkColor(taskItem.Color, supItem.color) {
			return supItem, nil
		}
	}

	return SupremeItem{}, errors.New("Unable to match item")
}

// Checkout Checks out a task. If there is an issue with
func Checkout(session *grequests.Session, task *Task, xcsrf string) (bool, string, error) {
	account := task.Account
	postData := map[string]string{
		"utf8":                     "✓",
		"authenticity_token":       xcsrf,
		"order[billing_name]":      account.Person.Firstname + " " + account.Person.Lastname,
		"order[email]":             account.Person.Email,
		"order[tel]":               account.Person.PhoneNumber,
		"order[billing_address]":   account.Address.Address1,
		"order[billing_address_2]": account.Address.Address2,
		"order[billing_zip]":       account.Address.Zipcode,
		"order[billing_city]":      account.Address.City,
		"order[billing_state]":     account.Address.State,
		"order[billing_country]":   account.Address.Country,
		"asec":                     "Rmasn",
		"same_as_billing_address":  "1",
		"store_credit_id":          "",
		"store_address":            "1",
		"credit_card[nlb]":         account.Card.Number,
		"credit_card[month]":       account.Card.Month,
		"credit_card[year]":        account.Card.Year,
		"credit_card[rvv]":         account.Card.Cvv,
		// "order[terms]":" 0", // Don't think we actually need this other one
		"order[terms]": "1",
		// "credit_card[vval]": "", // No idea what this is still doing here, old credit card info still currently alive on thier website
	}

	localRo := &grequests.RequestOptions{
		UserAgent: sharedUserAgent,
		Headers: map[string]string{
			"accept":           "*/*",
			"accept-encoding":  "gzip, deflate, br",
			"accept-language":  "en-US,en;q=0.9",
			"origin":           "https://www.supremenewyork.com",
			"referer":          "https://www.supremenewyork.com/checkout",
			"x-csrf-token":     xcsrf,
			"x-requested-with": "XMLHttpRequest",
			"content-type":     "application/x-www-form-urlencoded; charset=UTF-8",
			"dnt":              "1",
		},
		Data: postData,
	}

	resp, err := session.Post("https://www.supremenewyork.com/checkout.json", localRo)
	if err != nil {
		task.Log().Error().Err(err).Msg("Checkout Error")
		return false, "", err
	}
	respString := resp.String()

	task.Log().Debug().Msg("----------------RESPONSE----------------")
	task.Log().Debug().Msgf("%s", respString)
	task.Log().Debug().Msgf("%v", resp.RawResponse)

	task.Log().Debug().Msgf("----------------REQUEST----------------")
	task.Log().Debug().Msgf("%v", resp.RawResponse.Request)

	if resp.Ok != true {
		task.Log().Warn().Msgf("Checkout Request did not return OK: %d", resp.StatusCode)
		return false, "", err
	}

	checkoutResponse := handleCheckoutResponse(task, &respString)
	if checkoutResponse {
		return true, respString, nil
	}

	return false, "", nil
}
