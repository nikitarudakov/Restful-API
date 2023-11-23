package vote

import (
	"context"
	"errors"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/controller"
)

type GRPCVoteService struct {
	VoteUsecase controller.VoteManager
	CacheDB     *cache.Database
	UnimplementedVoteControllerServer
}

func (vs *GRPCVoteService) Vote(ctx context.Context, vote *VoteArg) (*Empty, error) {
	if err := validateVote(vote); err != nil {
		return &Empty{}, err
	}

	voteConverted := &model.Vote{
		Sender:  vote.Sender,
		Target:  vote.Target,
		Vote:    vote.Vote,
		VotedAt: vote.VotedAt,
	}

	_, err := vs.VoteUsecase.StoreVote(voteConverted)
	if err != nil {
		return &Empty{}, err
	}

	return &Empty{}, nil
}

func (vs *GRPCVoteService) RetractVote(ctx context.Context, retractArg *RetractArg) (*Empty, error) {
	var update = &model.Update{Nickname: retractArg.Target}

	if err := vs.VoteUsecase.RetractVote(update, retractArg.Sender); err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

func (vs *GRPCVoteService) GetRating(ctx context.Context, target *Target) (*Rating, error) {
	rating, err := vs.VoteUsecase.GetRating(target.Target)
	if err != nil {
		return nil, err
	}

	ratingConverted := Rating{Target: rating.Target, Rating: rating.Rating}

	return &ratingConverted, nil
}

func validateVote(v *VoteArg) error {
	sender := v.Sender
	target := v.Target

	if sender == target {
		return errors.New("self-voting is forbidden (bad request)")
	}
	vote := v.Vote
	if vote != -1 && vote != 1 {
		return errors.New("invalid vote")
	}

	return nil
}
