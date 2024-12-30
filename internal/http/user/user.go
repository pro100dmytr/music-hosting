package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models/https"
	"music-hosting/internal/models/services"
	"music-hosting/internal/service"
	"music-hosting/pkg/utils/jwtutils"
	"music-hosting/pkg/utils/userutils"
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

		c.JSON(http.StatusCreated, gin.H{
			"message": "User created",
			"userID":  userServ.ID,
		})
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

		userHttp := https.User{
			ID:         user.ID,
			Login:      user.Login,
			Email:      user.Email,
			Password:   user.Password,
			PlaylistID: user.PlaylistID,
		}

		c.JSON(http.StatusOK, gin.H{"user": userHttp})
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

		var usersHttp []https.User
		for _, user := range users {
			userHttp := https.User{
				ID:         user.ID,
				Login:      user.Login,
				Email:      user.Email,
				Password:   user.Password,
				PlaylistID: user.PlaylistID,
			}
			usersHttp = append(usersHttp, userHttp)
		}

		c.JSON(http.StatusOK, gin.H{"users": usersHttp})
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

		userHttp := https.User{
			ID:         updatedUser.ID,
			Login:      updatedUser.Login,
			Email:      updatedUser.Email,
			Password:   updatedUser.Password,
			PlaylistID: updatedUser.PlaylistID,
		}

		c.JSON(http.StatusOK, gin.H{"user": userHttp})
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

		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
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

		var usersHttp []https.User
		for _, user := range users {
			userHttp := https.User{
				ID:         user.ID,
				Login:      user.Login,
				Email:      user.Email,
				Password:   user.Password,
				PlaylistID: user.PlaylistID,
			}
			usersHttp = append(usersHttp, userHttp)
		}

		c.JSON(http.StatusOK, gin.H{"users": usersHttp})
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

		existingUser, err := h.service.GetUserByLogin(c.Request.Context(), user.Login)
		if err != nil {
			h.logger.Error("User not found", slog.Any("error", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if !userutils.CheckPasswordHash(user.Password, existingUser.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := jwtutils.GenerateToken(existingUser.ID)
		if err != nil {
			h.logger.Error("Failed to generate token", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
