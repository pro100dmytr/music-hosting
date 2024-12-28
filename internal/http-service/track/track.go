package track

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models"
	"music-hosting/internal/repository"
	"music-hosting/pkg/utils/trackutils"
	"net/http"
	"strconv"
)

type TrackHandler interface {
	CreateTrack() gin.HandlerFunc
	GetTrackID() gin.HandlerFunc
	GetAllTracks() gin.HandlerFunc
	UpdateTrack() gin.HandlerFunc
	DeleteTrack() gin.HandlerFunc
}

type Handler struct {
	store  *repository.TrackStorage
	logger *slog.Logger
}

func NewHandler(store *repository.TrackStorage, logger *slog.Logger) *Handler {
	return &Handler{store: store, logger: logger}
}

func (h *Handler) CreateTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var track models.Track

		if err := c.ShouldBindJSON(&track); err != nil {
			h.logger.Error("Invalid request", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := trackutils.ValidateTrack(&track); err != nil {
			h.logger.Error("Error validating track", slog.Any("Error", "Error validating track"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := h.store.Create(c.Request.Context(), &track)
		if err != nil {
			h.logger.Error("Error creating track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"track": track})
	}
}

func (h *Handler) GetTrackID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid track ID", slog.Any("Error", "Invalid user ID"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		track, err := h.store.Get(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error getting track", slog.Any("Error", "Error getting track"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, track)
	}
}

func (h *Handler) GetAllTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracks, err := h.store.GetAll(c.Request.Context())
		if err != nil {
			h.logger.Error("Error getting all tracks", slog.Any("Error", "Error getting all tracks"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, tracks)
	}
}

func (h *Handler) UpdateTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid track ID", slog.Any("Error", "Invalid user ID"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		track := &models.Track{}
		if err := c.ShouldBindJSON(track); err != nil {
			h.logger.Error("Error parsing request body", slog.Any("Error", "Error parsing request body"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := trackutils.ValidateTrack(track); err != nil {
			h.logger.Error("Error validating track", slog.Any("Error", "Error validating track"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := h.store.Update(c.Request.Context(), track, id); err != nil {
			h.logger.Error("Error updating track", slog.Any("Error", "Error updating track"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"track update": track})
	}
}

func (h *Handler) DeleteTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid track ID", slog.Any("Error", "Invalid user ID"))
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := h.store.Delete(c.Request.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Error deleting track", slog.Any("Error", "Error deleting track"))
				c.JSON(http.StatusNotFound, gin.H{"error": err})
				return
			}

			h.logger.Error("Error deleting track", slog.Any("Error", "Error deleting track"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"track delete": nil})
	}
}
