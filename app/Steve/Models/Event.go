package Models

// Event contains payload information about event-entity
type Event struct {
	// Type can me `message`
	Type string `json:"type"`
	// Subtype contains information about kind of activity. For example editing or deleting
	Subtype string `json:"subtype"`
	// User is slack ID of user who create this event
	User string `json:"user"`
	// Ts It's timestamp and ID (at the same time) of the message
	// You can use it as ID for creating threads
	Ts   string `json:"ts"`
	Team string `json:"team"`
	// Channel is slack channel id of chat which contains message (or where event occured)
	Channel string `json:"channel"`
	// Event timestamp (or ID)
	EventTs string `json:"event_ts"`
	// ChannelType can me `channel`
	ChannelType string `json:"channel_type"`
	// ThreadTs exists only if this event occured in thread
	ThreadTs string `json:"thread_ts"`
}