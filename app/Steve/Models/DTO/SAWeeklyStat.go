package DTO

// SAWeeklyStat contains statistic about week of work for SA and Ops
type SAWeeklyStat struct {
	// RequestsCount a number of requests which was sent to SA requests chat
	RequestsCount int
	// ReportedRequestsCount is a map where key is type of reported request and value is a number of reports
	ReportedRequestsCount map[string]int
	// Reactions contains count of each reaction
	Reactions map[string]int
}
