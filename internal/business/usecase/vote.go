package usecase

import (
	"errors"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"github.com/rs/zerolog/log"
	"time"
)

func (uu *UserUsecase) GetRating(target string) (*model.Rating, error) {
	rating, err := uu.vr.CalcTotalRating(target)
	return rating, err
}

func (uu *UserUsecase) RetractVote(u *model.Update, sender string) error {
	var vote = &model.Vote{
		Sender: sender,
		Target: u.Nickname,
	}

	vote, err := uu.vr.Find(vote, true, true)
	if err != nil {
		return errors.New("no vote has been recorded")
	}

	if err = uu.vr.Delete(vote, true, true); err != nil {
		return err
	}

	return nil
}

func (uu *UserUsecase) StoreVote(v *model.Vote) (*repository.VoteInsertResult, error) {
	lastStoredVote, err := uu.vr.Find(v, true, true)
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

	result, err := uu.vr.Create(v)
	return result, err
}

func validateTimespanBetweenVotes(newVote *model.Vote, lastVote *model.Vote) error {
	timeOfNewVote := time.Unix(*newVote.VotedAt, 0)
	timeOfLastVote := time.Unix(*lastVote.VotedAt, 0)

	timeDiff := timeOfNewVote.Sub(timeOfLastVote)

	if timeDiff.Hours() < 1 {
		return errors.New("voting allowed only once per hour")
	}

	return nil
}
