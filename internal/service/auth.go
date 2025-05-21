package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"music-hosting/internal/auth"
	"music-hosting/internal/models"
)

func ValidateUser(user *models.User) error {
	if user.Login == "" {
		return errors.New("login is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}

func HashPassword(password string) (string, string, error) {
	salt, err := GenerateSalt()
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

func ValidateTrack(track *models.Track) error {
	if track.Name == "" {
		return errors.New("name is required")
	}
	if track.Artist == "" {
		return errors.New("artist is required")
	}
	if track.URL == "" {
		return errors.New("url is required")
	}
	return nil
}
