package Controllers

import (
	"encoding/json"
	models "github.com/LastSprint/feedback_bot/Steve/Models"
	"io/ioutil"
	"log"
	"net/http"
)

type ReplyOnMessageService interface {
	Reply(event models.SlackEvent)
}

type EventHandlerController struct {
	ReplyOnMessageService
}

func (cnt *EventHandlerController) Init() {
	http.HandleFunc("/slack_event", cnt.handleChannelPush)
}

func (cnt *EventHandlerController) handleChannelPush(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	if cnt.handleVerification(w, body) {
		// it's verification request. skip it
		return
	}

	if err != nil {
		log.Println("[ERR] Couldn't read request body " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)

	go func() {

		var event models.SlackEvent

		if err := json.Unmarshal(body, &event); err != nil {
			log.Println("[ERR] Couldn't parse body " + err.Error())
			return
		}
		cnt.ReplyOnMessageService.Reply(event)
	}()
}

func (cnt *EventHandlerController) handleVerification(w http.ResponseWriter, body []byte) bool {
	type verification_token struct {
		Token     string
		Challenge string
		Type      string
	}

	var val verification_token

	json.Unmarshal(body, &val)

	if len(val.Challenge) == 0 {
		return false
	}

	w.Write([]byte(val.Challenge))

	return true
}