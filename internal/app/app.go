package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log/slog"
	_ "music-hosting/docs"
	config2 "music-hosting/internal/config"
	"music-hosting/internal/http/playlist"
	"music-hosting/internal/http/track"
	"music-hosting/internal/http/user"
	"music-hosting/internal/middleware"
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

	router.POST("/users/login", userHandler.Login())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userRoutes := router.Group("/api")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("/users", userHandler.GetAllUsers())
		userRoutes.GET("/users/:id", userHandler.GetUserID())
		userRoutes.GET("/users?offset=1&limit=10", userHandler.GetUserWithPagination())
		userRoutes.POST("/users", userHandler.CreateUser())
		userRoutes.PUT("/users/:id", userHandler.UpdateUser())
		userRoutes.DELETE("/users/:id", userHandler.DeleteUser())
	}

	trackRoutes := router.Group("/api")
	trackRoutes.Use(middleware.AuthMiddleware())
	{
		trackRoutes.GET("/tracks", trackHandler.GetAllTracks())
		trackRoutes.GET("/tracks/:id", trackHandler.GetTrackByID())
		trackRoutes.GET("/tracks/name/:name", trackHandler.GetTrackByName())
		trackRoutes.GET("/tracks/artist/:artist", trackHandler.GetTrackByArtist())
		trackRoutes.GET("/tracks?offset=1&limit=10", trackHandler.GetTracksWithPagination())
		trackRoutes.POST("/tracks", trackHandler.CreateTrack())
		trackRoutes.PUT("/tracks/:id", trackHandler.UpdateTrack())
		trackRoutes.DELETE("/tracks/:id", trackHandler.DeleteTrack())
	}

	playlistRoutes := router.Group("/api")
	playlistRoutes.Use(middleware.AuthMiddleware())
	{
		playlistRoutes.GET("/playlists", playlistHandler.GetAllPlaylists())
		playlistRoutes.GET("/playlists/:id", playlistHandler.GetPlaylistByID())
		playlistRoutes.GET("/playlists/name/:name", playlistHandler.GetPlaylistByName())
		playlistRoutes.GET("/playlists/userid/:id", playlistHandler.GetPlaylistByUserID())
		playlistRoutes.POST("/playlists", playlistHandler.CreatePlaylist())
		playlistRoutes.PUT("/playlists/:id", playlistHandler.UpdatePlaylist())
		playlistRoutes.DELETE("/playlists/:id", playlistHandler.DeletePlaylist())
	}

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := router.Run(serverAddr); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
