package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

const maxNumberVotesPerHistory = 10

func randomVoteInt() int {
	randomNumber := rand.Intn(2) // Generate a random number between 0 and 1
	var result int
	if randomNumber == 0 {
		result = -1
	} else {
		result = 1
	}

	return result
}

func GetUserUsecase() *UserUsecase {
	cfg, err := config.InitConfig(".usecaseConfig.json")
	if err != nil {
		panic(err)
	}

	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(db, &cfg.Database)

	return NewUserUsecase(r.Ur, r.Pr)
}

func TestUserUsecase_AddUpdateVoteSender(t *testing.T) {
	uu := GetUserUsecase()

	t.Run("User1 votes User2", func(t *testing.T) {
		now := time.Now().Unix()

		update := &model.Update{
			Nickname: "user1",
		}

		vote := &model.Vote{
			Sender:  "user1",
			Target:  "user2",
			Vote:    1,
			VotedAt: &now,
		}

		err := uu.AddUpdateVoteSender(update, vote)
		assert.Nil(t, err)

		err = uu.AddUpdateVoteSender(update, vote)
		assert.Error(t, err)
	})
}

func TestVoteHistory_removeVoteFromHistory(t *testing.T) {
	var histories []*VoteHistory

	for i := 1; i < maxNumberVotesPerHistory; i++ {
		var votes []*model.Vote

		now := time.Now().Unix()

		for j := 1; j <= i; j++ {
			votes = append(votes, &model.Vote{Vote: randomVoteInt(), VotedAt: &now})
		}

		index := 1
		histories = append(histories, &VoteHistory{history: votes, index: &index})
	}

	for i, vh := range histories {
		vh.removeVoteFromHistory()

		t.Logf("This history has %d votes", i+1)

		if i == 0 {
			assert.Nil(t, vh.history)
		} else {
			assert.Equal(t, len(vh.history), i)
		}
	}
}
