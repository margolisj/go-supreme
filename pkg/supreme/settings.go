package supreme

// ApplicationSettings The default application settings
type ApplicationSettings struct {
	StartTime    string `json:"startTime"`
	RefreshWait  int    `json:"refreshWait"`
	AtcWait      int    `json:"atcWait"`
	CheckoutWait int    `json:"checkoutWait"`
}
