package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models"
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
		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		newUser, err := h.service.CreateUser(c.Request.Context(), &user)
		if err != nil {
			h.logger.Error("Failed to create user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": newUser})
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
				h.logger.Warn("User not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Failed to retrieve user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func (h *UserHandler) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.service.GetAllUsers(c.Request.Context())
		if err != nil {
			h.logger.Error("Failed to retrieve users", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
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

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			h.logger.Error("Invalid request body", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		updatedUser, err := h.service.UpdateUser(c.Request.Context(), id, &user)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Warn("User not found", slog.Any("id", id))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Failed to update user", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": updatedUser})
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
				h.logger.Warn("User not found", slog.Any("id", id))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Failed to delete user", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	}
}

func (h *UserHandler) GetUserWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		min, err := strconv.Atoi(c.DefaultQuery("min", "1"))
		if err != nil || min < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min parameter"})
			return
		}

		max, err := strconv.Atoi(c.DefaultQuery("max", "10"))
		if err != nil || max < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max parameter"})
			return
		}

		offset := (min - 1) * max

		users, err := h.service.GetUsersWithPagination(c.Request.Context(), max, offset)
		if err != nil {
			h.logger.Error("Failed to retrieve users", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
