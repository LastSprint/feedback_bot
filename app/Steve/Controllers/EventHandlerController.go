package Controllers

import (
	"encoding/json"
	"fmt"
	models "github.com/LastSprint/feedback_bot/Steve/Models"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

type Dispatcher interface {
	Dispatch(event models.SlackEvent) error
}

type ConfusingShortcutService interface {
	Save(message models.MessageShortcutCallBackModel) error
}

// EventHandlerController handles slack events
type EventHandlerController struct {
	Dispatcher
	ConfusingShortcutService
}

// Init add http handlers for:
//	- `POST /slack_events/steve`
func (cnt *EventHandlerController) Init() {
	http.HandleFunc("/slack_events/steve", cnt.handleChannelPush)
	http.HandleFunc("/commands/ops/wtf", cnt.handleWtfCommand)
}

// handleChannelPush handles situation when somebody write something to channel (doesn't matter in thread or not)
func (cnt *EventHandlerController) handleChannelPush(w http.ResponseWriter, r *http.Request) {

	log.Println("[INFO] Start push handler")

	body, err := ioutil.ReadAll(r.Body)

	if cnt.verifyIfNeeded(w, body) {
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

		if err = cnt.Dispatcher.Dispatch(event); err != nil {
			fmt.Printf("[ERR] An error occurred for event %v\nError: %s", event, err.Error())
		}
	}()
}

// handleWtfCommand handles `/ops_wtf` shortcut
func (cnt *EventHandlerController) handleWtfCommand(_ http.ResponseWriter, r *http.Request) {

	log.Println("[INFO] Get request", r.URL)

	if err := r.ParseForm(); err != nil {
		log.Printf("[ERR] WTF command with request %s. Got error while parsing url form %s", r.URL, err.Error())
		return
	}

	var payload models.MessageShortcutCallBackModel

	if err := json.Unmarshal([]byte(r.Form.Get("payload")), &payload); err != nil {
		log.Printf("[ERR] WTF command with request %s. Got error while parsing payload to json %s", r.URL, err.Error())
		return
	}

	if err := cnt.ConfusingShortcutService.Save(payload); err != nil {
		log.Printf("[ERR] WTF command with request %s. Got error from service %s", r.URL, err.Error())
	}

	log.Printf("[INFO] Report from %s(%s) was written successfully", payload.User.Id, payload.User.Name)
}

// verifyIfNeeded check if request is slack-verification request
// if it's true then return challenge back to slack (write to ResponseWriter) and returns `true`
// if it isn't true then just returns `false`
func (cnt *EventHandlerController) verifyIfNeeded(w http.ResponseWriter, body []byte) bool {
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

func printErrorIfNeeded(err error) {
	if err != nil {
		fmt.Println("[FATAL]", err, string(runtime.ReadTrace()))
	}
}
