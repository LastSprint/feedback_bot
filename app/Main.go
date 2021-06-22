package main

import (
	crepo "github.com/LastSprint/feedback_bot/Common/Repo"
	cservices "github.com/LastSprint/feedback_bot/Common/Services"
	"github.com/LastSprint/feedback_bot/Rest"
	"github.com/LastSprint/feedback_bot/Steve/Controllers"
	"github.com/LastSprint/feedback_bot/Steve/Repo"
	"github.com/LastSprint/feedback_bot/Steve/Services"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"
)

type config struct {
	BotSlackId     string `env:"STEVE_SLACK_BOT_ID" envDefault:"U025GCFUK6F"`
	MessageToReply string `env:"STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY"`
	AuthToken      string `env:"STEVE_SLACK_BOT_AUTH_TOKEN,unset"`

	RestrictedAuthorsIds []string `env:"OPS_WTF_RESTRICTED_AUTHORS_IDS" envDefault:"UFH46AX6W"`
	AllowedReportersIds  []string `env:"ALLOWED_REPORTERS_IDS" envDefault:"UFH46AX6W"`
	AllowedChannelsIds   []string `env:"ALLOWED_CHANNELS_IDS" envDefault:"C0251ECG4QP"`

	FeedbackDbFilePath string `env:"FEEDBACK_BOT_DB_FILE_PATH"`

	MongoDBConnectionString string `env:"MONGODB_CONNECTION_STRING,unset" envDefault:"mongodb://root:root@127.0.0.1:6355"`

	SlackChannelIdForNotifications string `env:"SLACK_CHANNEL_ID_FOR_NOTIFICATIONS" envDefault:"UFH46AX6W"`
}

func main() {

	config := config{}

	log.Println("Config")
	log.Println("STEVE_SLACK_BOT_ID: ", config.BotSlackId)
	log.Println("STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY: ", config.MessageToReply)
	log.Println("OPS_WTF_RESTRICTED_AUTHORS_IDS: ", config.RestrictedAuthorsIds)
	log.Println("ALLOWED_REPORTERS_IDS: ", config.AllowedReportersIds)
	log.Println("ALLOWED_CHANNELS_IDS: ", config.AllowedChannelsIds)
	log.Println("SLACK_CHANNEL_ID_FOR_NOTIFICATIONS: ", config.SlackChannelIdForNotifications)

	if err := env.Parse(&config); err != nil {
		log.Fatal("[ERR] Couldn't parse config from env with error", err.Error())
		return
	}

	configureCTOFeedback(config)
	configureSteve(config)

	log.Println("[INFO] Started on :6654")

	log.Fatal(http.ListenAndServe(":6654", nil))
}

func configureSteve(c config) {

	slackRepo := &crepo.SlackRepo{
		AuthToken: c.AuthToken,
	}

	steve := Controllers.EventHandlerController{
		ReplyOnMessageService: &Services.ReplyOnMessageInThreadService{
			BotSlackId:     c.BotSlackId,
			MessageToReply: c.MessageToReply,
			SlackRepo:      slackRepo,
		},
		ConfusingShortcutService: &Services.ConfusingMessageService{
			ConfusingMessagesRepo: &Repo.ConfusingMessagesMongoDBRepo{
				ConnectionString: c.MongoDBConnectionString,
			},
			RestrictedAuthorsIds: c.RestrictedAuthorsIds,
			AllowedReportersIds:  c.AllowedReportersIds,
			AllowedChannels:      c.AllowedChannelsIds,
			NotificationService: &cservices.SlackNotificationService{
				SlackRepo:      slackRepo,
				SlackChannelId: c.SlackChannelIdForNotifications,
			},
		},
	}

	steve.Init()
}

func configureCTOFeedback(c config) {

	slackRepo := &crepo.SlackRepo{
		AuthToken: c.AuthToken,
	}

	controller := Rest.SlackController{
		NotificationService: &cservices.SlackNotificationService{
			SlackRepo:      slackRepo,
			SlackChannelId: c.SlackChannelIdForNotifications,
		},
	}

	controller.Init()
}
