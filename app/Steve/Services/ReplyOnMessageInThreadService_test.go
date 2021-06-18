package Services

import (
	"github.com/LastSprint/feedback_bot/Steve/Models"
	"testing"
)

// Test cases:
//
// -- Negative
//
//
//
// -- 

func TestReplyOnMessageInThreadService_Reply(t *testing.T) {
	type fields struct {
		BotSlackId     string
		MessageToReply string
		SlackRepo      SlackRepo
	}
	type args struct {
		event Models.SlackEvent
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := &ReplyOnMessageInThreadService{
				BotSlackId:     tt.fields.BotSlackId,
				MessageToReply: tt.fields.MessageToReply,
				SlackRepo:      tt.fields.SlackRepo,
			}
		})
	}
}
