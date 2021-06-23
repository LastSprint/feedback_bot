package Entry

const (
	// ReportTypeBadRequest must be used for reporting requests without enough information
	ReportTypeBadRequest string = "bad_request"
	// ReportTypeDidNotReadJenkinsLogs must be used for reporting such requests when people didn't even read logs. And the want that SA or DevOps read it instead.
	ReportTypeDidNotReadJenkinsLogs string = "did_not_read_jenkins_logs"
)