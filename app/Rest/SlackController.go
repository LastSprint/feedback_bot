package Rest

import (
	"net/http"
)

type DbService interface {
	WriteFeedbackMessage(string) error
}

type SlackController struct {
	DB DbService
}

func (c *SlackController) Init() {
	http.HandleFunc("/slack/cto/feedback", c.handleFeedbackCommand)
}


func (c* SlackController) handleFeedbackCommand(w http.ResponseWriter, r *http.Request) {
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

	if err := c.DB.WriteFeedbackMessage(messageText); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Thank you for feedback :cat-smiling:"))
}