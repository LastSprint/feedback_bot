package Services

import "fmt"

// SlackNotificationService provides possibility to post messages in predefined chat
type SlackNotificationService struct {
	SlackRepo

	// SlackChannelId is id of chat where you want to post messages
	SlackChannelId string
}

func (s *SlackNotificationService) Notify(event, message string) error {
	slackMsg := fmt.Sprintf("New Event `%s`\n\nMessage: `%s`", event, message)
	return s.SlackRepo.PostMessageToChat(slackMsg, s.SlackChannelId, "")
}
