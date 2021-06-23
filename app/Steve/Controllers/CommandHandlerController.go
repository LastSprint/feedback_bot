package Controllers

import (
	"encoding/json"
	"fmt"
	"github.com/LastSprint/feedback_bot/Steve/Models/DTO"
	"log"
	"net/http"
	"strings"
)

type SaWeekStatService interface {
	GatherStatistic() (*DTO.SAWeeklyStat, error)
}

// EventHandlerController handles slack events
type CommandHandlerController struct {
	SaWeekStatService
}

// Init add http handlers for:
//	- `POST /commands/ops/weekly_stat`
func (cnt *CommandHandlerController) Init() {
	http.HandleFunc("/commands/ops/weekly_stat", cnt.handleSaWeeklyStat)
}

func (cnt *CommandHandlerController) handleSaWeeklyStat(w http.ResponseWriter, r *http.Request) {
	stat, err := cnt.SaWeekStatService.GatherStatistic()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while gathering sa weekly stat " + err.Error()))
		return
	}

	strBuilder := strings.Builder{}
	strBuilder.WriteString("Всем привет!\nЭта неделя почти закончилась:party_blob:\n\nВот что случилось за это время:\n")
	strBuilder.WriteString(fmt.Sprintf("- Запросов сделано в #devops_and_sa: `%v`", stat.RequestsCount))
	if len(stat.ReportedRequestsCount) == 0 {
		strBuilder.WriteString("\n\nИ ничего не было зарепорчено. Это хорошо или плохо? :hm:")
		writeMessageInPublicVisibility(strBuilder.String(), w)
		return
	}

	strBuilder.WriteString("\n\nКакие запросы были зарепорчены:")

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
		strBuilder.WriteString(fmt.Sprintf("\n- %s `%v`", nameOfType, value))
	}

	writeMessageInPublicVisibility(strBuilder.String(), w)
}

func writeMessageInPublicVisibility(message string, w http.ResponseWriter) {
	value := struct {
		ResponseType string `json:"response_type"`
		Text         string `json:"text"`
		DeleteOriginal bool `json:"replace_original"`
	}{
		ResponseType: "in_channel",
		Text: message,
		DeleteOriginal: true,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[ERR] Couldn't write json to response")
		return
	}
}
