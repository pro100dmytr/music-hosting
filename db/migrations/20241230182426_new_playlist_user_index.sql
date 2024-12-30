-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF EXISTS idx_playlists_user_id ON playlists(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_playlists_user_id;
-- +goose StatementEnd
