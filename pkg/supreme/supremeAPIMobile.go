package supreme

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/levigross/grequests"
)

const mobileUserAgent = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0_1 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A402 Safari/604.1"

// SupremeItemMobile models the important information needed for the mobile
// supreme API
type SupremeItemMobile struct {
	name string
	id   int
}

type mobileStockResponse struct {
	ProductsAndCategories map[string][]mobileItem `json:"products_and_categories"`
	LastMobileAPIUpdate   string                  `json:"last_mobile_api_update"`
	ReleaseDate           string                  `json:"release_date"`
	ReleaseWeek           string                  `json:"release_week"`
}

type mobileItem struct {
	Name         string `json:"name"`
	ID           int    `json:"id"`
	ImageURL     string `json:"image_url"`
	ImageURLHi   string `json:"image_url_hi"`
	Price        int    `json:"price"`
	SalePrice    int    `json:"sale_price"`
	NewItem      bool   `json:"new_item"`
	Position     int    `json:"position"`
	CategoryName string `json:"category_name"`
}

// GetCollectionItemsMobile Gets the all the items from a specific category
func GetCollectionItemsMobile(session *grequests.Session, task *Task) (*[]SupremeItemMobile, error) {
	localRo := grequests.RequestOptions{
		UserAgent: mobileUserAgent,
		Headers: map[string]string{
			"accept-language": "en-US,en;q=0.9",
			"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		},
	}
	resp, err := session.Get("https://www.supremenewyork.com/mobile_stock.json", &localRo)
	if err != nil {
		task.Log().Error().Err(err)
		return nil, errors.New("Error getting mobile stock")
	}
	if resp.Ok != true {
		task.Log().Warn().Msgf("GetCollectionItemsMobile Request did not return OK: %d", resp.StatusCode)
		return nil, errors.New("GetCollectionItemsMobile request did not return OK")
	}

	var stock mobileStockResponse
	err = resp.JSON(&stock)
	if err != nil {
		return nil, err
	}

	targetCategory, ok := supremeCategoriesMobile[task.Item.Category]
	if !ok {
		return nil, errors.New("Catgeory was not found in supremeCategoriesMobile")
	}
	mobileItems, ok := stock.ProductsAndCategories[targetCategory]
	if !ok {
		return nil, errors.New("Category was incorrect")
	}

	var items []SupremeItemMobile
	for _, item := range mobileItems {
		items = append(items, SupremeItemMobile{item.Name, item.ID})
	}
	task.Log().Debug().Msgf("Items Found: %d in category %s", len(items), task.Item.Category)

	return &items, nil
}

type mobileStylesResponse struct {
	Styles []Style `json:"styles"`
}

// Style is the different colors and sizes of the item
type Style struct {
	ID int `json:"id"`
	// Name of the name of the style aka the color of the item
	Name  string       `json:"name"`
	Sizes []SizeMobile `json:"sizes"`
}

// SizeMobile is the size object in the style object
type SizeMobile struct {
	// Name of the name of the size, aka the size of the item
	Name       string `json:"name"`
	ID         int    `json:"id"`
	StockLevel int    `json:"stock_level"`
}

// GetSizeInfoMobile gets the size information from the item page
func GetSizeInfoMobile(session *grequests.Session, task *Task, item *SupremeItemMobile) (*[]Style, error) {
	localRo := grequests.RequestOptions{
		UserAgent: mobileUserAgent,
		Headers: map[string]string{
			"accept-language": "en-US,en;q=0.9",
			"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		},
	}
	resp, err := session.Get(fmt.Sprintf("https://www.supremenewyork.com/shop/%d.json", item.id), &localRo)
	if err != nil {
		task.Log().Error().Err(err)
		return nil, errors.New("Error during GetSizeInfoMobile")
	}
	if resp.Ok != true {
		task.Log().Warn().Msgf("GetSizeInfoMobile Request did not return OK: %d", resp.StatusCode)
		return nil, errors.New("GetSizeInfoMobile request did not return OK")
	}

	var styleResponse mobileStylesResponse
	err = resp.JSON(&styleResponse)
	if err != nil {
		return nil, errors.New("Error unmarshaling styleResponseMobile")
	}

	return &styleResponse.Styles, nil
}

func findItemMobile(taskItem TaskItem, itemsMobile *[]SupremeItemMobile) (SupremeItemMobile, error) {
	for _, item := range *itemsMobile {
		if checkKeywords(taskItem.Keywords, item.name) {
			return item, nil
		}
	}

	return SupremeItemMobile{}, errors.New("Unable to match Item")
}

// PickSizeMobile picks a size out of the style list
func PickSizeMobile(taskItem *TaskItem, style *Style) (int, bool, error) {
	// If the task item is an empty string, task was set up to target no-size item
	if taskItem.Size == "" {
		if len(style.Sizes) != 1 && style.Sizes[0].Name != "N/A" {
			return 0, false, errors.New("Unable to pick size, no task size specificed and style not N/A")
		}
		return style.Sizes[0].ID, style.Sizes[0].StockLevel == 1, nil
	}

	// Make sure we found sizes on the page before we check them
	for _, size := range style.Sizes {
		if strings.ToLower(taskItem.Size) == strings.ToLower(size.Name) {
			return size.ID, size.StockLevel == 1, nil
		}
	}

	return 0, false, errors.New("Unable to pick size, unable to find size in multiple styles")
}

type atcResponseMobile []struct {
	SizeID  string `json:"size_id"`
	InStock bool   `json:"in_stock"`
}

// AddToCartMobile adds the item to the cart
func AddToCartMobile(session *grequests.Session, task *Task, ID int, st int, s int) (bool, error) {
	localRo := grequests.RequestOptions{
		UserAgent: mobileUserAgent,
		Headers: map[string]string{
			"accept-language":  "en-US,en;q=0.9",
			"accept":           "application/json",
			"reffer":           "http://www.supremenewyork.com/mobile",
			"x-requested-with": "XMLHttpRequest",
			// "accept-encoding":  "gzip, deflate",
		},
		Data: map[string]string{
			"qty": "1",
			"st":  strconv.Itoa(st), // Style
			"s":   strconv.Itoa(s),  // Size ID
		},
	}
	resp, err := session.Post(fmt.Sprintf("https://www.supremenewyork.com/shop/%d/add.json", ID), &localRo)
	if err != nil {
		task.Log().Error().Err(err).Msg("Checkout Error")
		return false, err
	}
	respString := resp.String()

	task.Log().Debug().Msgf("ATC Response: %s", respString)

	if resp.Ok != true {
		task.Log().Warn().Msgf("ATC Request did not return OK: %d", resp.StatusCode)
		return false, errors.New("ATC request did not return OK")
	}

	var atcResponse atcResponseMobile
	task.Log().Debug().Msg(respString)
	if err := json.Unmarshal([]byte(respString), &atcResponse); err != nil {
		task.Log().Error().
			Str("response", respString).
			Msg("Unable to marshall json in atcMobile")
		return false, nil
	}

	if len(atcResponse) == 0 {
		return false, errors.New("ATC Unsuccessfull, empty response")
	}

	return true, nil
}

// CheckoutMobile checks out with the mobile api
func CheckoutMobile(session *grequests.Session, task *Task, cookieSub *string) (bool, string, error) {
	account := task.Account
	// %7B%2259765%22%3A1%7D => {"59765":1}
	postData := map[string]string{
		"store_credit_id":          "",
		"from_mobile":              "1",
		"cookie-sub":               *cookieSub,
		"same_as_billing_address":  "1",
		"order[billing_name]":      account.Person.Firstname + " " + account.Person.Lastname,
		"order[email]":             account.Person.Email,
		"order[tel]":               account.Person.PhoneNumber,
		"order[billing_address]":   account.Address.Address1,
		"order[billing_address_2]": account.Address.Address2,
		"order[billing_zip]":       account.Address.Zipcode,
		"order[billing_city]":      account.Address.City,
		"order[billing_state]":     account.Address.State,
		"order[billing_country]":   account.Address.Country,
		"store_address":            "1",
		"credit_card[cnb]":         account.Card.Number,
		"credit_card[month]":       account.Card.Month,
		"credit_card[year]":        account.Card.Year,
		"credit_card[rsusr]":       account.Card.Cvv,
		"order[terms]":             "1",
		// "g-recaptcha-response": gcap_response,
		"is_from_ios_native": "1",
	}

	localRo := &grequests.RequestOptions{
		UserAgent: mobileUserAgent,
		Headers: map[string]string{
			"accept":          "application/json",
			"accept-encoding": "gzip, deflate, br",
			"accept-language": "en-US,en;q=0.9",
			"origin":          "https://www.supremenewyork.com",
			"referer":         "https://www.supremenewyork.com/checkout",
		},
		Data: postData,
	}

	resp, err := session.Post("https://www.supremenewyork.com/checkout.json", localRo)

	if err != nil {
		task.Log().Error().Err(err).Msg("Checkout Mobile Error")
		return false, "", err
	}

	task.Log().Debug().Msg("----------------RESPONSE----------------")
	respString := resp.String()
	task.Log().Debug().Msgf("%s", respString)
	task.Log().Debug().Msgf("%v", resp.RawResponse)

	task.Log().Debug().Msgf("----------------REQUEST----------------")
	task.Log().Debug().Msgf("%v", resp.RawResponse.Request)

	if resp.Ok != true {
		task.Log().Warn().Msgf("Checkout request did not return OK: %d", resp.StatusCode)
		return false, "", errors.New("Checkout request did not return OK")
	}

	checkoutResponse := handleCheckoutResponse(task, &respString)
	if checkoutResponse {
		return true, respString, nil
	}

	return false, "", nil
}
