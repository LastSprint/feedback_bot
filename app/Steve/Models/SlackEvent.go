package Models

type SlackEvent struct {
	Token      string `json:"token"`
	TeamID     string `json:"team_id"`
	ApiAppId   string `json:"api_app_id"`
	EventValue Event  `json:"event"`
	Type       string `json:"type"`
	EventId    string `json:"event_id"`
	EventTime  int    `json:"event_time"`
}