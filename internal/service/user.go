package service

import (
	"context"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
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

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	err := ValidateUser(user)
	if err != nil {
		return err
	}

	hashedPassword, salt, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	repoUser := repository.User{
		Login:    user.Login,
		Email:    user.Email,
		Password: hashedPassword,
		Salt:     salt,
	}

	id, err := s.userRepo.Create(ctx, &repoUser)
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	repoUser, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:    repoUser.ID,
		Login: repoUser.Login,
		Email: repoUser.Email,
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, user *models.User) error {
	err := ValidateUser(user)
	if err != nil {
		return err
	}

	hashedPassword, salt, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	repoUser := &repository.User{
		Login:    user.Login,
		Email:    user.Email,
		Password: hashedPassword,
		Salt:     salt,
	}
	err = s.userRepo.Update(ctx, repoUser, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	err = s.userRepo.RemovePlaylistsFromUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUsersWithPagination(ctx context.Context, limit, offset string) ([]*models.User, error) {
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

	var users []*models.User
	for _, repoUser := range repoUsers {
		user := &models.User{
			ID:    repoUser.ID,
			Login: repoUser.Login,
			Email: repoUser.Email,
		}
		users = append(users, user)
	}

	return users, nil
}
