package Controllers

import (
	"encoding/json"
	"fmt"
	"github.com/LastSprint/feedback_bot/Steve/Models/DTO"
	"log"
	"net/http"
	"strings"
	"time"
)

type SaWeekStatService interface {
	GatherStatistic() (*DTO.SAWeeklyStat, error)
}

type TimeSpentAnalyticsService interface {
	GetAllIssuesWithTimeSpent(startTimeFrame time.Time, endTimeFrame time.Time, jiraUserNames []string) ([]DTO.UserWorkLog, error)
}

// EventHandlerController handles slack events
type CommandHandlerController struct {
	SaWeekStatService
	TimeSpentAnalyticsService
}

// Init add http handlers for:
//	- `POST /commands/ops/weekly_stat`
//  - `POST /command/analytics/user_project_spent`
func (cnt *CommandHandlerController) Init() {
	http.HandleFunc("/commands/ops/weekly_stat", cnt.handleSaWeeklyStat)
	http.HandleFunc("/command/analytics/user_project_spent", cnt.handleWorkLogAnalytics)
}

// handleSaWeeklyStat for `POST /commands/ops/weekly_stat`
func (cnt *CommandHandlerController) handleSaWeeklyStat(w http.ResponseWriter, _ *http.Request) {
	stat, err := cnt.SaWeekStatService.GatherStatistic()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while gathering sa weekly stat " + err.Error()))
		return
	}

	strBuilder := strings.Builder{}
	strBuilder.WriteString("Всем привет!\nЭта неделя почти закончилась:party_blob:\n\nВот что случилось за это время:\n")
	strBuilder.WriteString(fmt.Sprintf("- Запросов сделано в #devops_and_sa: `%v`", stat.RequestsCount))

	writeReportsIfPossible(&strBuilder, stat)
	writeReactionsIfPossible(&strBuilder, stat)

	writeMessageInPublicVisibility(strBuilder.String(), w)
}

func (cnt *CommandHandlerController) handleWorkLogAnalytics(w http.ResponseWriter, r *http.Request) {
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

	splited := strings.Split(messageText, " ")

	if len(splited) == 0 {
		w.Write([]byte("You should provide usernames"))
		return
	}

	timeFrom, err := time.Parse("2006-01-02", splited[0])
	userNames := splited[1:]
	if err != nil {
		timeFrom = time.Now().Add(time.Hour * 24 * 7 * -1)
		userNames = splited
		log.Printf("[WARN] Error while parsing timeFrom from value %s; Error: %s", splited[0], err.Error())
	}

	timeTo := time.Now()

	res, err := cnt.TimeSpentAnalyticsService.GetAllIssuesWithTimeSpent(timeFrom, timeTo, userNames)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if len(res) == 0 {
		w.Write([]byte("Result is empty, but seems like that there wasn't any error :hm:"))
		return
	}

	builder := strings.Builder{}

	const timeFormat = "2006-01-02"

	builder.WriteString(fmt.Sprintf("Результат за период %s - %s\n", timeFrom.Format(timeFormat), timeTo.Format(timeFormat)))

	for _, it := range res {

		if len(it.WorkLog) == 0 {
			builder.WriteString(fmt.Sprintf("`%s` ничего не затрекал\n", it.UserName))
			continue
		}

		builder.WriteString(fmt.Sprintf("Для `%s`:", it.UserName))

		for _, wl := range it.WorkLog {
			builder.WriteString(fmt.Sprintf("\n- %s: *%v*", wl.ProjectName, formatTime(wl.TimeSpentSecond)))
		}

		builder.WriteString("\n\n")
	}

	writeMessageInPublicVisibility(builder.String(), w)
}

func writeMessageInPublicVisibility(message string, w http.ResponseWriter) {
	value := struct {
		ResponseType   string `json:"response_type"`
		Text           string `json:"text"`
		DeleteOriginal bool   `json:"replace_original"`
	}{
		ResponseType:   "in_channel",
		Text:           message,
		DeleteOriginal: true,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[ERR] Couldn't write json to response")
		return
	}
}

func formatTime(seconds int) string {
	var minutes = seconds / 60

	if minutes >= 60 {
		hours := float32(minutes) / 60.0
		return fmt.Sprintf("%.2f ч.", hours)
	}
	return fmt.Sprintf("%d мин.", minutes)
}

func writeReportsIfPossible(builder *strings.Builder, stat *DTO.SAWeeklyStat) {
	if len(stat.ReportedRequestsCount) == 0 {
		builder.WriteString("\n\nИ ничего не было зарепорчено. Это хорошо или плохо? :hm:")
		return
	}

	builder.WriteString("\n\nКакие запросы были зарепорчены:")

	for key, value := range stat.ReportedRequestsCount {
		nameOfType := ""
		switch key {
		case DTO.ReportTypeBadRequest:
			nameOfType = "Непонятный запрос:"
		case DTO.ReportTypeDidNotReadJenkinsLogs:
			nameOfType = "Не читал(а) логи:"
		default:
			continue
		}
		builder.WriteString(fmt.Sprintf("\n- %s `%v`", nameOfType, value))
	}
}

func writeReactionsIfPossible(builder *strings.Builder, stat *DTO.SAWeeklyStat) {
	if len(stat.Reactions) == 0 {
		builder.WriteString("Никаких реакций не было :(")
		return
	}

	builder.WriteString("\n\nКакие реакции ставили SA:\n")

	for key, val := range stat.Reactions {
		str := fmt.Sprintf(":%s: :%v\n", key, val)
		builder.WriteString(str)
	}
}
