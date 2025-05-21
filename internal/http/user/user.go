package user

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"music-hosting/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUsersWithPagination(ctx context.Context, limit, offset string) ([]*models.User, error)
	UpdateUser(ctx context.Context, id int, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
	GetToken(ctx context.Context, login string, password string) (string, error)
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.UserRequest

		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userServ := models.User{
			Login:    user.Login,
			Email:    user.Email,
			Password: user.Password,
		}

		err := h.service.CreateUser(c.Request.Context(), &userServ)
		if err != nil {
			h.logger.Error("Failed to create user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, nil)
	}
}

func (h *Handler) GetUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := h.service.GetUser(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			h.logger.Error("Failed to retrieve user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		userResponse := models.UserResponse{
			ID:    user.ID,
			Login: user.Login,
			Email: user.Email,
		}

		c.JSON(http.StatusOK, userResponse)
	}
}

func (h *Handler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user models.UserRequest
		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request body", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userServ := models.User{
			Login:    user.Login,
			Email:    user.Email,
			Password: user.Password,
		}

		err = h.service.UpdateUser(c.Request.Context(), id, &userServ)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Failed to update user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		err = h.service.DeleteUser(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Failed to delete user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) GetUserWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		users, err := h.service.GetUsersWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Failed to retrieve users with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		var usersResponse []models.UserResponse
		for _, user := range users {
			userResponse := models.UserResponse{
				ID:    user.ID,
				Login: user.Login,
				Email: user.Email,
			}
			usersResponse = append(usersResponse, userResponse)
		}

		c.JSON(http.StatusOK, usersResponse)
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.UserRequest

		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		token, err := h.service.GetToken(c.Request.Context(), user.Login, user.Password)
		if err != nil {
			h.logger.Error("User not found", slog.Any("error", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		tokenResponse := models.TokenResponse{
			Token: token,
		}

		c.JSON(http.StatusOK, tokenResponse)
	}
}
