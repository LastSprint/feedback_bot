package Services

import (
	"github.com/LastSprint/feedback_bot/Common/Models/Jira/Tempo"
	"github.com/LastSprint/feedback_bot/Steve/Models/DTO"
	"github.com/pkg/errors"
	"time"
)

type WorkLogRepo interface {
	LoadWorkLogForSpecificUser(jiraUserName []string, from, to time.Time) ([]Tempo.WorkLogModel, error)
}

// TimeSpendDistributionService analyzes time spent distribution for specific users 
type TimeSpendDistributionService struct {
	WorkLogRepo
}

func (srv *TimeSpendDistributionService) GetAllIssuesWithTimeSpent(startTimeFrame time.Time, endTimeFrame time.Time, jiraUserNames []string) ([]DTO.UserWorkLog, error) {
	rawWorkLogs, err := srv.WorkLogRepo.LoadWorkLogForSpecificUser(jiraUserNames, startTimeFrame, endTimeFrame)

	if err != nil {
		return nil, errors.WithMessagef(err, "Couldn't load work log for %s", jiraUserNames)
	}

	if len(rawWorkLogs) == 0 {
		return nil, nil
	}

	// first key is username second key is project id
	userProjectSpent := map[string]map[string]int{}

	for _, it := range rawWorkLogs {

		_, ok := userProjectSpent[it.Worker]

		if !ok {
			// it's new user then we just add his log and continue
			userProjectSpent[it.Worker] = map[string]int{it.Issue.ProjectKey: it.TimeSpentSeconds}
			continue
		}

		_, ok = userProjectSpent[it.Worker][it.Issue.ProjectKey]

		if !ok {
			// user is exist but didn't have this project
			// add info and continue
			userProjectSpent[it.Worker][it.Issue.ProjectKey] = it.TimeSpentSeconds
			continue
		}

		// this project already exist in user projects

		userProjectSpent[it.Worker][it.Issue.ProjectKey] += it.TimeSpentSeconds
	}

	// if user doesn't have a work log, then he will not get into jiraUserNames
	// so we add him manually and show that his work log is empty
	for _, it := range jiraUserNames {
		if _, ok := userProjectSpent[it]; !ok {
			userProjectSpent[it] = map[string]int{}
		}
	}

	result := make([]DTO.UserWorkLog, len(userProjectSpent))

	counter := 0

	for userName, projectsAndSpent := range userProjectSpent {

		projects := make([]DTO.WorkLogInProjects, len(projectsAndSpent))

		i := 0

		for projectKey, spent := range projectsAndSpent {
			projects[i] = DTO.WorkLogInProjects{
				ProjectName:     projectKey,
				TimeSpentSecond: spent,
			}
			i++
		}

		result[counter] = DTO.UserWorkLog{
			UserName: userName,
			WorkLog: projects,
		}
		counter++
	}

	return result, nil
}