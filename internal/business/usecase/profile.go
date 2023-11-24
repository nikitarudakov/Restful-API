package usecase

import (
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/repository"
)

type ProfileUseCase struct {
	profileRepo repository.ProfileRepoController
}

func NewProfileUseCase(profileRepo repository.ProfileRepoController) *ProfileUseCase {
	return &ProfileUseCase{profileRepo: profileRepo}
}

// CreateProfile creates profile for model.User and stores it in DB with repository.ProfileRepository
func (pu *ProfileUseCase) CreateProfile(p *model.Profile) (*repository.InsertResult, error) {
	result, err := pu.profileRepo.InsertProfileToStorage(p)
	return result, err
}

// UpdateProfile updates profile of model.User in DB with repository.ProfileRepository
func (pu *ProfileUseCase) UpdateProfile(p *model.Update, profileName string) error {
	return pu.profileRepo.UpdateProfileInStorage(p, profileName)
}

func (pu *ProfileUseCase) DeleteProfile(profileName string) error {
	return pu.profileRepo.DeleteProfileFromStorage(profileName)
}

func (pu *ProfileUseCase) ListProfiles(page int64) ([]*model.Profile, error) {
	profiles, err := pu.profileRepo.ListProfilesFromStorage(page)
	return profiles, err
}
