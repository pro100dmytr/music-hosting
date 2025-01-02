package track

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"music-hosting/internal/models/dto"
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
	AddLike() gin.HandlerFunc
	RemoveLike() gin.HandlerFunc
	AddDislike() gin.HandlerFunc
	RemoveDislike() gin.HandlerFunc
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

		message := dto.MessageResponse{
			Message: "Created track",
		}

		c.JSON(http.StatusCreated, message)
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

		trackResponse := dto.TrackResponse{
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

func (h *TrackHandler) GetAllTracks() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracks, err := h.service.GetAllTracks(c.Request.Context())
		if err != nil {
			h.logger.Error("Error fetching all tracks", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching all tracks"})
			return
		}

		var tracksResponse []dto.TrackResponse
		for _, track := range tracks {
			trackResponse := dto.TrackResponse{
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
			ID:       id,
			Name:     track.Name,
			Artist:   track.Artist,
			URL:      track.URL,
			Likes:    track.Likes,
			Dislikes: track.Dislikes,
		}

		updatedTrack, err := h.service.UpdateTrack(c.Request.Context(), id, &trackServ)
		if err != nil {
			h.logger.Error("Error updating track", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating track"})
			return
		}

		trackResponse := dto.TrackResponse{
			ID:       updatedTrack.ID,
			Name:     updatedTrack.Name,
			Artist:   updatedTrack.Artist,
			URL:      updatedTrack.URL,
			Likes:    updatedTrack.Likes,
			Dislikes: updatedTrack.Dislikes,
		}

		c.JSON(http.StatusOK, gin.H{"track": trackResponse})
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

		message := dto.MessageResponse{
			Message: "Deleted track",
		}

		c.JSON(http.StatusOK, message)
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

		var tracksResponse []dto.TrackResponse
		for _, track := range tracks {
			trackResponse := dto.TrackResponse{
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

func (h *TrackHandler) GetTrackByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		tracks, err := h.service.GetTrackByName(c.Request.Context(), name)
		if err != nil {
			h.logger.Error("Error fetching track by name", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching track by name"})
			return
		}

		var tracksResponse []dto.TrackResponse
		for _, track := range tracks {
			trackResponse := dto.TrackResponse{
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

func (h *TrackHandler) GetTrackByArtist() gin.HandlerFunc {
	return func(c *gin.Context) {
		artist := c.Query("artist")
		tracks, err := h.service.GetTrackByArtist(c.Request.Context(), artist)
		if err != nil {
			h.logger.Error("Error fetching tracks by artist", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tracks by artist"})
			return
		}

		var tracksResponse []dto.TrackResponse
		for _, track := range tracks {
			trackResponse := dto.TrackResponse{
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

func (h *TrackHandler) AddLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.AddLike(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error adding like", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding like"})
			return
		}

		message := dto.MessageResponse{
			Message: "Added like",
		}

		c.JSON(http.StatusOK, message)
	}
}

func (h *TrackHandler) RemoveLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.RemoveLike(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error removing like", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing like"})
			return
		}

		message := dto.MessageResponse{
			Message: "Removed like",
		}

		c.JSON(http.StatusOK, message)
	}
}

func (h *TrackHandler) AddDislike() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.AddDislike(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error adding dislike", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding dislike"})
			return
		}

		message := dto.MessageResponse{
			Message: "Dislike like",
		}

		c.JSON(http.StatusOK, message)
	}
}

func (h *TrackHandler) RemoveDislike() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			h.logger.Error("Invalid ID", slog.Any("error", err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		err = h.service.RemoveDislike(c.Request.Context(), id)
		if err != nil {
			h.logger.Error("Error removing dislike", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing dislike"})
			return
		}

		message := dto.MessageResponse{
			Message: "Removed dislike",
		}

		c.JSON(http.StatusOK, message)
	}
}
