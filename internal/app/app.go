package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	config2 "music-hosting/internal/config"
	"music-hosting/internal/http-service/playlist"
	"music-hosting/internal/http-service/track"
	"music-hosting/internal/http-service/user"
	"music-hosting/internal/repository"
	"os"
)

func Run(config string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg, err := config2.LoadConfig(config)
	if err != nil {
		logger.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}

	userStorage, err := repository.NewUserStorage(cfg)
	if err != nil {
		logger.Error("Error creating user storage", slog.Any("error", err))
		os.Exit(1)
	}

	defer userStorage.Close()

	userHandler := user.NewHandler(userStorage, logger)

	router := gin.Default()

	userRoutes := router.Group("/api/user")
	{
		userRoutes.GET("/getall", userHandler.GetAllUsers())
		userRoutes.GET("/get/:id", userHandler.GetUserID())
		userRoutes.POST("/create", userHandler.CreateUser())
		userRoutes.PUT("/update/:id", userHandler.UpdateUser())
		userRoutes.DELETE("/delete/:id", userHandler.DeleteUser())
		//userRoutes.POST("/save/:id", userHandler.SaveTrack())
	}

	trackStorage, err := repository.NewTrackStorage(cfg)
	if err != nil {
		logger.Error("Error creating track storage", slog.Any("error", err))
		os.Exit(1)
	}

	defer trackStorage.Close()

	trackHandler := track.NewHandler(trackStorage, logger)

	trackRoutes := router.Group("/api/track")
	{
		trackRoutes.GET("/getall", trackHandler.GetAllTracks())
		trackRoutes.GET("/get/:id", trackHandler.GetTrackID())
		trackRoutes.POST("/create", trackHandler.CreateTrack())
		trackRoutes.PUT("/update/:id", trackHandler.UpdateTrack())
		trackRoutes.DELETE("/delete/:id", trackHandler.DeleteTrack())
	}

	playlistStorage, err := repository.NewPlaylistStorage(cfg)
	if err != nil {
		logger.Error("Error creating playlist storage", slog.Any("error", err))
		os.Exit(1)
	}

	defer playlistStorage.Close()

	playlistHandler := playlist.NewHandler(playlistStorage, logger)

	playlistRoutes := router.Group("/api/playlist")
	{
		playlistRoutes.GET("/getall", playlistHandler.GetAllPlaylists())
		playlistRoutes.GET("/get/:id", playlistHandler.GetPlaylistByID())
		playlistRoutes.POST("/create", playlistHandler.CreatePlaylist())
		playlistRoutes.PUT("/update/:id", playlistHandler.UpdatePlaylist())
		playlistRoutes.DELETE("/delete/:id", playlistHandler.DeletePlaylist())
	}

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := router.Run(serverAddr); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}

}
