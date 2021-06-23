package Tempo

type WorkLogModel struct {
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	TempoWorkLogId   int    `json:"tempoWorklogId"`
	Issue            Issue  `json:"issue"`
	Comment          string `json:"comment"`
	Worker           string `json:"worker"`
	Updater          string `json:"updater"`
}
