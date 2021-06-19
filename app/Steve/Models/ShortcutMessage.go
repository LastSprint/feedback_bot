package Models

// ShortcutMessage message on which the shortcut (MessageShortcutCallBackModel) was called
type ShortcutMessage struct {
	// Type can be `message`
	Type string `json:"type"`
	// Text text of the somebody's message
	Text string `json:"text"`
	// User is id of the user who posted the message
	User string `json:"user"`
	// Ts is message id
	Ts   string `json:"ts"`
	Team string `json:"team"`
	// ThreadTs exists if this message is inside thread
	ThreadTs string `json:"thread_ts"`
	// ParentUserId i have no idea what is it :D
	ParentUserId string `json:"parent_user_id"`
}
