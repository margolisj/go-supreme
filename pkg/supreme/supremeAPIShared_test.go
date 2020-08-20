// +build unit

package supreme

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleCheckoutResponse(t *testing.T) {
	task := testTask()
	tests := []struct {
		name             string
		checkoutResponse string
		out              bool
	}{
		{"queuing", "{\"status\":\"queued\",\"slug\":\"plfspd5ebvuam8hxj\"}", true},
		{"failed", "{\"status\":\"failed\",\"id\":36526604,\"page\":\"\\u003cdiv id=\\\"content\\\"\\u003e\\u003cdiv id=\\\"cart-header\\\"\\u003e\\u003cdiv id=\\\"tabs\\\"\\u003e\\u003cdiv class=\\\"tab\\\"\\u003e\\u003cb\\u003eEDIT / VIEW CART\\u003c/b\\u003e\\u003c/div\\u003e\\u003cdiv class=\\\"tab tab-payment \\\"\\u003e\\u003cb\\u003eSHIPPING / PAYMENT\\u003c/b\\u003e\\u003c/div\\u003e\\u003cdiv class=\\\"tab tab-confirmation selected\\\"\\u003e\\u003cb\\u003eCONFIRMATION\\u003c/b\\u003e\\u003c/div\\u003e\\u003c/div\\u003e\\u003c/div\\u003e\\u003cdiv class=\\\"failed\\\" id=\\\"confirmation\\\"\\u003e\\u003cp\\u003eUnfortunately, we cannot process your payment. This could be due to  your payment being declined by your card issuer. \\u003c/p\\u003e\\u003cp\\u003eYou have not been charged by Supreme, but in some cases your bank may still hold funds in your account. This is only a temporary hold that usually will be removed after 48 to 72 hours, depending on your individual credit card issuer.\\u003c/p\\u003e\\u003cp\\u003eIf you'd like to re-enter your card details \\u003ca id=\\\"back\\\" data-no-turbolink=\\\"true\\\" href=\\\"/checkout\\\"\\u003ego back\\u003c/a\\u003e and try again.\\u003c/p\\u003e\\u003c/div\\u003e\\u003cscript type=\\\"text/javascript\\\"\\u003e\\n  $(function() {\\n    ga_track('ecommerce:addTransaction',\\n      '36526604',           // order ID - required\\n      'Supreme',  // affiliation or store name\\n      '118.00',          // total - required\\n      '0.00',           // tax\\n      '10.00',              // shipping\\n      'Bala Cynwyd',       // city\\n      'PA',     // state or province\\n      'USA'             // country\\n    );\\n\\n      ga_track('ecommerce:addItem',\\n        '36526604',           // order ID - required\\n        '65429',           // SKU/code - required\\n        'Supreme\u00ae/The North Face\u00ae Leather Shoulder Bag',        // product name\\n        'Red N/A',   // category or variation\\n        '118.00',          // unit price - required\\n        '1'               // quantity - required\\n      );\\n\\n    ga_track('ecommerce:send'); //submits transaction to the Analytics servers\\n  });\\n\\n\\u003c/script\\u003e\\u003cscript\\u003ega_track('Purchase Attempt', 'mp_only', [{\\\"Release Date\\\":\\\"10/18/2018\\\",\\\"Release Week\\\":\\\"9FW18\\\",\\\"Total Cart Cost\\\":\\\"118.00\\\",\\\"Currency\\\":\\\"USD\\\",\\\"Shipping City\\\":\\\"Bala Cynwyd\\\",\\\"Shipping Country\\\":\\\"USA\\\",\\\"Season\\\":\\\"FW18\\\",\\\"Page Name\\\":\\\"Confirmation\\\",\\\"Category\\\":\\\"Bags\\\",\\\"Sold Out?\\\":false,\\\"Failure Reason\\\":null,\\\"Product Name\\\":\\\"Supreme\u00ae/The North Face\u00ae Leather Shoulder Bag\\\",\\\"Product Number\\\":\\\"FW18B7\\\",\\\"Product Color\\\":\\\"Red\\\",\\\"Product Size\\\":\\\"N/A\\\",\\\"Product Cost\\\":\\\"118.00\\\",\\\"Event Name\\\":\\\"Purchase Attempt\\\",\\\"Success?\\\":false}]);\\u003c/script\\u003e\\u003c/div\\u003e\"}", false},
		{"invalid card", "{\"status\":\"failed\",\"cart\":[{\"size_id\":\"59765\",\"in_stock\":true}],\"errors\":{\"order\":\"\",\"credit_card\":\"number is not a valid credit card number\"}}", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, handleCheckoutResponse(&task, &tt.checkoutResponse))
		})
	}
}

func TestHandleQueueResponse(t *testing.T) {
	task := testTask()
	tests := []struct {
		name          string
		queueresponse string
		isInQueue     bool
		queueSuccess  bool
	}{
		{"in queue", "{\"status\":\"queued\"}", true, false},
		{"paid", "{\"status\":\"paid\",\"id\":36468926,\"info\":{\"id\":36468926,\"billing_name\":\"Jax Blax\",\"email\":\"somethingWebsite+2zyh@gmail.com\",\"purchases\":[{\"image\":\"//assets.supremenewyork.com/156780/ca/UpYUEolh5WY.jpg\",\"product_name\":\"Supreme®/Hanes® Boxer Briefs (4 Pack)\",\"style_name\":\"White\",\"size_name\":\"Medium\",\"price\":3600,\"product_id\":171745,\"style_id\":21347,\"quantity\":1}],\"item_total\":3600,\"shipping_total\":1000,\"tax_total\":0,\"total\":4600,\"currency\":\"$\",\"store_credit\":0,\"discount_total\":0,\"created_at\":\"Oct 17 at 18:48\",\"service\":\"UPS Ground\",\"is_canada\":false,\"manual_review\":false},\"mpa\":[{\"Release Date\":\"10/11/2018\",\"Release Week\":\"8FW18\",\"Total Cart Cost\":\"36.00\",\"Currency\":\"USD\",\"Shipping City\":\"Philadelphia\",\"Shipping Country\":\"USA\",\"Season\":\"FW18\",\"Page Name\":\"Confirmation\",\"Category\":\"Accessories\",\"Sold Out?\":false,\"Failure Reason\":null,\"Product Name\":\"Supreme®/Hanes® Boxer Briefs (4 Pack)\",\"Product Number\":\"FW18A36\",\"Product Color\":\"White\",\"Product Size\":\"Medium\",\"Product Cost\":\"36.00\",\"Event Name\":\"Purchase Attempt\",\"Success?\":true}],\"mps\":{\"Release Date\":\"10/11/2018\",\"Release Week\":\"8FW18\",\"Total Cart Cost\":\"36.00\",\"Currency\":\"USD\",\"Shipping City\":\"Pittsburg\",\"Shipping Country\":\"USA\",\"Season\":\"FW18\",\"Page Name\":\"Confirmation\",\"Event Name\":\"Purchase Success\",\"Product Names\":[\"Supreme®/Hanes® Boxer Briefs (4 Pack)\"],\"Product Numbers\":[\"FW18A36\"],\"Product Colors\":[\"White\"],\"Product Sizes\":[\"Medium\"],\"Products\":[{\"Name\":\"Supreme®/Hanes® Boxer Briefs (4 Pack)\",\"Color\":\"White\",\"Size\":\"Medium\"}],\"# of Items\":1}}", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isInQueue, queueSuccess := handleQueueResponse(&task, &tt.queueresponse)
			assert.Equal(t, tt.isInQueue, isInQueue)
			assert.Equal(t, tt.queueSuccess, queueSuccess)
		})
	}
}
