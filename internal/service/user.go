package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"music-hosting/internal/auth"
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

func ValidateUser(user *models.User) error {
	if user.Login == "" {
		return errors.New("login is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func HashPassword(password string) (string, string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", "", err
	}

	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	hashedPassword := hash.Sum(nil)

	return hex.EncodeToString(hashedPassword), hex.EncodeToString(salt), nil
}

func CheckPassword(password, storedHash, storedSalt string) (bool, error) {
	storedHashBytes, err := hex.DecodeString(storedHash)
	if err != nil {
		return false, fmt.Errorf("failed to decode stored hash: %w", err)
	}

	storedSaltBytes, err := hex.DecodeString(storedSalt)
	if err != nil {
		return false, fmt.Errorf("failed to decode stored salt: %w", err)
	}

	hash := sha256.New()
	hash.Write(storedSaltBytes)
	hash.Write([]byte(password))
	computedHash := hash.Sum(nil)

	return bytes.Equal(computedHash, storedHashBytes), nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	user := &models.User{
		ID:    repoUser.ID,
		Login: repoUser.Login,
		Email: repoUser.Email,
	}

	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	repoUsers, err := s.userRepo.GetAll(ctx)
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
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	err := s.userRepo.Delete(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.ErrNoRows
		}
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

func (s *UserService) GetToken(ctx context.Context, login string, password string) (string, error) {
	if login == "" || password == "" {
		return "", fmt.Errorf("invalid login or password")
	}

	user, err := s.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	isValidPassword, err := CheckPassword(password, user.Password, user.Salt)
	if err != nil || !isValidPassword {
		return "", fmt.Errorf("invalid password")
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
