package userDao

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
)

type UserRepoController interface {
	FindUserInStorage(username *Username) (*model.User, error)
	InsertUserToStorage(user *User) (*model.InsertResult, error)
	DeleteUserFromStorage(username *Username) error
	UpdateUsernameInStorage(userWithNewUsername *User, toReplaceUsername *Username) error
	UpdatePasswordInStorage(u *User) error
}

type Server struct {
	UserRepo UserRepoController
	UnimplementedUserRepoServer
}

func (s *Server) FindUserInStorage(_ context.Context, username *Username) (*User, error) {
	user, err := s.UserRepo.FindUserInStorage(username)
	if err != nil {
		return nil, err
	}

	convertedUser := MarshalTogRPCUserObj(user)

	return convertedUser, err
}

func (s *Server) InsertUserToStorage(_ context.Context, user *User) (*InsertResult, error) {
	insertResult, err := s.UserRepo.InsertUserToStorage(user)
	if err != nil {
		return nil, err
	}

	convertedInsertResult := MarshalTogRPCInsertResult(insertResult)

	return convertedInsertResult, err
}

func (s *Server) DeleteUserFromStorage(_ context.Context, username *Username) (*Empty, error) {
	err := s.UserRepo.DeleteUserFromStorage(username)
	return nil, err
}

func (s *Server) UpdateUsernameInStorage(_ context.Context, user *UserReplace) (*Empty, error) {
	err := s.UserRepo.UpdateUsernameInStorage(user.User, user.ToUpdateUserWithUsername)
	return nil, err
}

func (s *Server) UpdatePasswordInStorage(_ context.Context, user *User) (*Empty, error) {
	err := s.UserRepo.UpdatePasswordInStorage(user)
	return nil, err
}
