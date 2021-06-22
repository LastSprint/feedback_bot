package Services

import "fmt"

// SlackRepo provides possibility to work with Slack
type SlackRepo interface {
	// PostMessageToChat should post `message` to specific `channel`
	// and if it's possible - post message exactly in specific thread (or create a new thread) by `threadId`
	PostMessageToChat(message, channel, threadId string) error
}

// SlackNotificationService provides possibility to post messages in predefined chat
type SlackNotificationService struct {
	SlackRepo

	// SlackChannelId is id of chat where you want to post messages
	SlackChannelId string
}

// Notify will post formatted message to slack
func (s *SlackNotificationService) Notify(event, message string) error {
	slackMsg := fmt.Sprintf("New Event `%s`\n\nMessage: `%s`", event, message)
	return s.SlackRepo.PostMessageToChat(slackMsg, s.SlackChannelId, "")
}
