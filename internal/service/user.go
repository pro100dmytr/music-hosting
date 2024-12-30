package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"music-hosting/internal/models/repositorys"
	"music-hosting/internal/models/services"
	"music-hosting/internal/repository"
	"music-hosting/pkg/utils/convertutils"
	"music-hosting/pkg/utils/userutils"
	"strconv"
)

type UserService struct {
	userRepo *repository.UserStorage
	logger   *slog.Logger
}

func NewUserService(userRepo *repository.UserStorage, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *services.User) error {
	err := userutils.ValidateUser(user)
	if err != nil {
		return err
	}

	hashedPassword, err := userutils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	repoUser := repositorys.User{
		Login:      user.Login,
		Email:      user.Email,
		Password:   hashedPassword,
		PlaylistID: convertutils.IntSliceConvertIntoString(user.PlaylistID),
	}

	id, err := s.userRepo.Create(ctx, &repoUser)
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (*services.User, error) {
	repoUser, err := s.userRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	playlistID := []int{}
	if repoUser.PlaylistID != "" {
		playlistID, err = convertutils.StringConvertIntoIntSlice(repoUser.PlaylistID)
		if err != nil {
			return nil, err
		}
	}

	user := &services.User{
		ID:         repoUser.ID,
		Login:      repoUser.Login,
		Email:      repoUser.Email,
		Password:   repoUser.Password,
		PlaylistID: playlistID,
	}

	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*services.User, error) {
	repoUsers, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var users []*services.User
	for _, repoUser := range repoUsers {
		playlistID := []int{}
		if repoUser.PlaylistID != "" {
			playlistID, err = convertutils.StringConvertIntoIntSlice(repoUser.PlaylistID)
			if err != nil {
				return nil, err
			}
		}

		user := &services.User{
			ID:         repoUser.ID,
			Login:      repoUser.Login,
			Email:      repoUser.Email,
			Password:   repoUser.Password,
			PlaylistID: playlistID,
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, user *services.User) (*services.User, error) {
	err := userutils.ValidateUser(user)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := userutils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	playlistIDString := convertutils.IntSliceConvertIntoString(user.PlaylistID)
	repoUser := &repositorys.User{
		Login:      user.Login,
		Email:      user.Email,
		Password:   hashedPassword,
		PlaylistID: playlistIDString,
	}

	err = s.userRepo.Update(ctx, repoUser, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	user.ID = id
	user.Password = hashedPassword
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	return nil
}

func (s *UserService) GetUsersWithPagination(ctx context.Context, limit, offset string) ([]*services.User, error) {
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		return nil, err
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil || offsetInt < 0 {
		return nil, err
	}

	repoUsers, err := s.userRepo.GetUsers(ctx, offsetInt, limitInt)
	if err != nil {
		return nil, err
	}

	var users []*services.User
	for _, repoUser := range repoUsers {
		playlistID, err := convertutils.StringConvertIntoIntSlice(repoUser.PlaylistID)
		if err != nil {
			return nil, err
		}

		user := &services.User{
			ID:         repoUser.ID,
			Login:      repoUser.Login,
			Email:      repoUser.Email,
			Password:   repoUser.Password,
			PlaylistID: playlistID,
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) GetUserByLogin(ctx context.Context, login string) (*services.User, error) {
	if login == "" {
		return nil, fmt.Errorf("invalid login")
	}

	user, err := s.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	users := &services.User{
		Login:    user.Login,
		Password: user.Password,
	}

	return users, nil
}
