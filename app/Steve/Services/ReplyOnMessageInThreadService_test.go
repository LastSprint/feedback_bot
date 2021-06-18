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
// - Bot replies in 2 specific channels

type SlackServiceStub struct {
	calls int
	replyChannel string
}

func (s *SlackServiceStub) PostMessageToChat(message, channel, threadId string) error {
	s.calls += 1
	s.replyChannel = channel
	return nil
}


// ---- Negative -----

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
		Type: "event_callback",
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

func TestReplyOnMessageInThreadService_Reply_DoesNotReplyInThread(t *testing.T) {
	// Arrange

	slack := SlackServiceStub{}

	botId := "123"

	service := ReplyOnMessageInThreadService{
		BotSlackId:     botId,
		MessageToReply: "123",
		SlackRepo:      &slack,
	}

	event := Models.SlackEvent{
		Type: "event_callback",
		EventValue: Models.Event{
			Type:        "message",
			Subtype:     "",
			User:        botId + "123",
			Ts:          "123",
			Channel:     devOpsAndSaChannelID,
			ThreadTs:    "123",
		},
	}

	// Act

	service.Reply(event)

	// Assert

	assert.Equal(t, 0, slack.calls, "Service must not call SlackRepo")
}

func TestReplyOnMessageInThreadService_Reply_DoesNotReplyOnEditingOrDeletingEvents(t *testing.T) {
	// Arrange

	slack := SlackServiceStub{}

	botId := "123"

	service := ReplyOnMessageInThreadService{
		BotSlackId:     botId,
		MessageToReply: "123",
		SlackRepo:      &slack,
	}

	event := Models.SlackEvent{
		Type: "event_callback",
		EventValue: Models.Event{
			Type:        "message",
			Subtype:     "editing",
			User:        botId + "123",
			Ts:          "123",
			Channel:     devOpsAndSaChannelID,
			ThreadTs:    "123",
		},
	}

	// Act

	service.Reply(event)
	event.EventValue.Subtype = "deleting"
	service.Reply(event)

	// Assert

	assert.Equal(t, 0, slack.calls, "Service must not call SlackRepo")
}

func TestReplyOnMessageInThreadService_Reply_DoesNotReplyInRandomChannel(t *testing.T) {
	// Arrange

	slack := SlackServiceStub{}

	botId := "123"

	service := ReplyOnMessageInThreadService{
		BotSlackId:     botId,
		MessageToReply: "123",
		SlackRepo:      &slack,
	}

	event := Models.SlackEvent{
		Type: "event_callback",
		EventValue: Models.Event{
			Type:        "message",
			User:        botId + "123",
			Ts:          "123",
			Channel:     "123",
		},
	}

	// Act

	service.Reply(event)

	// Assert

	assert.Equal(t, 0, slack.calls, "Service must not call SlackRepo")
}

// ---- Positive ----

func TestReplyOnMessageInThreadService_Reply_ToPerson(t *testing.T) {
	// Arrange

	slack := SlackServiceStub{}

	botId := "123"

	service := ReplyOnMessageInThreadService{
		BotSlackId:     botId,
		MessageToReply: "123",
		SlackRepo:      &slack,
	}

	event := Models.SlackEvent{
		Type: "event_callback",
		EventValue: Models.Event{
			Type:        "message",
			User:        botId + "123",
			Ts:          "123",
			Channel:     devOpsAndSaChannelID,
		},
	}

	// Act

	service.Reply(event)

	// Assert

	assert.Equal(t, 1, slack.calls, "Service must call SlackRepo 1 time")
}

func TestReplyOnMessageInThreadService_Reply_ToRightChannel(t *testing.T) {
	// Arrange

	slack := SlackServiceStub{}

	botId := "123"

	service := ReplyOnMessageInThreadService{
		BotSlackId:     botId,
		MessageToReply: "123",
		SlackRepo:      &slack,
	}

	event := Models.SlackEvent{
		Type: "event_callback",
		EventValue: Models.Event{
			Type:        "message",
			User:        botId + "123",
			Ts:          "123",
			Channel:     devOpsAndSaChannelID,
		},
	}

	// Act

	service.Reply(event)

	// Assert

	assert.Equal(t, 1, slack.calls, "Service must call SlackRepo 1 time")
	assert.Equalf(t, devOpsAndSaChannelID, slack.replyChannel, "Service must call SlackRepo with %s channel id", devOpsAndSaChannelID)
}

func TestReplyOnMessageInThreadService_Reply_ToTestChannel(t *testing.T) {
	// Arrange

	slack := SlackServiceStub{}

	botId := "123"

	service := ReplyOnMessageInThreadService{
		BotSlackId:     botId,
		MessageToReply: "123",
		SlackRepo:      &slack,
	}

	event := Models.SlackEvent{
		Type: "event_callback",
		EventValue: Models.Event{
			Type:        "message",
			User:        botId + "123",
			Ts:          "123",
			Channel:     steveTestChannelId,
		},
	}

	// Act

	service.Reply(event)

	// Assert

	assert.Equal(t, 1, slack.calls, "Service must call SlackRepo 1 time")
}
