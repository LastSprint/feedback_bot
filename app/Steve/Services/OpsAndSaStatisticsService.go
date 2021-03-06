package Services

import (
	"github.com/LastSprint/feedback_bot/Steve/Models/DTO"
	"github.com/pkg/errors"
	"log"
	"sync"
)

type SAStatisticRequestRepo interface {
	// GetCountForThisWeek returns just a number of requests which was written in public chat for SA requests for the week
	GetCountForThisWeek(channelID string) (int, error)
}

type SAStatReportsRepo interface {
	// GetCountForThisWeek will return a map where keys is from `Entry.ReportType` and values is a number of reports
	// so it returns a number of each type of reports
	GetCountForThisWeek(channelID string) (map[string]int, error)
}

type SAReactionsRepo interface {
	ReadReactions() (map[string]int, error)
}

// OpsAndSaStatisticsService get statistic about SA work and send it to specific slack channel
//
// What kind of statistics:
//	1. Count of requests in public chat
//  2. Count of reported requests
type OpsAndSaStatisticsService struct {
	SAStatisticRequestRepo
	SAStatReportsRepo
	SAReactionsRepo
	PublicSaRequestsChannelId string
}

func (srv *OpsAndSaStatisticsService) GatherStatistic() (*DTO.SAWeeklyStat, error) {

	var stat *DTO.SAWeeklyStat
	var statErr error

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		stat, statErr = srv.gatherReports()
		wg.Done()
	}()

	wg.Add(1)

	var reactions map[string]int
	var reactionsError error

	go func() {
		reactions, reactionsError = srv.gatherReactions()
		wg.Done()
	}()

	wg.Wait()

	if statErr != nil {
		return nil, statErr
	}

	if reactionsError != nil {
		return nil, reactionsError
	}

	stat.Reactions = reactions

	return stat, nil
}

func (srv *OpsAndSaStatisticsService) gatherReports() (*DTO.SAWeeklyStat, error) {
	requestsCount, err := srv.SAStatisticRequestRepo.GetCountForThisWeek(srv.PublicSaRequestsChannelId)
	if err != nil {
		return nil, errors.WithMessage(err, "request count repo")
	}

	reports, err := srv.SAStatReportsRepo.GetCountForThisWeek(srv.PublicSaRequestsChannelId)

	if err != nil {
		return nil, errors.WithMessage(err, "reports count repo")
	}

	convertedReport := map[string]int{}

	for key, val := range reports {
		converted := DTO.ReportTypeFromEntry(key)

		if len(converted) == 0 {
			log.Printf("[ERR] couldn't convert report type from entry %s to dto", key)
			continue
		}

		convertedReport[DTO.ReportTypeFromEntry(key)] = val
	}

	return &DTO.SAWeeklyStat{
		RequestsCount:         requestsCount,
		ReportedRequestsCount: convertedReport,
	}, nil
}

func (srv *OpsAndSaStatisticsService) gatherReactions() (map[string]int, error) {
	return srv.SAReactionsRepo.ReadReactions()
}
