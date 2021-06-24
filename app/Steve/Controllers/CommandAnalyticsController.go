package Controllers

import (
	"fmt"
	"github.com/LastSprint/feedback_bot/Steve/Models/DTO"
	"net/http"
	"strings"
)

type TechWidespreadAnalyze interface {
	Calculate() ([]DTO.TechnologyWidespreadAnalyseResult, error)
}

// CommandAnalyticsController contains endpoints for analytics
type CommandAnalyticsController struct {
	Templates TechWidespreadAnalyze
	SurfGen TechWidespreadAnalyze
}

func (cnt *CommandAnalyticsController) Init() {
	 http.HandleFunc("/commands/analyzes/tech_widespread", cnt.handleTechWidespreadAnalyse)
}

func (cnt *CommandAnalyticsController) handleTechWidespreadAnalyse(w http.ResponseWriter, r *http.Request) {
	res, err := cnt.Templates.Calculate()

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if len(res) == 0 {
		w.Write([]byte("Result is empty. Seems like that templates aren't used"))
		return
	}

	templatesMsg := analysisResultToFormattedText(res, "Templates")

	surfGenResult, err := cnt.SurfGen.Calculate()

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	if len(surfGenResult) == 0 {
		w.Write([]byte("Result is empty. Seems like that templates aren't used"))
		return
	}

	surfGenMsg := analysisResultToFormattedText(surfGenResult, "SurfGen")

	writeMessageInPublicVisibility(templatesMsg + surfGenMsg, w)
}

func analysisResultToFormattedText(data []DTO.TechnologyWidespreadAnalyseResult, header string) string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("*%s*\n", header))

	for _, it := range data {
		builder.WriteString(fmt.Sprintf("%s:\n", it.ProjectName))
		builder.WriteString(fmt.Sprintf("Использовалось раз: %v\n", it.UsageCount))
		builder.WriteString(fmt.Sprintf("Процент внедрения: %.5f", it.WidespreadPercent))
		builder.WriteString("\n\n")
	}
	return builder.String()
}