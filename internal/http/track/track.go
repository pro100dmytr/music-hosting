package track

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

type Handler interface {
	CreateTrack() gin.HandlerFunc
	GetTrackID() gin.HandlerFunc
	GetAllTracks() gin.HandlerFunc
	GetTrackForName() gin.HandlerFunc
	GetTrackForArtist() gin.HandlerFunc
	GetTrackWithPagination() gin.HandlerFunc
	UpdateTrack() gin.HandlerFunc
	DeleteTrack() gin.HandlerFunc
}

type TrackHandler struct {
	service *service.TrackService
	logger  *slog.Logger
}

func NewTrackHandler(service *service.TrackService, logger *slog.Logger) *TrackHandler {
	return &TrackHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TrackHandler) CreateTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var track https.Track
		if err := c.ShouldBindJSON(&track); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		trackServ := services.Track{
			Name:   track.Name,
			Artist: track.Artist,
			URL:    track.URL,
		}

		err := h.service.CreateTrack(c.Request.Context(), &trackServ)
		if err != nil {
			h.logger.Error("Error creating track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating track"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Created track"})
	}
}

func (h *TrackHandler) GetTrackByID() gin.HandlerFunc {
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
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Track not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
				return
			}

			h.logger.Error("Error fetching track by ID", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by ID"})
			return
		}

		trackHttp := https.Track{
			ID:     track.ID,
			Name:   track.Name,
			Artist: track.Artist,
			URL:    track.URL,
		}

		c.JSON(http.StatusOK, gin.H{"track": trackHttp})
	}
}

func (h *TrackHandler) GetAllTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracks, err := h.service.GetAllTracks(c.Request.Context())
		if err != nil {
			h.logger.Error("Error fetching all tracks", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching all tracks"})
			return
		}

		var tracksHttp []https.Track
		for _, track := range tracks {
			trackHttp := https.Track{
				ID:     track.ID,
				Name:   track.Name,
				Artist: track.Artist,
				URL:    track.URL,
			}
			tracksHttp = append(tracksHttp, trackHttp)
		}

		c.JSON(http.StatusOK, gin.H{"tracks": tracksHttp})
	}
}

func (h *TrackHandler) UpdateTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid track ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
			return
		}

		var track https.Track
		if err := c.ShouldBindJSON(&track); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		trackServ := services.Track{
			ID:     id,
			Name:   track.Name,
			Artist: track.Artist,
			URL:    track.URL,
		}

		updatedTrack, err := h.service.UpdateTrack(c.Request.Context(), id, &trackServ)
		if err != nil {
			h.logger.Error("Error updating track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating track"})
			return
		}

		trackHttp := https.Track{
			ID:     updatedTrack.ID,
			Name:   updatedTrack.Name,
			Artist: updatedTrack.Artist,
			URL:    updatedTrack.URL,
		}

		c.JSON(http.StatusOK, gin.H{"track": trackHttp})
	}
}

func (h *TrackHandler) DeleteTrack() gin.HandlerFunc {
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

		c.JSON(http.StatusOK, gin.H{"message": "Track deleted successfully"})
	}
}

func (h *TrackHandler) GetTracksWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		tracks, err := h.service.GetTracksWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Error fetching tracks with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks with pagination"})
			return
		}

		var tracksHttp []https.Track
		for _, track := range tracks {
			trackHttp := https.Track{
				ID:     track.ID,
				Name:   track.Name,
				Artist: track.Artist,
				URL:    track.URL,
			}
			tracksHttp = append(tracksHttp, trackHttp)
		}

		c.JSON(http.StatusOK, gin.H{"tracks": tracksHttp})
	}
}

func (h *TrackHandler) GetTrackByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		tracks, err := h.service.GetTrackByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error fetching track by name", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by name"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tracks": tracks})
	}
}

func (h *TrackHandler) GetTrackByArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		artist := c.Param("artist")
		tracks, err := h.service.GetTrackByArtist(c.Request.Context(), artist)
		if err != nil {
			h.logger.Error("Error fetching tracks by artist", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks by artist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tracks": tracks})
	}
}
