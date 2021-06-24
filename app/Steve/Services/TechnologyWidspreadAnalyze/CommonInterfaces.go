package TechnologyWidspreadAnalyze

import (
	"github.com/LastSprint/JiraGoIssues/models"
	"github.com/LastSprint/JiraGoIssues/services"
)

type JiraRepo interface {
	LoadIssues(convertible services.RequestConvertible) (models.IssueSearchWrapperEntity, error)
}