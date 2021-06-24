package TechnologyWidspreadAnalyze

import (
	"fmt"
	"github.com/LastSprint/JiraGoIssues/models"
	"github.com/LastSprint/feedback_bot/Common/Extensions/JiraGoIssues"
	"github.com/LastSprint/feedback_bot/Steve/Models/DTO"
	"github.com/pkg/errors"
	"log"
	"sync"
)

type TemplatesService struct {
	JiraRepo
}

func (srv *TemplatesService) Calculate() ([]DTO.TechnologyWidespreadAnalyseResult, error) {
	res, err := srv.LoadIssues(JiraGoIssues.RawJqlRequest{
		JQL:              "labels in (\"templates\")",
		AdditionalFields: nil,
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

	projectAndCount := map[string][]models.IssueEntity{}
	rwMutex := &sync.RWMutex{}
	wg := &sync.WaitGroup{}
	totalIssuesCountInProject := map[string]int{}

	for _, it := range res.Issues {
		_, ok := projectAndCount[it.Fields.Project.Key]

		if !ok {
			wg.Add(1)
			go func(projectId string) {
				res, err := srv.LoadIssues(JiraGoIssues.RawJqlRequest{
					JQL:              fmt.Sprintf("project = %s and type not in (Epic, Test, \"Service Task\")", projectId),
					AdditionalFields: nil,
					PageSize:         1,
				})

				// TODO: Handle error in host code (throw it)
				if err != nil {
					log.Println("[ERR] Can't load project info for ")
				}

				rwMutex.Lock()
				totalIssuesCountInProject[projectId] = res.Total
				rwMutex.Unlock()
				wg.Done()
			}(it.Fields.Project.Key)

			projectAndCount[it.Fields.Project.Key] = []models.IssueEntity{it}
			continue
		}
		projectAndCount[it.Fields.Project.Key] = append(projectAndCount[it.Fields.Project.Key], it)
	}

	wg.Wait()

	//if len(projectAndCount) != len(totalIssuesCountInProject) {
	//	return nil, errors.New("Couldn't load project info from jira. For more info look inside logs")
	//}

	result := make([]DTO.TechnologyWidespreadAnalyseResult, len(projectAndCount))

	i := 0
	for project, count := range projectAndCount {

		projectTotal, ok := totalIssuesCountInProject[project]

		if !ok {
			return nil, errors.Errorf("%s wasn't loaded from Jira. Look at logs for details", project)
		}

		result[i] = DTO.TechnologyWidespreadAnalyseResult{
			ProjectName:       project,
			UsageCount:        len(count),
			WidespreadPercent: float32(len(count)) / float32(projectTotal),
		}
		i++
	}

	return result, nil
}
