package voteDao

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
)

type VoteRepoController interface {
	FindVoteInStorage(v *Vote, isTarget bool, isSender bool) (*model.Vote, error)
	InsertVoteToStorage(v *Vote) (*model.InsertResult, error)
	DeleteVoteFromStorage(v *Vote, isTarget bool, isSender bool) error
	GetRating(target string) (*model.Rating, error)
}

type Server struct {
	VoteRepo VoteRepoController
	UnimplementedVoteRepoServer
}

func (s *Server) FindVoteInStorage(_ context.Context, voteFilter *VoteFilter) (*Vote, error) {
	vote, err := s.VoteRepo.FindVoteInStorage(voteFilter.Vote,
		voteFilter.IsTarget, voteFilter.IsTarget)

	if err != nil {
		return nil, err
	}

	voteConverted := MarshalTogRPCVote(vote)

	return voteConverted, err
}

func (s *Server) InsertVoteToStorage(_ context.Context, vote *Vote) (*InsertResult, error) {
	insertResult, err := s.VoteRepo.InsertVoteToStorage(vote)
	if err != nil {
		return nil, err
	}

	insertResultConverted := MarshalTogRPCInsertResult(insertResult)

	return insertResultConverted, err
}

func (s *Server) DeleteVoteFromStorage(_ context.Context, vote *VoteFilter) (*Empty, error) {
	err := s.VoteRepo.DeleteVoteFromStorage(vote.Vote, vote.IsTarget, vote.IsSender)
	return nil, err
}

func (s *Server) GetRating(_ context.Context, target *Target) (*Rating, error) {
	rating, err := s.VoteRepo.GetRating(target.Val)
	if err != nil {
		return nil, err
	}

	ratingConverted := MarshalTogRPCRating(rating)

	return ratingConverted, err
}
