package profileDao

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
)

type ProfileRepoController interface {
	FindProfileInStorage(profileName *ProfileName) (*model.Profile, error)
	InsertProfileToStorage(profileObj *ProfileObj) (*InsertResult, error)
	DeleteProfileFromStorage(profileName *ProfileName) error
	UpdateProfileInStorage(profileUpdate *ProfileUpdate) error
	ListProfilesFromStorage(page *Page) ([]model.Profile, error)
}

type Server struct {
	ProfileRepo ProfileRepoController
	UnimplementedProfileRepoServer
}

func (s *Server) FindProfileInStorage(_ context.Context,
	profileName *ProfileName) (*ProfileObj, error) {

	p, err := s.ProfileRepo.FindProfileInStorage(profileName)
	if err != nil {
		return nil, err
	}

	profileObj := MarshalTogRPCProfileObj(p)

	return profileObj, nil
}

func (s *Server) InsertProfileToStorage(_ context.Context,
	profileObj *ProfileObj) (*InsertResult, error) {

	insertResult, err := s.ProfileRepo.InsertProfileToStorage(profileObj)
	if err != nil {
		return nil, err
	}

	return insertResult, nil
}

func (s *Server) DeleteProfileFromStorage(_ context.Context,
	profileName *ProfileName) (*Empty, error) {

	err := s.ProfileRepo.DeleteProfileFromStorage(profileName)
	return nil, err
}

func (s *Server) UpdateProfileInStorage(_ context.Context,
	profileUpdate *ProfileUpdate) (*Empty, error) {

	err := s.ProfileRepo.UpdateProfileInStorage(profileUpdate)
	return nil, err
}

func (s *Server) ListProfilesFromStorage(_ context.Context,
	page *Page) (*SliceOfProfileObj, error) {
	var gRPCProfilesSlice []*ProfileObj

	profiles, err := s.ProfileRepo.ListProfilesFromStorage(page)
	if err != nil {
		return nil, err
	}

	for _, profile := range profiles {
		gRPCProfile := MarshalTogRPCProfileObj(&profile)

		gRPCProfilesSlice = append(gRPCProfilesSlice, gRPCProfile)
	}

	return &SliceOfProfileObj{Profiles: gRPCProfilesSlice}, err
}
