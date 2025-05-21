package app

import (
	"fmt"
	"log/slog"
	_ "music-hosting/docs"
	"music-hosting/internal/config"
	"music-hosting/internal/http/playlist"
	"music-hosting/internal/http/track"
	"music-hosting/internal/http/user"
	"music-hosting/internal/middleware"
	"music-hosting/internal/repository"
	"music-hosting/internal/service"
	"music-hosting/internal/storage/postgresql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func Run(configPath string) error {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	level := cfg.Logger.GetLogLevel()
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)

	db, err := postgresql.OpenConnection(&cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to create storage connection: %w", err)
	}
	defer db.Close()

	userStorage, err := repository.NewUserStorage(db)
	if err != nil {
		return fmt.Errorf("failed to create user storage: %w", err)
	}

	userSvc := service.NewUserService(userStorage, logger)
	userHandler := user.NewHandler(userSvc, logger)

	trackStorage, err := repository.NewTrackStorage(db)
	if err != nil {
		return fmt.Errorf("failed to create track storage: %w", err)
	}

	trackSvc := service.NewTrackService(trackStorage, logger)
	trackHandler := track.NewHandler(trackSvc, logger)

	playlistStorage, err := repository.NewPlaylistStorage(db)
	if err != nil {
		return fmt.Errorf("failed to create playlist storage: %w", err)
	}

	playlistSvc := service.NewPlaylistService(playlistStorage, logger)
	playlistHandler := playlist.NewHandler(playlistSvc, logger)

	router := gin.Default()

	router.POST("/users", userHandler.CreateUser())
	router.POST("/tracks", trackHandler.CreateTrack())
	router.POST("/playlists", playlistHandler.CreatePlaylist())
	router.POST("/login", userHandler.Login())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes := router.Group("/api/v1")
	routes.Use(middleware.Auth())
	{
		routes.GET("/users/:id", userHandler.GetUserID())
		routes.GET("/users", userHandler.GetUserWithPagination())
		routes.PUT("/users/:id", userHandler.UpdateUser())
		routes.DELETE("/users/:id", userHandler.DeleteUser())

		routes.GET("/tracks/:id", trackHandler.GetTrackByID())
		routes.GET("/tracks", trackHandler.GetTracks())
		routes.PUT("/tracks/:id", trackHandler.UpdateTrack())
		routes.DELETE("/tracks/:id", trackHandler.DeleteTrack())

		routes.GET("/playlists/:id", playlistHandler.GetPlaylistByID())
		routes.GET("/playlists", playlistHandler.GetPlaylists())
		routes.PUT("/playlists/:id", playlistHandler.UpdatePlaylist())
		routes.DELETE("/playlists/:id", playlistHandler.DeletePlaylist())
	}

	if err = router.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		return fmt.Errorf("Failed to start server: %w", err)
	}

	return nil
}
