package usecase

import (
	"errors"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type VoteUsecase struct {
	pr repository.ProfileRepoManager
	vr repository.VoteRepoManager
}

func NewVoteUsecase(pr repository.ProfileRepoManager, vr repository.VoteRepoManager) *VoteUsecase {
	return &VoteUsecase{pr: pr, vr: vr}
}

func (vc *VoteUsecase) GetRating(target string) (*model.Rating, error) {
	rating, err := vc.vr.GetRating(target)
	return rating, err
}

func (vc *VoteUsecase) RetractVote(u *model.Update, sender string) error {
	var vote = &model.Vote{
		Sender: sender,
		Target: u.Nickname,
	}

	vote, err := vc.vr.Find(vote, true, true)
	if err != nil {
		return errors.New("no vote has been recorded")
	}

	return vc.vr.Delete(vote, true, true)
}

func (vc *VoteUsecase) StoreVote(v *model.Vote) (*repository.VoteInsertResult, error) {
	lastStoredVote, err := vc.vr.Find(v, true, true)
	if err != nil {
		log.Warn().Str("service", "last stored vote").Err(err).Send()
	}

	if lastStoredVote != nil {
		if v.Vote == lastStoredVote.Vote {
			return nil, errors.New("voting for same user twice is not allowed")
		}

		if err = validateTimespanBetweenVotes(v, lastStoredVote); err != nil {
			return nil, err
		}
	}

	result, err := vc.vr.Create(v)
	return result, err
}

func validateTimespanBetweenVotes(newVote *model.Vote, lastVote *model.Vote) error {
	timeOfNewVote := time.Unix(newVote.VotedAt, 0)
	timeOfLastVote := time.Unix(lastVote.VotedAt, 0)

	timeDiff := timeOfNewVote.Sub(timeOfLastVote)

	if timeDiff.Hours() < 1 {
		return errors.New("voting allowed only once per hour")
	}

	return nil
}
