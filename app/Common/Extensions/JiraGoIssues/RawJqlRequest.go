package JiraGoIssues

import "github.com/LastSprint/JiraGoIssues/services"

// RawJqlRequest is a simple Jira search request which you can build with raw JQL
// Implements JiraGoIssues.service.RequestConvertible
type RawJqlRequest struct {
	JQL string
	AdditionalFields []services.JiraField
	PageSize int
}

func (r RawJqlRequest) GetUseOnlyAdditionalFields() bool {
	return false
}

func (r RawJqlRequest) MakeJiraRequest() string {
	return r.JQL
}

func (r RawJqlRequest) GetAdditionFields() []services.JiraField {
	return r.AdditionalFields
}

func (r RawJqlRequest) GetPageSize() int {
	return r.PageSize
}

