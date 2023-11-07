package usecase

import (
	"errors"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
)

func (uu *UserUsecase) AddUpdateVoteTarget(u *model.Update) error {
	var p model.Profile
	p.Nickname = u.Nickname

	profileFromDB, err := uu.pr.Find(&p)
	if err != nil {
		return err
	}

	if profileFromDB.Rating != nil {
		*u.Rating += *profileFromDB.Rating
	}

	return uu.pr.Update(u, p.Nickname)
}

type VoteHistory struct {
	has   bool
	index int
}

func (uu *UserUsecase) AddUpdateVoteSender(u *model.Update, vote *model.Vote) error {
	var p model.Profile
	p.Nickname = u.Nickname

	profileFromDB, err := uu.pr.Find(&p)
	if err != nil {
		return err
	}

	var voteHistory VoteHistory

	for i, v := range profileFromDB.Votes {
		if vote.Target == v.Target && vote.Vote == v.Vote {
			return errors.New("vote has already been stored for this target")
		} else if vote.Target == v.Target && vote.Vote != v.Vote {
			voteHistory = VoteHistory{has: true, index: i}
		}
	}

	var votes []model.Vote

	if voteHistory.has {
		votes = profileFromDB.Votes
		votes[voteHistory.index] = *vote

	} else {
		votes = append(profileFromDB.Votes, *vote)
	}

	/*	if len(profileFromDB.Votes) > 1 {
		timeOfNewVote := time.Unix(*vote.VotedAt, 0)

		lastVote := profileFromDB.Votes[len(profileFromDB.Votes)-1]

		timeOfLastVote := time.Unix(*lastVote.VotedAt, 0)

		timeDiff := timeOfNewVote.Sub(timeOfLastVote)

		if timeDiff.Hours() < 1 {
			return errors.New("voting allowed once per hour")
		}
	}*/

	u.Votes = votes

	return uu.pr.Update(u, p.Nickname)
}
