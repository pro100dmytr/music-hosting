package playlist

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models/https"
	"music-hosting/internal/models/services"
	"music-hosting/internal/service"
	"net/http"
	"strconv"
)

type PlaylistHandler interface {
	CreatePlaylist() gin.HandlerFunc
	GetPlaylistByID() gin.HandlerFunc
	GetAllPlaylists() gin.HandlerFunc
	GetPlaylistByName() gin.HandlerFunc
	GetPlaylistByUserID() gin.HandlerFunc
	UpdatePlaylist() gin.HandlerFunc
	DeletePlaylist() gin.HandlerFunc
}

type Handler struct {
	service *service.PlaylistService
	logger  *slog.Logger
}

func NewPlaylistHandler(service *service.PlaylistService, logger *slog.Logger) *Handler {
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
			Name:    playlist.Name,
			UserID:  playlist.UserID,
			TrackID: playlist.TrackID,
		}

		err := h.service.CreatePlaylist(c.Request.Context(), &playlistServ)
		if err != nil {
			h.logger.Error("Error creating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating playlist"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Playlist created"})
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

		playlistHttp := https.Playlist{
			ID:        playlist.ID,
			Name:      playlist.Name,
			UserID:    playlist.UserID,
			TrackID:   playlist.TrackID,
			CreatedAt: playlist.CreatedAt,
			UpdatedAt: playlist.UpdatedAt,
		}

		c.JSON(http.StatusOK, playlistHttp)
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

		var playlistsHttp []https.Playlist
		for _, playlist := range playlists {
			playlistHttp := https.Playlist{
				ID:        playlist.ID,
				Name:      playlist.Name,
				UserID:    playlist.UserID,
				TrackID:   playlist.TrackID,
				CreatedAt: playlist.CreatedAt,
				UpdatedAt: playlist.UpdatedAt,
			}

			playlistsHttp = append(playlistsHttp, playlistHttp)
		}

		c.JSON(http.StatusOK, playlistsHttp)
	}
}

func (h *Handler) GetPlaylistByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		playlists, err := h.service.GetPlaylistByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error getting playlist by name", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting playlist by name"})
			return
		}

		var playlistsHttp []https.Playlist
		for _, playlist := range playlists {
			playlistHttp := https.Playlist{
				ID:        playlist.ID,
				Name:      playlist.Name,
				UserID:    playlist.UserID,
				TrackID:   playlist.TrackID,
				CreatedAt: playlist.CreatedAt,
				UpdatedAt: playlist.UpdatedAt,
			}

			playlistsHttp = append(playlistsHttp, playlistHttp)
		}
		c.JSON(http.StatusOK, playlistsHttp)
	}
}

func (h *Handler) GetPlaylistByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Param("userID")
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

		var playlistsHttp []https.Playlist
		for _, playlist := range playlists {
			playlistHttp := https.Playlist{
				ID:        playlist.ID,
				Name:      playlist.Name,
				UserID:    playlist.UserID,
				TrackID:   playlist.TrackID,
				CreatedAt: playlist.CreatedAt,
				UpdatedAt: playlist.UpdatedAt,
			}

			playlistsHttp = append(playlistsHttp, playlistHttp)
		}

		c.JSON(http.StatusOK, playlistsHttp)
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
			TrackID:   playlist.TrackID,
			CreatedAt: playlist.CreatedAt,
			UpdatedAt: playlist.UpdatedAt,
		}

		updatedPlaylist, err := h.service.UpdatePlaylist(c.Request.Context(), id, &playlistServ)
		if err != nil {
			h.logger.Error("Error updating playlist", slog.Any("Error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating playlist"})
			return
		}

		playlistHttp := https.Playlist{
			ID:        updatedPlaylist.ID,
			Name:      updatedPlaylist.Name,
			UserID:    updatedPlaylist.UserID,
			TrackID:   updatedPlaylist.TrackID,
			CreatedAt: updatedPlaylist.CreatedAt,
			UpdatedAt: updatedPlaylist.UpdatedAt,
		}

		c.JSON(http.StatusOK, playlistHttp)
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

		c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted"})
	}
}
