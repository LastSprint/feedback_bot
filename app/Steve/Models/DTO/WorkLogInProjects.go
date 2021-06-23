package DTO

// WorkLogInProjects contains work log distributed between projects
// So each projects time which was spent on it
type WorkLogInProjects struct {
	ProjectName string
	TimeSpentSecond int
}