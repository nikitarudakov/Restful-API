package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/auth"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/repository"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/controller"
	"git.foxminded.ua/foxstudent106092/user-management/tools/hashing"
	customValidator "git.foxminded.ua/foxstudent106092/user-management/tools/validator"
	libValidator "github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var validator = &customValidator.CustomValidator{
	Validator: libValidator.New(libValidator.WithRequiredStructEnabled()),
}

type GRPCAuthServer struct {
	UserUseCase    controller.UserManager
	ProfileUseCase controller.ProfileManager
	Cfg            *config.Config
	UnimplementedAuthControllerServer
}

func (as *GRPCAuthServer) UpdatePassword(ctx context.Context, credentials *Credentials) (*Empty, error) {
	var u model.User

	u.Username = credentials.Username

	password, err := hashing.HashPassword(credentials.Password)
	u.Password = password
	if err != nil {
		return &Empty{}, fmt.Errorf("error hashing password: %w", err)
	}

	if err = as.UserUseCase.UpdatePassword(&u); err != nil {
		return &Empty{}, fmt.Errorf("error updating password: %w", err)
	}

	return &Empty{}, nil
}

func (as *GRPCAuthServer) Login(ctx context.Context, credentials *Credentials) (*Token, error) {
	var u model.User

	u.Username = credentials.Username
	password := credentials.Password

	userFromDB, err := as.UserUseCase.FindUser(&u)
	if err != nil {
		return nil, fmt.Errorf("user was not found: %w", err)
	}

	u.Role = userFromDB.Role

	if err = hashing.CheckPassword(userFromDB.Password, password); err != nil ||
		subtle.ConstantTimeCompare([]byte(u.Username), []byte(userFromDB.Username)) != 1 {

		return nil, fmt.Errorf("username/password is incorrect: %w", err)
	}

	authCfg := as.Cfg.Auth
	token, err := auth.GenerateJWTToken(&u, []byte(authCfg.SecretKey))
	if err != nil {
		return nil, err
	}

	return &Token{Token: token}, nil
}

func (as *GRPCAuthServer) Register(ctx context.Context, user *RegisterUser) (*Empty, error) {
	_, err := as.registerUser(user.Credentials, user.Role)
	if err != nil {
		return nil, err
	}

	_, err = as.registerProfile(user.Profile, user.Credentials.Username)
	if err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

func (as *GRPCAuthServer) registerUser(credentials *Credentials, role string) (*repository.InsertResult, error) {
	var u model.User

	u.Username = credentials.Username
	u.Password = credentials.Password
	u.Role = role

	log.Trace().Str("username", u.Username).
		Str("role", u.Role).
		Msg("register request")

	if err := validator.Validate(u); err != nil {
		return nil, err
	}

	if u.Role == "vote" {
		if subtle.ConstantTimeCompare([]byte(u.Username), []byte(as.Cfg.Admin.Username)) != 1 ||
			subtle.ConstantTimeCompare([]byte(u.Password), []byte(as.Cfg.Admin.Password)) != 1 {

			return nil, errors.New("error: unable to register ADMIN user")
		}
	}

	hashedPassword, err := hashing.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}
	u.Password = hashedPassword

	result, err := as.UserUseCase.CreateUser(&u)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (as *GRPCAuthServer) registerProfile(profile *Profile, username string) (*repository.InsertResult, error) {
	domainProfile := &model.Profile{
		Nickname:  profile.Nickname,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}

	// check if new Nickname was passed
	if domainProfile.Nickname == "" {
		domainProfile.Nickname = username // assign User username to Update nickname
	}

	if err := validator.Validate(profile); err != nil {
		return nil, err
	}

	result, err := as.ProfileUseCase.CreateProfile(domainProfile)
	if err != nil {
		return nil, err
	}

	return result, nil
}
