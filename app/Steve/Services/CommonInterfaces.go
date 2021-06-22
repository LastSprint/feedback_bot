package Services

// SlackRepo provides possibility to work with Slack
type SlackRepo interface {
	// PostMessageToChat should post `message` to specific `channel`
	// and if it's possible - post message exactly in specific thread (or create a new thread) by `threadId`
	PostMessageToChat(message, channel, threadId string) error
}
