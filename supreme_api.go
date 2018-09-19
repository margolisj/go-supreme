package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/levigross/grequests"
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

// SupremeItems a slice of supreme items
type SupremeItems []SupremeItem

//GetCollectionItems Gets the collection items
//TODO: make this match?
func GetCollectionItems(taskItem taskItem, inStockOnly bool) *SupremeItems {
	collectionURL := "https://www.supremenewyork.com/shop/all/" + taskItem.category

	resp, err := grequests.Get(collectionURL, defaultRo)

	if err != nil {
		log.Println(err)
	}
	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	// log.Println(resp.String())
	doc, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	var items SupremeItems
	doc.Find(".inner-article").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		soldOut := s.Find("a .sold_out_tag").Size() == 0
		if inStockOnly && !soldOut {
			return
		}
		nameSelector := s.Find("h1 .name-link")
		name := nameSelector.Text()
		url, _ := nameSelector.Attr("href")
		color := s.Find("p .name-link").Text()
		fmt.Printf("%s %s %t %s\n", name, color, soldOut, url)
		items = append(items, SupremeItem{name, color, url})
	})

	log.Println("Found items in collection")
	return &items
}

// GetSizeInfo Gets st and size options for an item
func GetSizeInfo(session *grequests.Session, itemURLSuffix string) (string, *map[string]string, string, string) {
	itemURL := "https://www.supremenewyork.com" + itemURLSuffix

	resp, err := grequests.Get(itemURL, defaultRo)

	if err != nil {
		log.Println(err)
	}
	if resp.Ok != true {
		log.Println("Request did not return OK")
	}
	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	// log.Println(resp.String())
	doc, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	st, _ := doc.Find("#st").Last().Attr("value")

	sizesToID := make(map[string]string)
	doc.Find("#s > option").Each(func(i int, s *goquery.Selection) {
		size := s.Text()
		value, _ := s.Attr("value")
		sizesToID[size] = value
		fmt.Printf("%s %s\n", size, value)
	})

	addURL, _ := doc.Find("#cart-addf").Attr("action")

	var xcsrf string
	doc.Find("[name=\"csrf-token\"]").Each(func(i int, s *goquery.Selection) {
		xcsrf, _ = s.Attr("content")
	})

	return st, &sizesToID, addURL, xcsrf
}

// AddToCart adds an item to the cart
func AddToCart(session *grequests.Session, addURL string, xcsrf string, st string, s string) bool {
	// localRo := copy(&defaultRo)
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
		log.Fatal("Error addding to cart", err)
		return false
	}

	if resp.Ok != true {
		log.Println("ATC Req did not return OK")
		log.Println(resp.RawResponse.Request)
		log.Println(resp.RawResponse)
		return false
	}

	return true
}

func fakeFirstCheckout(session *grequests.Session, xcsrf string, account *account) {
	queryData := map[string]string{
		"utf8":                     "✓",
		"authenticity_token":       xcsrf,
		"order[billing_name]":      account.person.firstname + " " + account.person.lastname,
		"order[email]":             account.person.email,
		"order[tel]":               account.person.phoneNumber,
		"order[billing_address]":   account.address.address1,
		"order[billing_address_2]": account.address.address2,
		"order[billing_zip]":       account.address.zipcode,
		"order[billing_city]":      account.address.city,
		"order[billing_state]":     account.address.state,
		"order[billing_country]":   account.address.country,
		"asec":                     "Rmasn",
		"same_as_billing_address":  "1",
		"store_credit_id":          "",
		"credit_card[nlb]":         "",
		"credit_card[month]":       account.card.month,
		"credit_card[year]":        account.card.year,
		"credit_card[rvv]":         "",
		"order[terms]":             "0",
		"g-recaptcha-response":     "",
		"credit_card[vval]":        "",
		"cnt":                      "1",
	}
	ro := &grequests.RequestOptions{
		Params: queryData,
	}

	resp, err := grequests.Get("https://www.supremenewyork.com/checkout.js", ro)

	if err != nil {
		log.Fatal("Js Checkout Error: ", err)
	}

	if resp.Ok != true {
		log.Println("Js Checkout request did not return OK")
		log.Println(resp.RawResponse.Request)
	}

}

// Checkout Checks out a task
func Checkout(session *grequests.Session, xcsrf string, account *account) bool {
	// fakeFirstCheckout(session, xcsrf, account)

	postData := map[string]string{
		"utf8":                     "✓",
		"authenticity_token":       xcsrf,
		"order[billing_name]":      account.person.firstname + " " + account.person.lastname,
		"order[email]":             account.person.email,
		"order[tel]":               account.person.phoneNumber,
		"order[billing_address]":   account.address.address1,
		"order[billing_address_2]": account.address.address2,
		"order[billing_zip]":       account.address.zipcode,
		"order[billing_city]":      account.address.city,
		"order[billing_state]":     account.address.state,
		"order[billing_country]":   account.address.country,
		"asec":                     "Rmasn",
		"same_as_billing_address":  "1",
		"store_credit_id":          "",
		"store_address":            "1",
		"credit_card[nlb]":         account.card.number,
		"credit_card[month]":       account.card.month,
		"credit_card[year]":        account.card.year,
		"credit_card[rvv]":         account.card.cvv,
		// "order[terms]":" 0", // Don't think we actually need this other one
		"order[terms]": "1",
		// "credit_card[vval]": "", // No idea what this is still doing here
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
		log.Fatal("Checkout Error: ", err)
		return false
	}

	if resp.Ok != true {
		log.Println("Request did not return OK")
		log.Println(resp.RawResponse.Request)
		log.Println(resp.RawResponse)
	}

	// var inter interface{}
	// err = resp.JSON(inter)
	// if err != nil {
	// 	log.Fatal("Error marshalling json", err)
	// } else {
	// 	log.Println(inter)
	// }
	log.Println("----------------RESPONSE----------------")
	log.Println(resp.String())
	log.Println(resp.RawResponse)

	log.Println("----------------REQUEST----------------")
	log.Println(resp.RawResponse.Request)
	log.Println(resp.RawResponse.Request.PostForm)

	return true
}

// func queue(session *grequests.Session, json interface) (bool, error) {

// 	if val, ok := json["status"]; ok {
// 		//do something here
// 	}

// 	return false, error
// }

func checkKeywords(keywords []string, supremeItemName string) bool {
	for _, keyword := range keywords {
		if !strings.Contains(strings.ToLower(supremeItemName), strings.ToLower(keyword)) {
			return false
		}
	}
	return true
}

func checkColor(taskItemColor string, supremeItemColor string) bool {
	if taskItemColor == "" {
		return true
	}
	return strings.ToLower(taskItemColor) == strings.ToLower(supremeItemColor)
}

// FindItem finds a task itime in the slice of supreme items
func FindItem(taskItem taskItem, supremeItems SupremeItems) (SupremeItem, error) {

	for _, supItem := range supremeItems {
		if checkKeywords(taskItem.keywords, supItem.name) && checkColor(taskItem.color, supItem.color) {
			return supItem, nil
		}
	}

	return SupremeItem{}, errors.New("Unable to match item")
}
