package barkserver

type Message struct {
	Title             string `json:"title"`
	Body              string `json:"body"`
	DeviceKey         string `json:"device_key"`
	Category          string `json:"category"`
	Badge             int    `json:"badge,omitempty"`
	Sound             string `json:"sound,omitempty"`
	Icon              string `json:"icon,omitempty"`
	Group             string `json:"group,omitempty"`
	URL               string `json:"url,omitempty"`
	Level             string `json:"level,omitempty"`
	AutomaticallyCopy string `json:"automaticallyCopy,omitempty"`
	Copy              string `json:"copy,omitempty"`
	IsArchive         string `json:"isArchive,omitempty"`
}
