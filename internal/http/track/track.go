package track

import (
	"context"
	"log/slog"
	"music-hosting/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateTrack(ctx context.Context, track *models.Track) error
	GetTrackByID(ctx context.Context, id int) (*models.Track, error)
	GetTracks(ctx context.Context, name, artist string, playlistID, offset, limit int) ([]*models.Track, error)
	UpdateTrack(ctx context.Context, track *models.Track) error
	DeleteTrack(ctx context.Context, id int) error
}

type Handler struct {
	service Service
	logger  *slog.Logger
}

func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) CreateTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var track models.TrackRequest
		if err := c.ShouldBindJSON(&track); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		trackServ := models.Track{
			Name:     track.Name,
			Artist:   track.Artist,
			URL:      track.URL,
			Likes:    track.Likes,
			Dislikes: track.Dislikes,
		}

		err := h.service.CreateTrack(c.Request.Context(), &trackServ)
		if err != nil {
			h.logger.Error("Error creating track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating track"})
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) GetTrackByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid track ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
			return
		}

		track, err := h.service.GetTrackByID(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error fetching track by ID", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by ID"})
			return
		}

		trackResponse := models.TrackResponse{
			ID:       track.ID,
			Name:     track.Name,
			Artist:   track.Artist,
			URL:      track.URL,
			Likes:    track.Likes,
			Dislikes: track.Dislikes,
		}

		c.JSON(http.StatusOK, trackResponse)
	}
}

func (h *Handler) UpdateTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid track ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
			return
		}

		var track models.TrackRequest
		if err := c.ShouldBindJSON(&track); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		trackServ := models.Track{
			ID:       id,
			Name:     track.Name,
			Artist:   track.Artist,
			URL:      track.URL,
			Likes:    track.Likes,
			Dislikes: track.Dislikes,
		}

		err = h.service.UpdateTrack(c.Request.Context(), &trackServ)
		if err != nil {
			h.logger.Error("Error updating track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating track"})
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) DeleteTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid track ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
			return
		}

		err = h.service.DeleteTrack(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error deleting track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting track"})
			return
		}

		c.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) GetTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		artist := c.Query("artist")
		playlistIDStr := c.Query("playlistID")
		offsetStr := c.DefaultQuery("offset", "0")
		limitStr := c.DefaultQuery("limit", "10")

		var playlistID, offset, limit int
		var err error

		if playlistIDStr != "" {
			playlistID, err = strconv.Atoi(playlistIDStr)
			if err != nil {
				h.logger.Error("Invalid playlist ID", slog.Any("error", err))
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
				return
			}
		}

		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			h.logger.Error("Invalid offset", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
			return
		}

		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			h.logger.Error("Invalid limit", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
			return
		}

		tracks, err := h.service.GetTracks(c.Request.Context(), name, artist, playlistID, offset, limit)
		if err != nil {
			h.logger.Error("Error fetching tracks", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks"})
			return
		}

		var tracksResponse []models.TrackResponse
		for _, track := range tracks {
			trackResponse := models.TrackResponse{
				ID:       track.ID,
				Name:     track.Name,
				Artist:   track.Artist,
				URL:      track.URL,
				Likes:    track.Likes,
				Dislikes: track.Dislikes,
			}
			tracksResponse = append(tracksResponse, trackResponse)
		}

		c.JSON(http.StatusOK, tracksResponse)
	}
}
