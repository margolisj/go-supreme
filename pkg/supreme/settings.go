package supreme

// Versioning and information for keygen.sh
const (
	keygenAccountID string = "e99bd6f7-900f-4bed-a440-f445fc572fc6"
	keygenProductID string = "a7e001f3-3194-4927-88eb-dd37366ab8ed"
	version         string = "0.0.8"
)

type applicationSettings struct {
	StartTime    string `json:"startTime"`
	RefreshWait  int    `json:"refreshWait"`
	AtcWait      int    `json:"atcWait"`
	CheckoutWait int    `json:"checkoutWait"`
}

// appSettings are the default application settings
var appSettings = applicationSettings{
	"",
	300,
	800,
	4500,
}
