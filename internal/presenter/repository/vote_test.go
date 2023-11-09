package repository

import (
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func GetVoteRepository() *VoteRepository {
	cfg, err := config.InitConfig(".repoConfig.json")
	if err != nil {
		panic(err)
	}

	db, err := datastore.NewDB(&cfg.Database)
	if err != nil {
		panic(err)
	}

	cacheDB, err := cache.NewCacheDatabase(&cfg.Cache)

	return NewVoteRepository(db.Collection(cfg.Database.VoteRepo), cacheDB)
}

func TestVoteRepository_CalcTotalRating(t *testing.T) {
	target := "user2"

	vr := GetVoteRepository()

	t.Run("calculate agg rating for user2", func(t *testing.T) {
		rating, err := vr.GetRating(target)

		t.Logf("%+v\n", rating)
		t.Logf("%d", rating.Rating)
		t.Log(err)

		assert.Nil(t, err)
		assert.IsType(t, &model.Rating{}, rating)
	})

	t.Run("get rating for user2 from CACHE", func(t *testing.T) {
		rating, err := vr.GetRating(target)

		t.Logf("%+v\n", rating)
		t.Logf("%d", rating.Rating)
		t.Log(err)

		assert.Nil(t, err)
		assert.IsType(t, &model.Rating{}, rating)
	})
}
