package usecase

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/grpc/profileDao"
)

type ProfileUseCase struct {
	grpcClient profileDao.ProfileRepoClient
}

func NewProfileUseCase(client profileDao.ProfileRepoClient) *ProfileUseCase {
	return &ProfileUseCase{grpcClient: client}
}

// CreateProfile creates profile for model.User and stores it in DB with repository.ProfileRepository
func (pu *ProfileUseCase) CreateProfile(p *model.Profile) (*profileDao.InsertResult, error) {
	profileObj := profileDao.MarshalTogRPCProfileObj(p)

	insertResult, err := pu.grpcClient.InsertProfileToStorage(context.Background(), profileObj)

	return insertResult, err
}

// UpdateProfile updates profile of model.User in DB with repository.ProfileRepository
func (pu *ProfileUseCase) UpdateProfile(p *model.Update, profileName string) error {
	profileUpdate := profileDao.MarshalTogRPCProfileUpdate(p)
	profileUpdate.ProfileName = &profileDao.ProfileName{Name: profileName}

	_, err := pu.grpcClient.UpdateProfileInStorage(context.Background(), profileUpdate)

	return err
}

func (pu *ProfileUseCase) DeleteProfile(profileName string) error {
	_, err := pu.grpcClient.DeleteProfileFromStorage(context.Background(),
		&profileDao.ProfileName{Name: profileName})

	return err
}

func (pu *ProfileUseCase) ListProfiles(page int64) ([]*model.Profile, error) {
	var profiles []*model.Profile

	gRPCProfilesSlice, err := pu.grpcClient.ListProfilesFromStorage(context.Background(),
		&profileDao.Page{Num: page})

	for _, profile := range gRPCProfilesSlice.Profiles {
		profiles = append(profiles, profileDao.UnmarshalTogRPCProfileObj(profile))
	}

	return profiles, err
}
