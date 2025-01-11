-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS playlist_tracks (
    id SERIAL PRIMARY KEY,
    playlist_id INTEGER NOT NULL REFERENCES
    playlists(id) ON DELETE CASCADE, // TODO: make one-line statement
    track_id INTEGER NOT NULL REFERENCES
    tracks(id) ON DELETE CASCADE // TODO: make one-line statement
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS playlist_tracks;
-- +goose StatementEnd
