package apialerts

type Event struct {
	Channel string   `json:"channel"`
	Event   string   `json:"event"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Tags    []string `json:"tags"`
	Link    string   `json:"link"`
}

type Result struct {
	Workspace string
	Channel   string
	Warnings  []string
}
