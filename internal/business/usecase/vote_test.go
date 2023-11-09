package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/registry"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/repository"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func GetVoteUsecase() *VoteUsecase {
	cfg, err := config.InitConfig(".usecaseConfig.json")
	if err != nil {
		panic(err)
	}

	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	cacheDB, err := cache.NewCacheDatabase(&cfg.Cache)
	if err != nil {
		panic(err)
	}

	r := registry.NewRegistry(db, &cfg.Database, cacheDB)

	return NewVoteUsecase(r.Pr, r.Vr)
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

	t.Run("try retract vote of user 1 from user 2 twice", func(t *testing.T) {
		err := uu.RetractVote(update, sender)

		t.Log(err)

		assert.Error(t, err)
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

	vote2 := &model.Vote{
		Sender:  "user1",
		Target:  "user2",
		Vote:    1,
		VotedAt: now,
	}

	t.Run("try store vote twice", func(t *testing.T) {
		result, err := uu.StoreVote(vote2)

		t.Log(err)

		assert.ErrorContains(t, err, "twice is not allowed")
		assert.Nil(t, result)
	})

	vote3 := &model.Vote{
		Sender:  "user1",
		Target:  "user2",
		Vote:    -1,
		VotedAt: now,
	}

	t.Run("try store vote again within 1 hour timespan", func(t *testing.T) {
		result, err := uu.StoreVote(vote3)

		t.Log(err)

		assert.ErrorContains(t, err, "only once per hour")
		assert.Nil(t, result)
	})
}
