package playlist

import (
	"context"
	"log/slog"
	"music-hosting/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreatePlaylist(ctx context.Context, playlist *models.Playlist) error
	GetPlaylistByID(ctx context.Context, id int) (*models.Playlist, error)
	GetPlaylists(ctx context.Context, name string, userID int) ([]*models.Playlist, error)
	UpdatePlaylist(ctx context.Context, playlist *models.Playlist, trackIDs []int) error
	DeletePlaylist(ctx context.Context, id int) error
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) CreatePlaylist() gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlist models.CreatePlaylistRequest
		if err := c.ShouldBindJSON(&playlist); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			h.logger.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
			return
		}

		playlistServ := models.Playlist{
			Name:   playlist.Name,
			UserID: userID.(int),
		}

		err := h.service.CreatePlaylist(c.Request.Context(), &playlistServ)
		if err != nil {
			h.logger.Error("Error creating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating playlist"})
			return
		}

		c.JSON(http.StatusCreated, nil)
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

		playlist, err := h.service.GetPlaylistByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error getting playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlist"})
			return
		}

		if playlist == nil {
			h.logger.Info("Playlist not found", slog.Int("playlistID", id))
			c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			h.logger.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
			return
		}

		playlistResponse := models.PlaylistResponse{
			ID:        playlist.ID,
			Name:      playlist.Name,
			UserID:    userID.(int),
			Tracks:    playlist.Tracks,
			CreatedAt: playlist.CreatedAt,
			UpdatedAt: playlist.UpdatedAt,
		}

		c.JSON(http.StatusOK, playlistResponse)
	}
}

func (h *Handler) GetPlaylists() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		userIDQuery := c.Query("userid")

		var userID int
		var err error

		if userIDQuery != "" {
			userID, err = strconv.Atoi(userIDQuery)
			if err != nil {
				h.logger.Error("Invalid user ID", slog.Any("Error", err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}
		}

		playlists, err := h.service.GetPlaylists(c.Request.Context(), name, userID)
		if err != nil {
			h.logger.Error("Error getting playlists", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlists"})
			return
		}

		var playlistsResponse []models.PlaylistResponse
		for _, playlist := range playlists {
			playlistResponse := models.PlaylistResponse{
				ID:        playlist.ID,
				Name:      playlist.Name,
				UserID:    playlist.UserID,
				Tracks:    playlist.Tracks,
				CreatedAt: playlist.CreatedAt,
				UpdatedAt: playlist.UpdatedAt,
			}
			playlistsResponse = append(playlistsResponse, playlistResponse)
		}

		c.JSON(http.StatusOK, playlistsResponse)
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

		var playlistRequest models.CreatePlaylistRequest
		if err := c.ShouldBindJSON(&playlistRequest); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			h.logger.Error("User ID not found in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
			return
		}

		playlist := models.Playlist{
			ID:     id,
			Name:   playlistRequest.Name,
			UserID: userID.(int),
		}

		err = h.service.UpdatePlaylist(c.Request.Context(), &playlist, playlistRequest.TrackIDs)
		if err != nil {
			h.logger.Error("Error updating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update playlist"})
			return
		}

		c.JSON(http.StatusOK, nil)
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

		err = h.service.DeletePlaylist(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error deleting playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting playlist"})
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}
