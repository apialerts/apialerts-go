package apialerts

type Event struct {
	Message string         `json:"message"`
	Channel string         `json:"channel,omitempty"`
	Event   string         `json:"event,omitempty"`
	Title   string         `json:"title,omitempty"`
	Tags    []string       `json:"tags,omitempty"`
	Link    string         `json:"link,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
}

type Result struct {
	Workspace string
	Channel   string
	Warnings  []string
}
