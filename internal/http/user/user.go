package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models/dto"
	"music-hosting/internal/models/https"
	"music-hosting/internal/models/services"
	"music-hosting/internal/service"
	"net/http"
	"strconv"
)

type Handler interface {
	CreateUser() gin.HandlerFunc
	GetUserID() gin.HandlerFunc
	GetAllUsers() gin.HandlerFunc
	GetUserWithPagination() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
}

type UserHandler struct {
	service *service.UserService
	logger  *slog.Logger
}

func NewHandler(service *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

func (h *UserHandler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user https.User

		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userServ := services.User{
			Login:      user.Login,
			Email:      user.Email,
			Password:   user.Password,
			PlaylistID: user.PlaylistID,
		}

		err := h.service.CreateUser(c.Request.Context(), &userServ)
		if err != nil {
			h.logger.Error("Failed to create user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		userResponse := dto.UserResponse{
			ID: userServ.ID,
		}

		c.JSON(http.StatusCreated, userResponse.ID)
	}
}

func (h *UserHandler) GetUserID() gin.HandlerFunc {
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

		userResponse := dto.UserResponse{
			ID:          user.ID,
			Login:       user.Login,
			Email:       user.Email,
			PlaylistsID: user.PlaylistID,
		}

		c.JSON(http.StatusOK, userResponse)
	}
}

func (h *UserHandler) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.service.GetAllUsers(c.Request.Context())
		if err != nil {
			h.logger.Error("Failed to retrieve users", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		var usersResponse []dto.UserResponse
		for _, user := range users {
			userResponse := dto.UserResponse{
				ID:          user.ID,
				Login:       user.Login,
				Email:       user.Email,
				PlaylistsID: user.PlaylistID,
			}
			usersResponse = append(usersResponse, userResponse)
		}

		c.JSON(http.StatusOK, usersResponse)
	}
}

func (h *UserHandler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var user https.User
		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request body", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userServ := services.User{
			Login:      user.Login,
			Email:      user.Email,
			Password:   user.Password,
			PlaylistID: user.PlaylistID,
		}

		updatedUser, err := h.service.UpdateUser(c.Request.Context(), id, &userServ)
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

		userResponse := dto.UserResponse{
			ID:          updatedUser.ID,
			Login:       updatedUser.Login,
			Email:       updatedUser.Email,
			PlaylistsID: updatedUser.PlaylistID,
		}

		c.JSON(http.StatusOK, userResponse)
	}
}

func (h *UserHandler) DeleteUser() gin.HandlerFunc {
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

		message := dto.MessageResponse{
			Message: "User deleted",
		}

		c.JSON(http.StatusOK, message)
	}
}

func (h *UserHandler) GetUserWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		users, err := h.service.GetUsersWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Failed to retrieve users with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		var usersResponse []dto.UserResponse
		for _, user := range users {
			userResponse := dto.UserResponse{
				ID:          user.ID,
				Login:       user.Login,
				Email:       user.Email,
				PlaylistsID: user.PlaylistID,
			}
			usersResponse = append(usersResponse, userResponse)
		}

		c.JSON(http.StatusOK, usersResponse)
	}
}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user https.User

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

		tokenResponse := dto.TokenResponse{
			Token: token,
		}

		c.JSON(http.StatusOK, tokenResponse)
	}
}
