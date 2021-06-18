package Models

type Event struct {
	Type        string `json:"type"`
	Subtype     string `json:"subtype"`
	User        string `json:"user"`
	Ts          string `json:"ts"`
	Team        string `json:"team"`
	Channel     string `json:"channel"`
	EventTs     string `json:"event_ts"`
	ChannelType string `json:"channel_type"`
	ThreadTs    string `json:"thread_ts"`
}