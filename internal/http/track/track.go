package track

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/service"
	"music-hosting/pkg/utils/trackutils"
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
		var track models.Track
		if err := c.ShouldBindJSON(&track); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := trackutils.ValidateTrack(&track); err != nil {
			h.logger.Error("Error validating track", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error validating track"})
			return
		}

		createdTrack, err := h.service.CreateTrack(c.Request.Context(), &track)
		if err != nil {
			h.logger.Error("Error creating track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating track"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"track": createdTrack})
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
			h.logger.Error("Error fetching track by ID", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"track": track})
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

		c.JSON(http.StatusOK, gin.H{"tracks": tracks})
	}
}

func (h *TrackHandler) GetTrackByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		track, err := h.service.GetTrackByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error fetching track by name", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by name"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"track": track})
	}
}

func (h *TrackHandler) GetTrackByArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		artist := c.Param("artist")
		track, err := h.service.GetTrackByArtist(c.Request.Context(), artist)
		if err != nil {
			h.logger.Error("Error fetching track by artist", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by artist"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"track": track})
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

		var track models.Track
		if err := c.ShouldBindJSON(&track); err != nil {
			h.logger.Error("Invalid request body", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if err := trackutils.ValidateTrack(&track); err != nil {
			h.logger.Error("Error validating track", slog.Any("track", track))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error validating track"})
			return
		}

		err = h.service.UpdateTrack(c.Request.Context(), &track, id)
		if err != nil {
			h.logger.Error("Error updating track", slog.Int("id", id), slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating track"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"track": track})
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
			h.logger.Error("Error deleting track", slog.Int("id", id), slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting track"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Track deleted"})
	}
}

func (h *TrackHandler) GetTracksWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		minStr := c.DefaultQuery("min", "1")
		maxStr := c.DefaultQuery("max", "10")

		min, err := strconv.Atoi(minStr)
		if err != nil || min < 1 {
			h.logger.Error("Invalid pagination parameters", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
			return
		}

		max, err := strconv.Atoi(maxStr)
		if err != nil || max < 1 {
			h.logger.Error("Invalid pagination parameters", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
			return
		}

		tracks, err := h.service.GetTracksWithPagination(c.Request.Context(), max, (min-1)*max)
		if err != nil {
			h.logger.Error("Error fetching tracks with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks with pagination"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"tracks": tracks})
	}
}
