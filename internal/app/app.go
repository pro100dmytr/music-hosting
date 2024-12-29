package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	config2 "music-hosting/internal/config"
	"music-hosting/internal/http/playlist"
	"music-hosting/internal/http/track"
	"music-hosting/internal/http/user"
	"music-hosting/internal/repository"
	"music-hosting/internal/service"
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

	userService := service.NewUserService(userStorage, logger)
	userHandler := user.NewHandler(userService, logger)

	trackStorage, err := repository.NewTrackStorage(cfg)
	if err != nil {
		logger.Error("Error creating track storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer trackStorage.Close()

	trackService := service.NewTrackService(trackStorage, logger)
	trackHandler := track.NewTrackHandler(trackService, logger)

	playlistStorage, err := repository.NewPlaylistStorage(cfg)
	if err != nil {
		logger.Error("Error creating playlist storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer playlistStorage.Close()

	playlistService := service.NewPlaylistService(playlistStorage, logger)
	playlistHandler := playlist.NewPlaylistHandler(playlistService, logger)

	router := gin.Default()

	userRoutes := router.Group("/api/user")
	{
		userRoutes.GET("/getall", userHandler.GetAllUsers())
		userRoutes.GET("/get/:id", userHandler.GetUserID())
		userRoutes.GET("/getwith", userHandler.GetUserWithPagination())
		userRoutes.POST("/create", userHandler.CreateUser())
		userRoutes.PUT("/update/:id", userHandler.UpdateUser())
		userRoutes.DELETE("/delete/:id", userHandler.DeleteUser())
	}

	trackRoutes := router.Group("/api/track")
	{
		trackRoutes.GET("/getall", trackHandler.GetAllTracks())
		trackRoutes.GET("/get/id/:id", trackHandler.GetTrackByID())
		trackRoutes.GET("/get/name/:name", trackHandler.GetTrackByName())
		trackRoutes.GET("/get/artist/:artist", trackHandler.GetTrackByArtist())
		trackRoutes.GET("/getwith", trackHandler.GetTracksWithPagination())
		trackRoutes.POST("/create", trackHandler.CreateTrack())
		trackRoutes.PUT("/update/:id", trackHandler.UpdateTrack())
		trackRoutes.DELETE("/delete/:id", trackHandler.DeleteTrack())
	}

	playlistRoutes := router.Group("/api/playlist")
	{
		playlistRoutes.GET("/getall", playlistHandler.GetAllPlaylists())
		playlistRoutes.GET("/get/id/:id", playlistHandler.GetPlaylistByID())
		playlistRoutes.GET("/get/name/:name", playlistHandler.GetPlaylistByName())
		playlistRoutes.GET("/get/userid/:userid", playlistHandler.GetPlaylistByUserID())
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
