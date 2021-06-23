package Repo

import (
	"bytes"
	"encoding/json"
	"github.com/LastSprint/feedback_bot/Common/Models/Jira/Tempo"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// JiraTempoTimesheetRepo is a repo which works with Tempo Timesheet REST API (server installation)
type JiraTempoTimesheetRepo struct {
	// ApiToken is a token for accessing jira REST API (basic auth)
	ApiToken string
}

const requestTimeFormat = "2006-01-02"

func (repo *JiraTempoTimesheetRepo) LoadWorkLogForSpecificUser(jiraUserName []string, from, to time.Time) ([]Tempo.WorkLogModel, error) {

	reqBody := struct {
		From string `json:"from"`
		To string `json:"to"`
		Worker []string `json:"worker"`
	}{
		From: from.Format(requestTimeFormat),
		To: to.Format(requestTimeFormat),
		Worker: jiraUserName,
	}

	body, err := json.Marshal(reqBody)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://jira.surfstudio.ru/rest/tempo-timesheets/4/worklogs/search", bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	headers := req.Header

	headers.Set("Authorization", "Basic " + repo.ApiToken)
	headers.Set("Content-Type", "Application/json")

	req.Header = headers

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.WithMessage(err, "Couldn't read response body for " + req.URL.String())
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("[WARN] Couldn't close response body in request %s with error %s", req.URL.String(), err.Error())
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Got %v status code on %s request. Response body: %s", response.StatusCode, req.URL.String(), string(respBody))
	}

	result := []Tempo.WorkLogModel{}

	if err = json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.WithMessagef(err, "Couldn't parse response on %s with error", req.URL.String())
	}

	return result, nil
}