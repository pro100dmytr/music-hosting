package playlist

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
	"net/http"
	"strconv"
)

type PlaylistHandler interface {
	CreatePlaylist() gin.HandlerFunc
	GetPlaylistByID() gin.HandlerFunc
	GetAllPlaylists() gin.HandlerFunc
	UpdatePlaylist() gin.HandlerFunc
	DeletePlaylist() gin.HandlerFunc
}

type Handler struct {
	store  *repository.PlaylistStorage
	logger *slog.Logger
}

func NewHandler(store *repository.PlaylistStorage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) CreatePlaylist() gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlist models.Playlist
		if err := c.ShouldBindJSON(&playlist); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if playlist.Name == "" {
			h.logger.Error("Playlist name is required", slog.Any("Error", "Name is required"))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
			return
		}

		err := h.store.Create(c.Request.Context(), &playlist)
		if err != nil {
			h.logger.Error("Error creating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"playlist": playlist})
	}
}

func (h *Handler) GetPlaylistByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid playlist ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
			return
		}

		playlist, err := h.store.Get(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error getting playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, playlist)
	}
}

func (h *Handler) GetAllPlaylists() gin.HandlerFunc {
	return func(c *gin.Context) {
		playlists, err := h.store.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error getting playlists", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, playlists)
	}
}

func (h *Handler) UpdatePlaylist() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid playlist ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
			return
		}

		var playlist models.Playlist
		if err := c.ShouldBindJSON(&playlist); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = h.store.Update(c.Request.Context(), id, &playlist)
		if err != nil {
			h.logger.Error("Error updating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"playlist": playlist})
	}
}

func (h *Handler) DeletePlaylist() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid playlist ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
			return
		}

		err = h.store.Delete(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error deleting playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted"})
	}
}
