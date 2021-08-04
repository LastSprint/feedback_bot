package Handlers

import (
	"fmt"
	"github.com/LastSprint/feedback_bot/Steve/Models"
	"github.com/LastSprint/feedback_bot/Steve/Services"
	"github.com/LastSprint/feedback_bot/Steve/Services/Dispatcher"
	"reflect"
)

type AddReactionService interface {
	Add(model Services.AddReactionServiceModel) error
}

type DeleteReactionService interface {
	Remove(model Services.AddReactionServiceModel) error
}

type ReactionHandler struct {
	AddReactionService
	DeleteReactionService
}

func (h *ReactionHandler) Handle(event Models.SlackEvent) error {

	reactionModel := Services.AddReactionServiceModel{
		Reaction:  event.EventValue.Reaction,
		Channel:   event.EventValue.Item.ChannelId,
		MessageId: event.EventValue.Item.Ts,
	}

	switch Dispatcher.EventTypeKey(event.EventValue.Type) {
	case Dispatcher.EventTypeKeyAddReaction:
		return h.AddReactionService.Add(reactionModel)
	case Dispatcher.EventTypeKeyRemoveReaction:
		return h.DeleteReactionService.Remove(reactionModel)
	default:
		return fmt.Errorf("event couldn't be handled by %s because event type is uknown", reflect.TypeOf(h))
	}
}
