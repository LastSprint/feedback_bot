package Services

import (
	"github.com/LastSprint/feedback_bot/Steve/Models"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test cases:
//
// -- Negative
//
// - Bot doesn't reply to itself
// - Bot doesn't reply in thread
// - Bot doesn't reply on message editing or deleting
// - Bot doesn't reply in any random channel
//
// -- Positive
//
// - Bot replies to person
// - Bot replies in thread linked to event message
// - Bot replies on message
// - Bot replies in 2 specific channels

type SlackServiceStub struct {
	calls int
}

func (s *SlackServiceStub) PostMessageToChat(message, channel, threadId string) error {
	s.calls += 1
	return nil
}

func TestReplyOnMessageInThreadService_Reply_DoesNotReplyToItself(t *testing.T) {
	// Arrange

	slack := SlackServiceStub{}

	botId := "123"

	service := ReplyOnMessageInThreadService{
		BotSlackId:     botId,
		MessageToReply: "123",
		SlackRepo:      &slack,
	}

	event := Models.SlackEvent{
		EventValue: Models.Event{
			Type:        "message",
			Subtype:     "",
			User:        botId,
			Ts:          "123",
			Channel:     devOpsAndSaChannelID,
			ThreadTs:    "",
		},
	}

	// Act

	service.Reply(event)

	// Assert

	assert.Equal(t, 0, slack.calls, "Service must not call SlackRepo")
}