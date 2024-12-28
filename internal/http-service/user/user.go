package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
	"music-hosting/pkg/utils/userutils"
	"net/http"
	"strconv"
)

type UserHandler interface {
	CreateUser() gin.HandlerFunc
	GetUserID() gin.HandlerFunc
	GetAllUsers() gin.HandlerFunc
	UpdateUser() gin.HandlerFunc
	DeleteUser() gin.HandlerFunc
	SaveTrack() gin.HandlerFunc
}

type Handler struct {
	store  *repository.UserStorage
	logger *slog.Logger
}

func NewHandler(store *repository.UserStorage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := &models.User{}

		if err := c.ShouldBindJSON(user); err != nil {
			h.logger.Error("Invalid request", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"err": err})
			return
		}

		if err := userutils.ValidateUser(user); err != nil {
			h.logger.Error("Failed to validate user", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		hashedPassword, err := userutils.HashPassword(user.Password)
		if err != nil {
			h.logger.Error("Failed to hash password", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		user.Password = hashedPassword

		err = h.store.Create(c.Request.Context(), user)
		if err != nil {
			h.logger.Error("Failed to create user", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": user})
	}
}

func (h *Handler) GetUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"err": err})
			return
		}

		user, err := h.store.Get(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"err": err})
				return
			}

			h.logger.Error("Failed to get user", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func (h *Handler) GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.store.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error getting users", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting users"})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func (h *Handler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"err": err})
			return
		}

		user := &models.User{}
		if err := c.ShouldBindJSON(user); err != nil {
			h.logger.Error("Invalid request", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"err": err})
			return
		}

		err = userutils.ValidateUser(user)
		if err != nil {
			h.logger.Error("Failed to validate user", slog.Any("err", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		hashedPassword, err := userutils.HashPassword(user.Password)
		if err != nil {
			h.logger.Error("Failed to hash password", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		user.Password = hashedPassword

		err = h.store.Update(c.Request.Context(), user, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("Error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			h.logger.Error("Error update user", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})

	}
}

func (h *Handler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"err": err})
			return
		}

		err = h.store.Delete(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("User not found", slog.Any("Error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			h.logger.Error("Error delete user", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": nil})
	}
}

//func (h *Handler) SaveTrack() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		idStr := c.Param("id")
//		id, err := strconv.Atoi(idStr)
//		if err != nil {
//			h.logger.Error("Invalid user id parameter", slog.Any("error", err))
//			c.JSON(http.StatusBadRequest, gin.H{"err": err})
//			return
//		}
//
//		var req struct {
//			PlaylistID int    `json:"playlist_id"`
//			Name       string `json:"name"`
//			Artist     string `json:"artist"`
//			URL        string `json:"url"`
//		}
//
//		if err := c.ShouldBindJSON(&req); err != nil {
//			h.logger.Error("Invalid request", slog.Any("err", err))
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//			return
//		}
//
//		track := &models.Track{
//			Name:       req.Name,
//			Artist:     req.Artist,
//			URL:        req.URL,
//			PlaylistID: req.PlaylistID,
//		}
//
//		if err := h.store.Save(c.Request.Context(), id, track); err != nil {
//			h.logger.Error("Failed to save track", slog.Any("err", err))
//			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
//			return
//		}
//
//		c.JSON(http.StatusOK, gin.H{"message": "Track saved successfully"})
//	}
//}
