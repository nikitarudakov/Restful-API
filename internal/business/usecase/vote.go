package usecase

import (
	"context"
	"errors"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/voteDao"
	"github.com/rs/zerolog/log"
	"time"
)

type VoteUsecase struct {
	grpcClient voteDao.VoteRepoClient
}

func NewVoteUsecase(client voteDao.VoteRepoClient) *VoteUsecase {
	return &VoteUsecase{grpcClient: client}
}

func (vc *VoteUsecase) GetRating(target string) (*model.Rating, error) {
	gRPCRating, err := vc.grpcClient.GetRating(context.Background(), &voteDao.Target{Val: target})
	if err != nil {
		return nil, err
	}

	return voteDao.UnmarshalToRating(gRPCRating), err
}

func (vc *VoteUsecase) RetractVote(u *model.Update, sender string) error {
	var vote = &model.Vote{
		Sender: sender,
		Target: u.Nickname,
	}

	gRPCVoteFilter := voteDao.MarshalTogRPCVoteFilter(vote, true, true)

	_, err := vc.grpcClient.FindVoteInStorage(context.Background(), gRPCVoteFilter)
	if err != nil {
		return errors.New("no vote has been recorded")
	}

	_, err = vc.grpcClient.DeleteVoteFromStorage(context.Background(), gRPCVoteFilter)
	return err
}

func (vc *VoteUsecase) StoreVote(v *model.Vote) (*voteDao.InsertResult, error) {
	gRPCVoteFilter := voteDao.MarshalTogRPCVoteFilter(v, true, true)

	lastStoredVote, err := vc.grpcClient.FindVoteInStorage(context.Background(), gRPCVoteFilter)
	if err != nil {
		log.Warn().Str("service", "FindVoteInStorage").Err(err).Send()
	}

	if lastStoredVote != nil {
		if v.Vote == lastStoredVote.Vote {
			return nil, errors.New("voting for same user twice is not allowed")
		}

		if err = validateTimespanBetweenVotes(v, lastStoredVote); err != nil {
			return nil, err
		}
	}

	gRPCVote := voteDao.MarshalTogRPCVote(v)

	insertResult, err := vc.grpcClient.InsertVoteToStorage(context.Background(), gRPCVote)

	return insertResult, err
}

func validateTimespanBetweenVotes(newVote *model.Vote, lastVote *voteDao.Vote) error {
	timeOfNewVote := time.Unix(newVote.VotedAt, 0)
	timeOfLastVote := time.Unix(lastVote.VotedAt, 0)

	timeDiff := timeOfNewVote.Sub(timeOfLastVote)

	if timeDiff.Hours() < 1 {
		return errors.New("voting allowed only once per hour")
	}

	return nil
}
