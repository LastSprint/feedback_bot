package DTO

type TechnologyWidespreadAnalyseResult struct {
	// ProjectName is name of jira project in which the technology is used
	ProjectName string
	// UsageCount is "how many times the technology were used"
	UsageCount int
	// WidespreadPercent are percents of issues which was done by the technology
	WidespreadPercent float32
}
