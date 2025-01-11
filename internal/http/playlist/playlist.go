package playlist

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"music-hosting/internal/models/dto"
	"music-hosting/internal/models/https"
	"music-hosting/internal/models/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TODO: delete interface
type PlaylistHandler interface {
	CreatePlaylist() gin.HandlerFunc
	GetPlaylistByID() gin.HandlerFunc
	GetAllPlaylists() gin.HandlerFunc
	GetPlaylistByName() gin.HandlerFunc
	GetPlaylistByUserID() gin.HandlerFunc
	UpdatePlaylist() gin.HandlerFunc
	DeletePlaylist() gin.HandlerFunc
}

// TODO: add other methods
type Service interface {
	CreatePlaylist(ctx context.Context, playlist *services.Playlist) error
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

// TODO: rename to NewHandler
func NewPlaylistHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) CreatePlaylist() gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlist https.Playlist
		if err := c.ShouldBindJSON(&playlist); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
			return
		}

		playlistServ := services.Playlist{
			Name:   playlist.Name,
			UserID: playlist.UserID,
		}

		err := h.service.CreatePlaylist(c.Request.Context(), &playlistServ)
		if err != nil {
			h.logger.Error("Error creating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating playlist"})
			return
		}

		message := dto.MessageResponse{
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

		playlistResponse := dto.PlaylistResponse{
			ID:        playlist.ID,
			Name:      playlist.Name,
			UserID:    playlist.UserID,
			TrackID:   playlist.TracksID,
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

		var playlistsResponse []dto.PlaylistResponse
		for _, playlist := range playlists {
			playlistResponse := dto.PlaylistResponse{
				ID:        playlist.ID,
				Name:      playlist.Name,
				UserID:    playlist.UserID,
				TrackID:   playlist.TracksID,
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

		playlists, err := h.service.GetPlaylistByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error getting playlist by name", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlist by name"})
			return
		}

		var playlistsResponse []dto.PlaylistResponse
		for _, playlist := range playlists {
			playlistResponse := dto.PlaylistResponse{
				ID:        playlist.ID,
				Name:      playlist.Name,
				UserID:    playlist.UserID,
				TrackID:   playlist.TracksID,
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

		playlists, err := h.service.GetPlaylistByUserID(c.Request.Context(), userID)
		if err != nil {
			h.logger.Error("Error getting playlist by user ID", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlist by user ID"})
			return
		}

		var playlistsResponse []dto.PlaylistResponse
		for _, playlist := range playlists {
			playlistResponse := dto.PlaylistResponse{
				ID:        playlist.ID,
				Name:      playlist.Name,
				UserID:    playlist.UserID,
				TrackID:   playlist.TracksID,
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

		var playlist https.Playlist
		if err := c.ShouldBindJSON(&playlist); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing request body"})
			return
		}

		playlistServ := services.Playlist{
			Name:      playlist.Name,
			UserID:    playlist.UserID,
			TracksID:  playlist.TrackID,
			CreatedAt: playlist.CreatedAt,
			UpdatedAt: playlist.UpdatedAt,
		}

		updatedPlaylist, err := h.service.UpdatePlaylist(c.Request.Context(), id, &playlistServ)
		if err != nil {
			h.logger.Error("Error updating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating playlist"})
			return
		}

		playlistResponse := dto.PlaylistResponse{
			ID:        updatedPlaylist.ID,
			Name:      updatedPlaylist.Name,
			UserID:    updatedPlaylist.UserID,
			TrackID:   updatedPlaylist.TracksID,
			CreatedAt: updatedPlaylist.CreatedAt,
			UpdatedAt: updatedPlaylist.UpdatedAt,
		}

		c.JSON(http.StatusOK, playlistResponse)
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

		message := dto.MessageResponse{
			Message: "Deleted playlist",
		}

		c.JSON(http.StatusOK, message)
	}
}
