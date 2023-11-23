package user

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/controller"
	"github.com/rs/zerolog/log"
)

type GRPCUserService struct {
	UserUseCase    controller.UserManager
	ProfileUseCase controller.ProfileManager
	CacheDb        *cache.Database
	UnimplementedUserControllerServer
}

func (us *GRPCUserService) DeleteUserAndProfile(ctx context.Context, username *Username) (*Empty, error) {
	err := us.ProfileUseCase.DeleteProfile(username.Username)
	return &Empty{}, err
}

func (us *GRPCUserService) ListProfiles(ctx context.Context, page *Page) (*Profiles, error) {
	profiles, err := us.ProfileUseCase.ListProfiles(page.Page)

	var convertedProfiles []*Profile

	for _, profile := range profiles {
		convertedProfile := &Profile{
			Nickname:  profile.Nickname,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
			DeletedAt: profile.DeletedAt,
		}

		convertedProfiles = append(convertedProfiles, convertedProfile)
	}

	return &Profiles{Profiles: convertedProfiles}, err
}

func (us *GRPCUserService) UpdateUserAndProfile(ctx context.Context, updateArg *UpdateArg) (*Empty, error) {
	update := &model.Update{
		Nickname:  updateArg.Update.Nickname,
		FirstName: updateArg.Update.FirstName,
		LastName:  updateArg.Update.LastName,
	}

	nickname := update.Nickname
	username := updateArg.Username.Username

	err := us.ProfileUseCase.UpdateProfile(update, username)
	if err != nil {
		return nil, err
	}

	if nickname != "" && nickname != username {
		if err = us.UserUseCase.UpdateUsername(nickname, username); err != nil {
			log.Error().Err(err).
				Str("old_username", username).
				Str("new_username", nickname).
				Msg("error updating username")
		}
	}

	return &Empty{}, nil
}
