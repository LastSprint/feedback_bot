package Handlers

import models "github.com/LastSprint/feedback_bot/Steve/Models"

type ReplyOnMsgInThreadService interface {
	Reply(event models.SlackEvent)
}

type MessageEventHandlerReplyInThread struct {
	Service ReplyOnMsgInThreadService
}

func (h *MessageEventHandlerReplyInThread) Handle(event models.SlackEvent) error {
	h.Service.Reply(event)
	return nil
}
