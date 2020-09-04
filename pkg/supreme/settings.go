package supreme

// ApplicationSettings The default application settings
type ApplicationSettings struct {
	DefaultStartTime     string        `json:"startTime"`
	DefaultDelaySettings DelaySettings `json:"defaultDelaySettings"`
}

// DefaultApplicationSettings Default settings used mostly in the task section.
var DefaultApplicationSettings = ApplicationSettings{
	"",
	DelaySettings{
		150,
		500,
		3500,
	},
}
