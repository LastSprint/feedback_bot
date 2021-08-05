package Services

import (
	"fmt"
	"reflect"
)

type AddReactionServiceModel struct {
	Reaction       string
	Channel        string
	MessageId      string
	ReporterUserId string
}

type MsgReactionRepo interface {
	AddReactionIfNotAddedPreviously(reaction, channelId, messageId string) error
	RemoveReactionIfPossible(reaction, channelId, messageId string) error
}

type ReactionService struct {
	MsgReactionRepo
}

const (
	saUser         = "U01CKPEPK2Q"
	superAdminUser = "UFH46AX6W"
	devOpsReaction = "devops"
	jarvisReaction = "jarvis"
)

func (srv *ReactionService) Add(model AddReactionServiceModel) error {

	if model.ReporterUserId != saUser && model.ReporterUserId != superAdminUser {
		return fmt.Errorf("only SA's reactions counted and other's are skipped")
	}

	switch model.Reaction {
	case devOpsReaction, jarvisReaction:
		break
	default:
		return fmt.Errorf("reaction %s coudn't be handled by %s because is unknown reaction", model.Reaction, reflect.TypeOf(srv))
	}

	if err := srv.MsgReactionRepo.AddReactionIfNotAddedPreviously(model.Reaction, model.Channel, model.MessageId); err != nil {
		return fmt.Errorf("can't write reaction to repo with error\n%w", err)
	}

	return nil
}

func (srv *ReactionService) Remove(model AddReactionServiceModel) error {

	if model.ReporterUserId != saUser && model.ReporterUserId != superAdminUser {
		return fmt.Errorf("only SA's reactions counted and other's are skipped")
	}

	switch model.Reaction {
	case devOpsReaction, jarvisReaction:
		break
	default:
		return fmt.Errorf("reaction %s coudn't be handled by %s because is unknown reaction", model.Reaction, reflect.TypeOf(srv))
	}

	if err := srv.MsgReactionRepo.RemoveReactionIfPossible(model.Reaction, model.Channel, model.MessageId); err != nil {
		return fmt.Errorf("can't delete reaction from repo with error\n%w", err)
	}

	return nil
}
