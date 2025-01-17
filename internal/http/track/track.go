package track

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
	CreateTrack(ctx context.Context, track *models.Track) error
	GetTrackByID(ctx context.Context, id int) (*models.Track, error)
	GetAllTracks(ctx context.Context) ([]*models.Track, error)
	GetTracksByName(ctx context.Context, name string) ([]*models.Track, error)
	GetTracksByArtist(ctx context.Context, artist string) ([]*models.Track, error)
	GetTracksByPlaylistID(ctx context.Context, playlistID int) ([]*models.Track, error)
	GetTracksWithPagination(ctx context.Context, offset, limit string) ([]*models.Track, error)
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

		message := models.MessageResponse{
			Message: "Created track",
		}

		c.JSON(http.StatusCreated, message)
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
			if errors.Is(err, sql.ErrNoRows) {
				h.logger.Error("Track not found", slog.Any("error", err))
				c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
				return
			}

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

func (h *Handler) GetAllTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracks, err := h.service.GetAllTracks(c.Request.Context())
		if err != nil {
			h.logger.Error("Error fetching all tracks", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching all tracks"})
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

		response := models.MessageResponse{
			Message: "Updated track",
		}

		c.JSON(http.StatusOK, response)
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

		message := models.MessageResponse{
			Message: "Deleted track",
		}

		c.JSON(http.StatusOK, message)
	}
}

func (h *Handler) GetTracksWithPagination() gin.HandlerFunc {
	return func(c *gin.Context) {
		offset := c.DefaultQuery("offset", "0")
		limit := c.DefaultQuery("limit", "10")

		tracks, err := h.service.GetTracksWithPagination(c.Request.Context(), limit, offset)
		if err != nil {
			h.logger.Error("Error fetching tracks with pagination", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks with pagination"})
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

func (h *Handler) GetTracksByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		tracks, err := h.service.GetTracksByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error fetching track by name", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by name"})
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

func (h *Handler) GetTrackByArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		artist := c.Query("artist")
		tracks, err := h.service.GetTracksByArtist(c.Request.Context(), artist)
		if err != nil {
			h.logger.Error("Error fetching tracks by artist", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks by artist"})
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

func (h *Handler) GetTracksByPlaylistID() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Query("playlistID")
		playlistID, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid playlist ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
			return
		}

		tracks, err := h.service.GetTracksByPlaylistID(c.Request.Context(), playlistID)
		if err != nil {
			h.logger.Error("Error fetching tracks by playlist", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks by playlist"})
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
