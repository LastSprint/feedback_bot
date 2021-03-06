package Rest

import (
	"log"
	"net/http"
)

type NotificationService interface {
	Notify(event, message string) error
}

type SlackController struct {
	NotificationService
}

func (c *SlackController) Init() {
	http.HandleFunc("/slack/cto/feedback", c.handleFeedbackCommand)
}

func (c *SlackController) handleFeedbackCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("This http method isn't supported"))
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// https://api.slack.com/interactivity/slash-commands
	messageText := r.FormValue("text")

	if len(messageText) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Seems like that message is empty :("))
		return
	}

	if err := c.NotificationService.Notify("CTO_Feedback", messageText); err != nil {
		log.Printf("[ERR] Couldn't sent feedback to slack with error %s. The message: %s", err.Error(), messageText)
		w.Write([]byte("Couldn't write the message. Something went wrong with feedback bot :("))
		return
	}

	w.Write([]byte("Thank you for feedback :cat-smiling:"))
}
