package Tempo

type Issue struct {
	ReporterKey               string   `json:"reporterKey"`
	IssueStatus               string   `json:"issueStatus"`
	OriginalEstimateSeconds   int      `json:"originalEstimateSeconds"`
	EstimatedRemainingSeconds int      `json:"estimatedRemainingSeconds"`
	InternalIssue             bool     `json:"internalIssue"`
	Summary                   string   `json:"summary"`
	IssueType                 string   `json:"issueType"`
	ProjectId                 int      `json:"projectId"`
	//Components                []string `json:"components"`
	ProjectKey                string   `json:"projectKey"`
	IconUrl                   string   `json:"iconUrl"`
	Versions                  []string `json:"versions"`
	Key                       string   `json:"key"`
	ID                        int      `json:"id"`
}
