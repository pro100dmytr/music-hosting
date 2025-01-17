package playlist

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
	CreatePlaylist(ctx context.Context, playlist *models.Playlist) error
	GetPlaylistByID(ctx context.Context, id int) (*models.Playlist, error)
	GetAllPlaylists(ctx context.Context) ([]*models.Playlist, error)
	GetPlaylistsByName(ctx context.Context, name string) ([]*models.Playlist, error)
	GetPlaylistsByUserID(ctx context.Context, userID int) ([]*models.Playlist, error)
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
		var playlist models.PlaylistRequest
		if err := c.ShouldBindJSON(&playlist); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
			return
		}

		playlistServ := models.Playlist{
			Name:   playlist.Name,
			UserID: playlist.UserID,
		}

		err := h.service.CreatePlaylist(c.Request.Context(), &playlistServ)
		if err != nil {
			h.logger.Error("Error creating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating playlist"})
			return
		}

		message := models.MessageResponse{
			Message: "Created playlist",
		}

		c.JSON(http.StatusCreated, message)
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
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Playlist not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
				return
			}

			h.logger.Error("Error getting playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlist"})
			return
		}

		playlistResponse := models.PlaylistResponse{
			ID:        playlist.ID,
			Name:      playlist.Name,
			UserID:    playlist.UserID,
			Tracks:    playlist.Tracks,
			CreatedAt: playlist.CreatedAt,
			UpdatedAt: playlist.UpdatedAt,
		}

		c.JSON(http.StatusOK, playlistResponse)
	}
}

func (h *Handler) GetAllPlaylists() gin.HandlerFunc {
	return func(c *gin.Context) {
		playlists, err := h.service.GetAllPlaylists(c.Request.Context())
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

func (h *Handler) GetPlaylistByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		playlists, err := h.service.GetPlaylistsByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error getting playlist by name", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlist by name"})
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

func (h *Handler) GetPlaylistByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			h.logger.Error("Invalid user ID", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		playlists, err := h.service.GetPlaylistsByUserID(c.Request.Context(), userID)
		if err != nil {
			h.logger.Error("Error getting playlist by user ID", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlist by user ID"})
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

		var playlist models.PlaylistRequest
		if err := c.ShouldBindJSON(&playlist); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
			return
		}

		playlistServ := models.Playlist{
			ID:     id,
			Name:   playlist.Name,
			UserID: playlist.UserID,
		}

		trackIDs := playlist.TrackIDs
		err = h.service.UpdatePlaylist(c.Request.Context(), &playlistServ, trackIDs)
		if err != nil {
			h.logger.Error("Error updating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating playlist"})
			return
		}

		response := models.MessageResponse{
			Message: "Updated playlist",
		}

		c.JSON(http.StatusOK, response)
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
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Playlist not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
				return
			}
			h.logger.Error("Error deleting playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting playlist"})
			return
		}

		message := models.MessageResponse{
			Message: "Deleted playlist",
		}

		c.JSON(http.StatusOK, message)
	}
}
