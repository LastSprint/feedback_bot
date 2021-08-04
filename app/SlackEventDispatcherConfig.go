package main

import (
	"github.com/LastSprint/feedback_bot/Steve/Repo"
	"github.com/LastSprint/feedback_bot/Steve/Services"
	"github.com/LastSprint/feedback_bot/Steve/Services/Dispatcher"
	"github.com/LastSprint/feedback_bot/Steve/Services/Handlers"
)

func configureSlackEventDispatcher(c config) *Dispatcher.SlackEventDispatcher {
	dispatcher := &Dispatcher.SlackEventDispatcher{}

	dispatcher.RegisterHandler(
		Dispatcher.EventTypeKeyMessage,
		&Handlers.MessageEventHandlerReplyInThread{Service: &Services.ReplyOnMessageInThreadService{
			BotSlackId:         c.BotSlackId,
			AllowedChannelsIds: c.SupportAutomationChannelToReply,
			RequestsRepo: &Repo.RequestsMongoDBRepo{
				ConnectionString: c.MongoDBConnectionString,
			},
		}},
	)

	reactionService := &Services.ReactionService{
		MsgReactionRepo: &Repo.ReactionsMongoDBRepo{
			ConnectionString: c.MongoDBConnectionString,
		},
	}
	reactionsHandler := &Handlers.ReactionHandler{
		AddReactionService:    reactionService,
		DeleteReactionService: reactionService,
	}

	dispatcher.RegisterHandler(Dispatcher.EventTypeKeyAddReaction, reactionsHandler)
	dispatcher.RegisterHandler(Dispatcher.EventTypeKeyRemoveReaction, reactionsHandler)

	return dispatcher
}
