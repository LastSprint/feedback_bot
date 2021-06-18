package main

import (
	"github.com/LastSprint/feedback_bot/DB"
	"github.com/LastSprint/feedback_bot/Rest"
	"github.com/LastSprint/feedback_bot/Steve/Controllers"
	"github.com/LastSprint/feedback_bot/Steve/Repo"
	"github.com/LastSprint/feedback_bot/Steve/Services"
	"log"
	"net/http"
	"os"
)

func main() {

	path := os.Getenv("FEEDBACK_BOT_DB_FILE_PATH")

	steve := Controllers.EventHandlerController{
		&Services.ReplyOnMessageInThreadService{
			BotSlackId:     os.Getenv("STEVE_SLACK_BOT_ID"),
			MessageToReply: os.Getenv("STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY"),
			SlackRepo:      &Repo.SlackRepo{
				AuthToken: os.Getenv("STEVE_SLACK_BOT_AUTH_TOKEN"),
			},
		},
	}

	db := DB.FileDB{FilePath: path}
	controller := Rest.SlackController{DB: &db}
	controller.Init()
	steve.Init()
	log.Fatal(http.ListenAndServe(":6654", nil))
}