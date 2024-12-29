package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
	"music-hosting/pkg/utils/userutils"
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

func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("User not found", slog.Any("error", err))
			return nil, sql.ErrNoRows
		}

		s.logger.Error("User not found", slog.Any("error", err))
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		s.logger.Error("failed to get all users", slog.Any("error", err))
		return nil, err
	}
	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := userutils.ValidateUser(user)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := userutils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("failed to hash password", slog.Any("error", err))
		return nil, err
	}
	user.Password = hashedPassword

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", slog.Any("error", err))
		return nil, err
	}
	user.ID = id
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, user *models.User) (*models.User, error) {
	err := userutils.ValidateUser(user)
	if err != nil {
		s.logger.Error("failed to validate user", slog.Any("error", err))
		return nil, err
	}

	hashedPassword, err := userutils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error("failed to hash password", slog.Any("error", err))
		return nil, err
	}
	user.Password = hashedPassword

	err = s.userRepo.Update(ctx, user, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("failed to update user", slog.Any("error", err))
			return nil, sql.ErrNoRows
		}

		s.logger.Error("failed to update user", slog.Any("error", err))
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Error("failed to delete user", slog.Any("error", err))
			return sql.ErrNoRows
		}

		s.logger.Error("failed to delete user", slog.Any("error", err))
		return err
	}

	return nil
}

func (s *UserService) GetUsersWithPagination(ctx context.Context, limit, offset int) ([]*models.User, error) {
	users, err := s.userRepo.GetUsers(ctx, limit, offset)
	if err != nil {
		s.logger.Error("failed to get users", slog.Any("error", err))
		return nil, err
	}

	return users, nil
}
