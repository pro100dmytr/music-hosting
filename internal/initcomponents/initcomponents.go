package initcomponents

//type AppComponents struct {
//	UserHandler     *user.Handler
//	TrackHandler    *track.Handler
//	PlaylistHandler *playlist.Handler
//}
//
//func InitializeAppComponents(db *sql.DB, logger *slog.Logger) (*AppComponents, error) {
//	userStorage, err := repository.NewUserStorage(db)
//	if err != nil {
//		return nil, err
//	}
//	trackStorage, err := repository.NewTrackStorage(db)
//	if err != nil {
//		return nil, err
//	}
//	playlistStorage, err := repository.NewPlaylistStorage(db)
//	if err != nil {
//		return nil, err
//	}
//
//	userService := service.NewUserService(userStorage, logger)
//	trackService := service.NewTrackService(trackStorage, logger)
//	playlistService := service.NewPlaylistService(playlistStorage, logger)
//
//	userHandler := user.NewHandler(userService, logger)
//	trackHandler := track.NewTrackHandler(trackService, logger)
//	playlistHandler := playlist.NewPlaylistHandler(playlistService, logger)
//
//	return &AppComponents{
//		UserHandler:     userHandler,
//		TrackHandler:    trackHandler,
//		PlaylistHandler: playlistHandler,
//	}, nil
//}
