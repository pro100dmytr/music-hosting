package app

import (
	"fmt"
	"log/slog"
	_ "music-hosting/docs"
	"music-hosting/internal/config"
	"music-hosting/internal/database/postgresql"
	"music-hosting/internal/http/playlist"
	"music-hosting/internal/http/track"
	"music-hosting/internal/http/user"
	"music-hosting/internal/middleware"
	"music-hosting/internal/repository"
	"music-hosting/internal/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func Run(configPath string) error {
	serverCfg, dbCfg, loggerCfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	var level slog.Level
	switch loggerCfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)

	db, err := postgresql.OpenConnection(dbCfg)
	if err != nil {
		logger.Error("Error creating database connection", slog.Any("error", err))
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	defer db.Close()

	userStorage, err := repository.NewUserStorage(db)
	if err != nil {
		logger.Error("Error creating user storage", slog.Any("error", err))
		return fmt.Errorf("failed to create user storage: %w", err)
	}

	userService := service.NewUserService(userStorage, logger)
	userHandler := user.NewHandler(userService, logger)

	trackStorage, err := repository.NewTrackStorage(db)
	if err != nil {
		logger.Error("Error creating track storage", slog.Any("error", err))
		return fmt.Errorf("failed to create track storage: %w", err)
	}

	trackService := service.NewTrackService(trackStorage, logger)
	trackHandler := track.NewHandler(trackService, logger)

	playlistStorage, err := repository.NewPlaylistStorage(db)
	if err != nil {
		logger.Error("Error creating playlist storage", slog.Any("error", err))
		return fmt.Errorf("failed to create playlist storage: %w", err)
	}

	playlistService := service.NewPlaylistService(playlistStorage, logger)
	playlistHandler := playlist.NewHandler(playlistService, logger)

	router := gin.Default()

	router.POST("/users/create", userHandler.CreateUser())
	router.POST("/tracks/create", trackHandler.CreateTrack())
	router.POST("/playlists/create", playlistHandler.CreatePlaylist())
	router.POST("/users/login", userHandler.Login())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	Routes := router.Group("/api/v1")
	Routes.Use(middleware.AuthMiddleware())
	{
		Routes.GET("/users", userHandler.GetAllUsers())
		Routes.GET("/users/:id", userHandler.GetUserID())
		Routes.GET("/users?offset=1&limit=10", userHandler.GetUserWithPagination())
		Routes.PUT("/users/:id", userHandler.UpdateUser())
		Routes.DELETE("/users/:id", userHandler.DeleteUser())

		Routes.GET("/tracks", trackHandler.GetAllTracks())
		Routes.GET("/tracks/:id", trackHandler.GetTrackByID())
		Routes.GET("/tracks?name=<track_name>", trackHandler.GetTracksByName())
		Routes.GET("/tracks?artist=<artist>", trackHandler.GetTrackByArtist())
		Routes.GET("/tracks?playlistID=<playlistID>", trackHandler.GetTracksByPlaylistID())
		Routes.GET("/tracks?offset=1&limit=10", trackHandler.GetTracksWithPagination())
		Routes.PUT("/tracks/:id", trackHandler.UpdateTrack())
		Routes.DELETE("/tracks/:id", trackHandler.DeleteTrack())

		Routes.GET("/playlists", playlistHandler.GetAllPlaylists())
		Routes.GET("/playlists/:id", playlistHandler.GetPlaylistByID())
		Routes.GET("/playlists?name=<playlist_name>", playlistHandler.GetPlaylistByName())
		Routes.GET("/playlists?userid=<user_id>", playlistHandler.GetPlaylistByUserID())
		Routes.PUT("/playlists/:id", playlistHandler.UpdatePlaylist())
		Routes.DELETE("/playlists/:id", playlistHandler.DeletePlaylist())
	}

	if err = router.Run(fmt.Sprintf(":%s", serverCfg.Port)); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}
