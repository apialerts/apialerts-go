package model

type APIAlertsEvent struct {
	Channel string   `json:"channel"`
	Message string   `json:"message"`
	Tags    []string `json:"tags"`
	Link    string   `json:"link"`
}
