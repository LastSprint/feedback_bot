package Models

import "github.com/LastSprint/feedback_bot/Steve/Models/Entry"

// MessageShortcutCallBackModel is entity which is used by slack to notify about shortcuts execution
type MessageShortcutCallBackModel struct {
	// Type can be `message_action`
	Type     string `json:"type"`
	Token    string `json:"token"`
	ActionTs string `json:"action_ts"`

	Team struct {
		Id     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`

	User struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		TeamId   string `json:"team_id"`
		Name     string `json:"name"`
	} `json:"user"`

	Channel struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	IsEnterpriseInstall bool `json:"is_enterprise_install"`

	// CallbackId name of the shortcut
	CallbackId  string          `json:"callback_id"`
	TriggerId   string          `json:"trigger_id"`
	ResponseUrl string          `json:"response_url"`
	MessageTs   string          `json:"message_ts"`
	Message     ShortcutMessage `json:"message"`
}

func (model MessageShortcutCallBackModel) GetReportType() string {
	switch model.CallbackId {
	case "ops_wtf": return Entry.ReportTypeBadRequest
	case "ops_didnt_read_jenkins_log": return Entry.ReportTypeDidNotReadJenkinsLogs
	}

	return ""
}