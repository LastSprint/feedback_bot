package Models

// SlackEvent generic slack event entity
// contains meta and payload `Event`
type SlackEvent struct {
	// Token is a sign for validation that event was sent by slack and not by offender
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	ApiAppId string `json:"api_app_id"`
	// EventValue is payload
	EventValue Event `json:"event"`
	// Type can be `event_callback`
	Type      string `json:"type"`
	EventId   string `json:"event_id"`
	EventTime int    `json:"event_time"`
}
