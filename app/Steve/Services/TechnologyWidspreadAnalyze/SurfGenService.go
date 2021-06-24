package TechnologyWidspreadAnalyze

import (
	"fmt"
	"github.com/LastSprint/JiraGoIssues/models"
	"github.com/LastSprint/JiraGoIssues/services"
	"github.com/LastSprint/feedback_bot/Common/Extensions/JiraGoIssues"
	"github.com/LastSprint/feedback_bot/Common/Utils"
	"github.com/LastSprint/feedback_bot/Steve/Models/DTO"
	"github.com/pkg/errors"
	"log"

)

const serviceIssueLabel = "service"
const surfGenIssueLabel = "surfgen"

type SurfGenService struct {
	JiraRepo
}

func (srv *SurfGenService) Calculate() ([]DTO.TechnologyWidespreadAnalyseResult, error) {
	res, err := srv.LoadIssues(JiraGoIssues.RawJqlRequest{
		JQL:              fmt.Sprintf("labels in (\"%s\")", serviceIssueLabel),
		AdditionalFields: []services.JiraField{services.JiraField("labels")},
		PageSize:         5000,
	})

	if err != nil {
		return nil, errors.WithMessage(err, "Couldn't load issues with label `templates`")
	}

	if res.Total == 0 {
		return nil, nil
	}

	if len(res.Issues) != res.Total {
		log.Printf("[WARN] When load `templates` issue loaded %v but total count is %v", len(res.Issues), res.Total)
	}

	projectAndIssues := map[string][]models.IssueEntity{}

	for _, it := range res.Issues {
		_, ok := projectAndIssues[it.Fields.Project.Key]

		if !ok {
			projectAndIssues[it.Fields.Project.Key] = []models.IssueEntity{it}
			continue
		}
		projectAndIssues[it.Fields.Project.Key] = append(projectAndIssues[it.Fields.Project.Key], it)
	}

	result := make([]DTO.TechnologyWidespreadAnalyseResult, len(projectAndIssues))

	i := 0
	for project, issues := range projectAndIssues {

		totalServicesCount := 0
		serviceWithSurfGen := 0

		for _, it := range issues {
			totalServicesCount++

			if Utils.Contains(it.Fields.Labels, surfGenIssueLabel) {
				serviceWithSurfGen += 1
			}
		}

		result[i] = DTO.TechnologyWidespreadAnalyseResult{
			ProjectName:       project,
			UsageCount:        serviceWithSurfGen,
			WidespreadPercent: float32(serviceWithSurfGen) / float32(totalServicesCount),
		}
		i++
	}

	return result, nil
}
