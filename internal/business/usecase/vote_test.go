package usecase

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/voteDao"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"testing"
	"time"
)

type VoteRepoMock struct {
	mock.Mock
}

func (s *VoteRepoMock) FindVoteInStorage(ctx context.Context, in *voteDao.VoteFilter, opts ...grpc.CallOption) (*voteDao.Vote, error) {
	return &voteDao.Vote{Sender: "Test-sender"}, nil
}

func (s *VoteRepoMock) InsertVoteToStorage(ctx context.Context, in *voteDao.Vote, opts ...grpc.CallOption) (*voteDao.InsertResult, error) {
	return &voteDao.InsertResult{ObjectID: []byte("test-id")}, nil
}

func (s *VoteRepoMock) DeleteVoteFromStorage(ctx context.Context, in *voteDao.VoteFilter, opts ...grpc.CallOption) (*voteDao.Empty, error) {
	return &voteDao.Empty{}, nil
}

func (s *VoteRepoMock) GetRating(ctx context.Context, in *voteDao.Target, opts ...grpc.CallOption) (*voteDao.Rating, error) {
	return &voteDao.Rating{Target: "Test-target", Rating: 15}, nil
}

func GetVoteUsecase() *VoteUsecase {
	voteClient := &VoteRepoMock{}

	return NewVoteUsecase(voteClient)
}

func TestUserUsecase_RetractVote(t *testing.T) {
	uu := GetVoteUsecase()

	update := &model.Update{Nickname: "user2"}
	sender := "user1"

	t.Run("retract vote of user 1 from user 2", func(t *testing.T) {
		err := uu.RetractVote(update, sender)

		t.Log(err)

		assert.Nil(t, err)
	})
}

func TestUserUsecase_StoreVote(t *testing.T) {
	uu := GetVoteUsecase()

	now := time.Now().Unix()

	vote1 := &model.Vote{
		Sender:  "user1",
		Target:  "user2",
		Vote:    1,
		VotedAt: now,
	}

	t.Run("store vote", func(t *testing.T) {
		result, err := uu.StoreVote(vote1)

		t.Logf("%+v\n", result)

		assert.Nil(t, err)
		assert.IsType(t, &repository.VoteInsertResult{}, result)
	})
}
