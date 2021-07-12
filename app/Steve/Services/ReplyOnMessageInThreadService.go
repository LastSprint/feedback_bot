package Services

import (
	"github.com/LastSprint/feedback_bot/Common/Utils"
	models "github.com/LastSprint/feedback_bot/Steve/Models"
	"log"
)

type RequestsRepo interface {
	IncrementRequestCount(channelID string) error
}

type ReplyOnMessageInThreadService struct {
	BotSlackId string
	RequestsRepo
	AllowedChannelsIds []string
}

func (srv *ReplyOnMessageInThreadService) Reply(event models.SlackEvent) {

	if event.EventValue.User == srv.BotSlackId {
		log.Println("[INFO] This is message from this bot")
		// it's id of this bot. It shouldn't reply to itself
		return
	}

	if event.EventValue.ThreadTs != "" {
		log.Println("[INFO] This is message inside thread")
		// empty ThreadTs means that message wasn't in thread. We don't want to answer on messages form threads
		return
	}

	if event.EventValue.Subtype != "" {
		// subtype is for editing
		log.Println("[INFO] This is message with subtype " + event.EventValue.Subtype)
		return
	}

	if event.Type != "event_callback" {
		log.Printf("[WARN] Got event with type %s\n", event.Type)
		return
	}

	if event.EventValue.Type != "message" {
		log.Printf("[WARN] Got event with value type %s\n", event.EventValue.Type)
		return
	}

	if !Utils.Contains(srv.AllowedChannelsIds, event.EventValue.Channel) {
		log.Printf("[INFO] Avert replying on message from channel %s", event.EventValue.Channel)
		return
	}

	// we decided to stop posting messages automatically

	//if err := srv.SlackRepo.PostMessageToChat(srv.MessageToReply, event.EventValue.Channel, event.EventValue.Ts); err != nil {
	//	log.Printf("[ERR] Got error while posting message to chat %s", err.Error())
	//}

	log.Printf("[INFO] Success reply on event: %v", event)

	if srv.RequestsRepo == nil {
		log.Println("[ERR] MongoDB repo in ReplyOnMessageInThreadService is nill. Abort writing operation")
		return
	}

	if err := srv.RequestsRepo.IncrementRequestCount(event.EventValue.Channel); err != nil {
		log.Printf("[ERR] while incrementing request count %s", err.Error())
	}
}
