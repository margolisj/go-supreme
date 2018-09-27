package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/levigross/grequests"
	log "github.com/sirupsen/logrus"
)

const sharedUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"

var defaultRo = &grequests.RequestOptions{
	UserAgent: sharedUserAgent,
	Headers: map[string]string{
		"accept-language": "en-US,en;q=0.9",
	},
}

//SupremeItem an item found on the supreme webpage
type SupremeItem struct {
	name  string
	color string
	url   string
}

type checkoutJSON struct {
	Status string `json:"status"`
	Slug   string `json:"slug"`
	Errors string `json:"errors"`
}

// SupremeItems a slice of supreme items
type SupremeItems []SupremeItem

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

// GetCollectionItems Gets the collection items from a specific category. If inStockOnly is true then
// the function will only return instock items
func GetCollectionItems(taskItem taskItem, inStockOnly bool) (*SupremeItems, error) {
	collectionURL := "https://www.supremenewyork.com/shop/all/" + taskItem.Category
	resp, err := grequests.Get(collectionURL, defaultRo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp.Ok != true {
		return nil, errors.New("GetCollectionItems request did not return OK")
	}

	// Build goquery doc and find each article
	doc, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	items := parseCategoryPage(doc, inStockOnly)

	return items, nil
}

func parseCategoryPage(doc *goquery.Document, inStockOnly bool) *SupremeItems {
	var items SupremeItems
	doc.Find(".inner-article").Each(func(i int, s *goquery.Selection) {
		// First check sold out status
		soldOut := s.Find("a .sold_out_tag").Size() != 0
		if inStockOnly && soldOut { // Ignore soldout items
			return
		}
		nameSelector := s.Find("h1 .name-link")
		name := nameSelector.Text()
		url, _ := nameSelector.Attr("href")
		color := s.Find("p .name-link").Text()
		items = append(items, SupremeItem{name, color, url})
	})

	return &items
}

// GetSizeInfo Gets st and size options for an item
func GetSizeInfo(session *grequests.Session, itemURLSuffix string) (string, SizeResponse, string, string, error) {
	itemURL := "https://www.supremenewyork.com" + itemURLSuffix
	resp, err := grequests.Get(itemURL, defaultRo)
	if err != nil {
		log.Error(err)
		return "", SizeResponse{}, "", "", err
	}
	if resp.Ok != true {
		return "", SizeResponse{}, "", "", errors.New("GetSizeInfo request did not return OK")
	}

	// Build goquery doc and find each size and style codes
	doc, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	if err != nil {
		log.Error(err)
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
func PickSize(taskItem taskItem, sizes SizeResponse) (string, error) {
	if taskItem.Size == "" {
		return sizes.singleSizeID, nil
	}
	for size, sizeID := range *sizes.multipleSizes {
		if strings.ToLower(taskItem.Size) == strings.ToLower(size) {
			return sizeID, nil
		}
	}

	return "", errors.New("Unable to find size")
}

// AddToCart adds an item to the cart
func AddToCart(session *grequests.Session, addURL string, xcsrf string, st string, s string) (bool, error) {
	localRo := &grequests.RequestOptions{
		UserAgent: sharedUserAgent,
		Headers: map[string]string{
			"accept":           "*/*;q=0.5, text/javascript, application/javascript, application/ecmascript, application/x-ecmascript",
			"accept-encoding":  "gzip, deflate, br",
			"accept-language":  "en-US,en;q=0.9",
			"referer":          addURL,
			"x-csrf-token":     xcsrf,
			"x-requested-with": "XMLHttpRequest",
			"dnt":              "1",
			"origin":           "https://www.supremenewyork.com",
		},
		Data: map[string]string{
			"utf8":   "✓",
			"st":     st,
			"s":      s, // Size
			"commit": "add to cart",
		},
	}

	resp, err := session.Post(
		"https://www.supremenewyork.com"+addURL,
		localRo, // Add ref
	)

	if err != nil {
		log.Error("Error addding to cart", err)
		return false, err
	}

	if resp.Ok != true {
		log.Warn(resp.RawResponse.Request)
		log.Warn(resp.RawResponse)
		return false, errors.New("ATC Req did not return OK")
	}

	return false, nil
}

// FindItem finds a task itime in the slice of supreme items
func findItem(taskItem taskItem, supremeItems SupremeItems) (SupremeItem, error) {
	for _, supItem := range supremeItems {
		if checkKeywords(taskItem.Keywords, supItem.name) && checkColor(taskItem.Color, supItem.color) {
			return supItem, nil
		}
	}

	return SupremeItem{}, errors.New("Unable to match item")
}

func checkKeywords(keywords []string, supremeItemName string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(strings.ToLower(supremeItemName), strings.ToLower(keyword)) {
			return false
		}
	}
	return true
}

// checkColor checks the supreme item color to see if it contains the task color
func checkColor(taskItemColor string, supremeItemColor string) bool {
	if taskItemColor == "" {
		return true
	}
	return strings.Contains(strings.ToLower(strings.TrimSpace(supremeItemColor)), strings.ToLower(taskItemColor))
}

// Checkout Checks out a task. If there is an issue with
func Checkout(i int, session *grequests.Session, xcsrf string, account *Account) (bool, error) {

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
			"accept":           "application/json",
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
		log.Error("Checkout Error: ", err)
		return false, err
	}

	log.Debug("----------------RESPONSE----------------")
	respString := resp.String()
	log.Debugf("%d %v", i, respString)
	log.Debugf("%d %v", i, resp.RawResponse)

	log.Debug("----------------REQUEST----------------")
	log.Debugf("%d %v", i, resp.RawResponse.Request)

	if resp.Ok != true {
		log.Warnf("%d Checkout request did not return OK", i)
		return false, err
	}

	// TODO: Is there a response that doesn't queue? If not we can get rid of redundent
	// return false logic below
	if strings.Contains(respString, "queued") {
		return queue(session, respString)
	} else if strings.Contains(respString, "failed") {
		log.WithFields(log.Fields{
			"reason":   "failed",
			"response": respString,
		}).Error("Checkout failed")
		return false, nil
	} else if strings.Contains(respString, "outOfStock") {
		log.WithFields(log.Fields{
			"reason":   "outOfStock",
			"response": respString,
		}).Error("Checkout failed")
		return false, nil
	}

	return true, nil
}

func queue(session *grequests.Session, respString string) (bool, error) {
	var queueJSON checkoutJSON
	if err := json.Unmarshal([]byte(respString), &queueJSON); err != nil {
		log.WithFields(log.Fields{
			"response": respString,
		}).Error("Unable to marshall json")
		return false, nil
	}

	time.Sleep(10000)

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

	resp, err := session.Get(fmt.Sprintf("https://www.supremenewyork.com/checkout/%s/status.json", queueJSON.Slug), localRo)
	if err != nil {
		log.Error("Queue Error: ", err)
		return false, err
	}

	if resp.Ok != true {
		log.Warn("Queue did not return OK")
		log.Warn(resp.RawResponse.Request)
		log.Warn(resp.RawResponse)
		return false, errors.New("Queue did not return OK")
	}

	if strings.Contains(respString, "queued") {
		return queue(session, resp.String())
	} else if strings.Contains(respString, "failed") {
		log.WithFields(log.Fields{
			"reason":   "failed",
			"response": respString,
		}).Error("Queue failed")
		return false, nil
	} else if strings.Contains(respString, "outOfStock") {
		log.WithFields(log.Fields{
			"reason":   "outOfStock",
			"response": respString,
		}).Error("Queue failed")
		return false, nil
	}

	log.WithFields(log.Fields{
		"respString": respString,
	}).Info("Queue successful")

	return true, nil
}
