package main

import (
	crepo "github.com/LastSprint/feedback_bot/Common/Repo"
	cservices "github.com/LastSprint/feedback_bot/Common/Services"
	"github.com/LastSprint/feedback_bot/Rest"
	"github.com/LastSprint/feedback_bot/Steve/Controllers"
	"github.com/LastSprint/feedback_bot/Steve/Repo"
	"github.com/LastSprint/feedback_bot/Steve/Services"
	"github.com/LastSprint/feedback_bot/Steve/Services/TechnologyWidspreadAnalyze"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"

	jiraService "github.com/LastSprint/JiraGoIssues/services"
)

type config struct {
	BotSlackId     string `env:"STEVE_SLACK_BOT_ID" envDefault:"U025GCFUK6F"`
	MessageToReply string `env:"STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY"`
	AuthToken      string `env:"STEVE_SLACK_BOT_AUTH_TOKEN,unset"`

	RestrictedAuthorsIds []string `env:"OPS_WTF_RESTRICTED_AUTHORS_IDS" envDefault:"UFH46AX6W"`
	AllowedReportersIds  []string `env:"ALLOWED_REPORTERS_IDS" envDefault:"UFH46AX6W"`
	AllowedChannelsIds   []string `env:"ALLOWED_CHANNELS_IDS" envDefault:"C0251ECG4QP"`

	DevOpsAndSAChannelId string `env:"STEVE_DEVOPS_AND_SA_CHANNEL_ID" envDefault:"CFSF56EHK"`

	// SupportAutomationChannelToReply is array of channels in which the bot can post message as reply on `post message` event
	// for details look at ReplyOnMessageInThreadService
	SupportAutomationChannelToReply []string `env:"SUPPORT_AUTOMATION_CHANNELS_TO_REPLY" envDefault:"CFSF56EHK,C0251ECG4QP"`

	FeedbackDbFilePath string `env:"FEEDBACK_BOT_DB_FILE_PATH"`

	MongoDBConnectionString string `env:"MONGODB_CONNECTION_STRING,unset" envDefault:"mongodb://root:root@127.0.0.1:6355"`

	SlackChannelIdForNotifications string `env:"SLACK_CHANNEL_ID_FOR_NOTIFICATIONS" envDefault:"UFH46AX6W"`

	JiraAuthApiToken string `env:"JIRA_AUTH_API_TOKEN,unset"`
}

func main() {

	config := config{}

	if err := env.Parse(&config); err != nil {
		log.Fatal("[ERR] Couldn't parse config from env with error", err.Error())
		return
	}

	log.Println("Config")
	log.Println("STEVE_SLACK_BOT_ID: ", config.BotSlackId)
	log.Println("STEVE_SLACK_BOT_DEVOPS_INFO_MESSAGE_TO_REPLY: ", config.MessageToReply)
	log.Println("OPS_WTF_RESTRICTED_AUTHORS_IDS: ", config.RestrictedAuthorsIds)
	log.Println("ALLOWED_REPORTERS_IDS: ", config.AllowedReportersIds)
	log.Println("ALLOWED_CHANNELS_IDS: ", config.AllowedChannelsIds)
	log.Println("SLACK_CHANNEL_ID_FOR_NOTIFICATIONS: ", config.SlackChannelIdForNotifications)

	configureCTOFeedback(config)
	configureSteveEvents(config)
	configureSteveCommands(config)
	configureSteveAnalyzesCommands(config)

	log.Println("[INFO] Started on :6654")

	log.Fatal(http.ListenAndServe(":6654", nil))
}

func configureSteveEvents(c config) {

	slackRepo := &crepo.SlackRepo{
		AuthToken: c.AuthToken,
	}

	slackEventDispatcher := configureSlackEventDispatcher(c)

	steveEvents := Controllers.EventHandlerController{
		Dispatcher: slackEventDispatcher,
		ConfusingShortcutService: &Services.ConfusingMessageService{
			ConfusingMessagesRepo: &Repo.ConfusingMessagesMongoDBRepo{
				ConnectionString: c.MongoDBConnectionString,
			},
			RestrictedAuthorsIds: c.RestrictedAuthorsIds,
			AllowedReportersIds:  c.AllowedReportersIds,
			AllowedChannels:      c.AllowedChannelsIds,
			MessageToReply:       c.MessageToReply,
			SlackRepo:            slackRepo,
			NotificationService: &cservices.SlackNotificationService{
				SlackRepo:      slackRepo,
				SlackChannelId: c.SlackChannelIdForNotifications,
			},
		},
	}

	steveEvents.Init()
}

func configureSteveCommands(c config) {
	steve := Controllers.CommandHandlerController{
		SaWeekStatService: &Services.OpsAndSaStatisticsService{
			SAStatisticRequestRepo: &Repo.RequestsMongoDBRepo{
				ConnectionString: c.MongoDBConnectionString,
			},
			SAStatReportsRepo: &Repo.ConfusingMessagesMongoDBRepo{
				ConnectionString: c.MongoDBConnectionString,
			},
			PublicSaRequestsChannelId: c.DevOpsAndSAChannelId,
		},
		TimeSpentAnalyticsService: &Services.TimeSpendDistributionService{
			WorkLogRepo: &crepo.JiraTempoTimesheetRepo{
				ApiToken: c.JiraAuthApiToken,
			},
		},
	}

	steve.Init()
}

func configureSteveAnalyzesCommands(c config) {

	jiraRepo := jiraService.NewJiraIssueLoader("https://jira.surfstudio.ru/rest/api/2/search", c.JiraAuthApiToken)

	steve := Controllers.CommandAnalyticsController{
		Templates: &TechnologyWidspreadAnalyze.TemplatesService{
			JiraRepo: jiraRepo,
		},
		SurfGen: &TechnologyWidspreadAnalyze.SurfGenService{
			JiraRepo: jiraRepo,
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
