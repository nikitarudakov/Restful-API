package usecase

import (
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"github.com/rs/zerolog/log"
	"time"
)

type VoteHistory struct {
	hasVote bool
	index   *int
	vote    *model.Vote
	history []*model.Vote
}

type HistoryOptions struct {
	add           bool
	voteToCompare *model.Vote
}

func (vh *VoteHistory) removeVoteFromHistory() {
	history := vh.history
	index := *vh.index

	if len(history) == 1 || vh.index == nil {
		vh.history = nil
		return
	}

	history[index] = history[len(history)-1]
	history[len(history)-1] = &model.Vote{}

	historyTruncated := history[:len(history)-1]

	vh.history = historyTruncated
}

func (uu *UserUsecase) RetractVote(uSender *model.Update, uTarget *model.Update) error {
	vote, err := uu.retractVoteSender(uSender, uTarget.Nickname)
	if err != nil {
		return err
	}

	return uu.retractVoteTarget(uTarget, vote)
}

func (uu *UserUsecase) retractVoteSender(uSender *model.Update, targetUsername string) (*model.Vote, error) {
	var pSender model.Profile
	pSender.Nickname = uSender.Nickname

	senderFromDB, err := uu.pr.Find(&pSender)
	if err != nil {
		return nil, err
	}

	vh, err := findVoteInVoteHistory(senderFromDB, targetUsername)
	if err != nil {
		return nil, err
	}

	if !vh.hasVote {
		return nil, errors.New("error retracting non-existent vote")
	}

	vh.removeVoteFromHistory()

	uSender.Votes = vh.history

	if err = uu.pr.Update(uSender, pSender.Nickname); err != nil {
		return nil, err
	}

	return vh.vote, nil
}

func (uu *UserUsecase) retractVoteTarget(uTarget *model.Update, vote *model.Vote) error {
	var pTarget model.Profile
	pTarget.Nickname = uTarget.Nickname

	targetFromDB, err := uu.pr.Find(&pTarget)
	if err != nil {
		return err
	}

	uTarget.Rating = targetFromDB.Rating

	if uTarget.Rating == nil {
		return errors.New("error target has no voting records")
	}

	*uTarget.Rating -= vote.Vote

	return uu.pr.Update(uTarget, pTarget.Nickname)
}

func findVoteInVoteHistory(subj *model.Profile, target string, opts ...HistoryOptions) (*VoteHistory, error) {
	var voteHistory VoteHistory
	var historyOptions *HistoryOptions

	if len(subj.Votes) == 0 {
		return &VoteHistory{hasVote: false}, nil
	}

	if len(opts) > 0 {
		historyOptions = &opts[0]
	}

	voteHistory.history = subj.Votes

	for i, v := range subj.Votes {
		if historyOptions != nil && historyOptions.add {
			if v.Target == target && v.Vote != historyOptions.voteToCompare.Vote {
				return &VoteHistory{hasVote: true, index: &i, vote: v, history: subj.Votes}, nil
			}

			if v.Target == target && v.Vote == historyOptions.voteToCompare.Vote {
				return nil, errors.New("error: voting twice for same target is forbidden")
			}

			continue
		}

		if v.Target == target {
			return &VoteHistory{hasVote: true, index: &i, vote: v, history: subj.Votes}, nil
		}
	}

	return &VoteHistory{hasVote: false}, nil
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

	senderFromDB, err := uu.pr.Find(&p)
	if err != nil {
		return err
	}

	opts := HistoryOptions{add: true, voteToCompare: vote}

	voteHistory, err := findVoteInVoteHistory(senderFromDB, vote.Target, opts)
	if err != nil {
		return err
	}

	log.Trace().Str("vote_history", fmt.Sprintf("%v", voteHistory)).Send()

	var history []*model.Vote
	if voteHistory.hasVote {
		history = voteHistory.history
		history[*voteHistory.index] = vote
	} else {
		history = append(senderFromDB.Votes, vote)
	}

	if err = validateTimespanBetweenVotes(history, vote); err != nil {
		return err
	}

	u.Votes = history

	return uu.pr.Update(u, p.Nickname)
}

func validateTimespanBetweenVotes(profileVotes []*model.Vote, newVote *model.Vote) error {
	timeOfNewVote := time.Unix(*newVote.VotedAt, 0)

	lastVote := profileVotes[len(profileVotes)-1]

	timeOfLastVote := time.Unix(*lastVote.VotedAt, 0)

	timeDiff := timeOfNewVote.Sub(timeOfLastVote)

	fmt.Println(timeDiff)

	if timeDiff.Hours() < 1 {
		return errors.New("voting allowed once per hour")
	}

	return nil
}
