package profile

import "context"

type ProfileDAOServer struct {
	UnsafeProfileRepoServer
}

func (ps *ProfileDAOServer) FindProfileInStorage(context.Context, *ProfileName) (
	*ProfileObj, error) {

	return nil, nil
}

func (ps *ProfileDAOServer) InsertProfileToStorage(context.Context,
	*ProfileObj) (*ProfileObj, error) {

	return nil, nil
}

func (ps *ProfileDAOServer) DeleteProfileFromStorage(context.Context,
	*ProfileName) (*ProfileObj, error) {

	return nil, nil
}

func (ps *ProfileDAOServer) UpdateProfileInStorage(context.Context,
	*ProfileUpdate) (*ProfileObj, error) {

	return nil, nil
}

func (ps *ProfileDAOServer) ListProfilesFromStorage(context.Context,
	*ProfileName) (*ProfileObj, error) {

	return nil, nil
}
