package usecase

import (
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"time"
)

type VoteHistory struct {
	has     bool
	index   int
	history []model.Vote
}

func (uu *UserUsecase) RetractVote(uTarget *model.Update, uSender *model.Update) error {
	var pTarget model.Profile
	pTarget.Nickname = uTarget.Nickname

	targetFromDB, err := uu.pr.Find(&pTarget)
	if err != nil {
		return err
	}

	// get current target rating
	uTarget.Rating = targetFromDB.Rating

	var pSender model.Profile
	pSender.Nickname = uSender.Nickname

	senderFromDB, err := uu.pr.Find(&pSender)
	if err != nil {
		return err
	}

	var voteHistory VoteHistory
	voteHistory.history = senderFromDB.Votes

	for i, v := range voteHistory.history {
		if v.Target == uTarget.Nickname {
			if uTarget.Rating == nil {
				return errors.New("this target has no rating record")
			}

			voteHistory.has = true
			voteHistory.index = i
		}
	}

	if !voteHistory.has {
		return errors.New("error retracting non-existent vote")
	}

	fmt.Println(voteHistory.history)

	*uTarget.Rating -= voteHistory.history[voteHistory.index].Vote

	err = uu.pr.Update(uTarget, pTarget.Nickname)
	if err != nil {
		return err
	}

	voteHistory.removeVoteFromHistory()

	uSender.Votes = voteHistory.history

	err = uu.pr.Update(uSender, pSender.Nickname)
	if err != nil {
		return err
	}

	return nil
}

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

	if len(profileFromDB.Votes) > 1 {
		if err = validateTimeSpanBetweenVotes(profileFromDB.Votes, vote); err != nil {
			return err
		}
	}

	u.Votes = votes

	return uu.pr.Update(u, p.Nickname)
}

func validateTimeSpanBetweenVotes(profileVotes []model.Vote, newVote *model.Vote) error {
	timeOfNewVote := time.Unix(*newVote.VotedAt, 0)

	lastVote := profileVotes[len(profileVotes)-1]

	timeOfLastVote := time.Unix(*lastVote.VotedAt, 0)

	timeDiff := timeOfNewVote.Sub(timeOfLastVote)

	if timeDiff.Hours() < 1 {
		return errors.New("voting allowed once per hour")
	}

	return nil
}

func (vh *VoteHistory) removeVoteFromHistory() {
	if len(vh.history) == 1 {
		vh.history = []model.Vote{}
		return
	}

	vh.history[vh.index] = vh.history[len(vh.history)-1]
	vh.history[len(vh.history)-1] = model.Vote{}
	vh.history = vh.history[:len(vh.history)-1]
}
