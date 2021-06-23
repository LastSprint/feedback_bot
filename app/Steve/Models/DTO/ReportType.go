package DTO

import "github.com/LastSprint/feedback_bot/Steve/Models/Entry"

const (
	// ReportTypeBadRequest must be used for reporting requests without enough information
	ReportTypeBadRequest string = "ops_wtf"
	// ReportTypeDidNotReadJenkinsLogs must be used for reporting such requests when people didn't even read logs. And the want that SA or DevOps read it instead.
	ReportTypeDidNotReadJenkinsLogs string = "ops_didnt_read_jenkins_log"
)

// ReportTypeFromEntry converts `Entry.ReportType` to `DTO.ReportType`
func ReportTypeFromEntry(val string) string {
	switch val {
	case Entry.ReportTypeBadRequest: return ReportTypeBadRequest
	case Entry.ReportTypeDidNotReadJenkinsLogs: return ReportTypeDidNotReadJenkinsLogs
	}

	return ""
}