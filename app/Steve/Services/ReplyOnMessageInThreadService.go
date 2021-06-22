package Services

import (
	models "github.com/LastSprint/feedback_bot/Steve/Models"
	"log"
)

const devOpsAndSaChannelID = "CFSF56EHK"
const steveTestChannelId = "C0251ECG4QP"

type ReplyOnMessageInThreadService struct {
	BotSlackId     string
	MessageToReply string
	SlackRepo
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

	if event.EventValue.Channel != steveTestChannelId && event.EventValue.Channel != devOpsAndSaChannelID {
		log.Printf("[WARN] Got event from channel %s\n", event.EventValue.Channel)
		return
	}

	if err := srv.SlackRepo.PostMessageToChat(srv.MessageToReply, event.EventValue.Channel, event.EventValue.Ts); err != nil {
		log.Printf("[ERR] Got error while posting message to chat %s", err.Error())
	}
}
