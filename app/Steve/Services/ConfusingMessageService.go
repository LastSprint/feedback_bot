package Services

import (
	"github.com/LastSprint/feedback_bot/Steve/Models"
	"github.com/LastSprint/feedback_bot/Steve/Models/Entry"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type ConfusingMessagesRepo interface {
	Save(message Entry.ConfusingMessage) error
}

type ConfusingMessageService struct {
	ConfusingMessagesRepo

	// AllowedAuthorsIds is array of userIds of people whose messages can be reported
	RestrictedAuthorsIds []string
	// AllowedReportersIds is array of userIds of people why can report
	AllowedReportersIds []string
	// AllowedChannels is array of ids if channels from which reported can be done
	AllowedChannels []string
}

func (srv *ConfusingMessageService) Save(message Models.MessageShortcutCallBackModel) error {

	// validation

	if len(message.Channel.Id) == 0 {
		return errors.Errorf("Channel is empty in %v", message)
	}

	if !contains(srv.AllowedChannels, message.Channel.Id) {
		return errors.Errorf("Report from restricted channel %s", message.Channel.Id)
	}

	if len(message.User.Id) == 0 {
		return errors.Errorf("User is empty in %v", message)
	}

	if !contains(srv.AllowedReportersIds, message.User.Id) {
		return errors.Errorf("Report from restricted user %s", message.User.Id)
	}

	if len(message.Message.User) == 0 {
		return errors.Errorf("User who post is empty in %v", message)
	}

	if contains(srv.RestrictedAuthorsIds, message.Message.User) {
		return errors.Errorf("Reporting of this author is forbidden %s", message.Message.User)
	}

	// saving

	// we create uniq id by concatenate `AuthorId_MessageTS_ChannelId`
	// because MessageTS isn't uniq for all channels (two channels can have 2 messages with same TS)
	// and if something will be broken and 1 channel will have two messages with same TS we add AuthorID
	// (because 1 person can't post 2 messages exactly at the same time)
	messageID := strings.Join([]string{message.Message.User, message.Message.Ts, message.Channel.Id}, "_")

	entry := Entry.ConfusingMessage{
		AuthorId:       message.Message.User,
		ReporterUserId: message.User.Id,
		Text:           message.Message.Text,
		MessageId:      messageID,
		ReportDate:     time.Now(),
	}

	return srv.ConfusingMessagesRepo.Save(entry)
}

func contains(source []string, value string) bool {
	for _, it := range source {
		if it == value {
			return true
		}
	}

	return false
}
