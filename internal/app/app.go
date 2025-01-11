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

// TODO: Run function should return error instead of os.Exit(1)
func Run(configPath string) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // TODO: move default level to config.
	}))

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Error("Error loading config", slog.Any("error", err))
		os.Exit(1)
	}

	// TODO: create db connection here
	// TODO: pass DBConfig as a parameter instead of whole config
	db, err := postgresql.OpenConnection(cfg)
	// TODO: handle error

	// TODO: pass connection to repository
	userStorage, err := repository.NewUserStorage(cfg)
	if err != nil {
		logger.Error("Error creating user storage", slog.Any("error", err))
		os.Exit(1)
	}
	defer userStorage.Close()

	userService := service.NewUserService(userStorage, logger)
	userHandler := user.NewHandler(userService, logger)

	// TODO: pass sql.DB as parameter instead of creation one inside
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

	// TODO: use one group for all routes
	userRoutes := router.Group("/api/v1")
	userRoutes.Use(middleware.AuthMiddleware())

	userRoutes.GET("/users", userHandler.GetAllUsers())
	userRoutes.GET("/users/:id", userHandler.GetUserID())
	userRoutes.GET("/users?offset=1&limit=10", userHandler.GetUserWithPagination())
	userRoutes.POST("/users", userHandler.CreateUser())
	userRoutes.PUT("/users/:id", userHandler.UpdateUser())
	userRoutes.DELETE("/users/:id", userHandler.DeleteUser())

	trackRoutes := router.Group("/api")
	trackRoutes.Use(middleware.AuthMiddleware())
	trackRoutes.GET("/tracks", trackHandler.GetAllTracks())
	trackRoutes.GET("/tracks/:id", trackHandler.GetTrackByID())
	trackRoutes.GET("/tracks?name=<track_name>", trackHandler.GetTrackByName())
	trackRoutes.GET("/tracks?artist=<artist>", trackHandler.GetTrackByArtist())
	trackRoutes.GET("/tracks?offset=1&limit=10", trackHandler.GetTracksWithPagination())
	trackRoutes.POST("/tracks", trackHandler.CreateTrack())
	trackRoutes.PUT("/tracks/:id", trackHandler.UpdateTrack())
	trackRoutes.DELETE("/tracks/:id", trackHandler.DeleteTrack())
	trackRoutes.PATCH("/tracks/:id/like", trackHandler.AddLike())
	trackRoutes.DELETE("/tracks/:id/like", trackHandler.RemoveLike())
	trackRoutes.PATCH("/tracks/:id/dislike", trackHandler.AddDislike())
	trackRoutes.DELETE("/tracks/:id/dislike", trackHandler.RemoveDislike())

	playlistRoutes := router.Group("/api")
	playlistRoutes.Use(middleware.AuthMiddleware())
	{
		playlistRoutes.GET("/playlists", playlistHandler.GetAllPlaylists())
		playlistRoutes.GET("/playlists/:id", playlistHandler.GetPlaylistByID())
		playlistRoutes.GET("/playlists?name=<playlist_name>", playlistHandler.GetPlaylistByName())
		playlistRoutes.GET("/playlists?userid=<user_id>", playlistHandler.GetPlaylistByUserID())
		playlistRoutes.POST("/playlists", playlistHandler.CreatePlaylist())
		playlistRoutes.PUT("/playlists/:id", playlistHandler.UpdatePlaylist())
		playlistRoutes.DELETE("/playlists/:id", playlistHandler.DeletePlaylist())
	}

	if err = router.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		logger.Error("Failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
}
