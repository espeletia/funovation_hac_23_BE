-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS clips (
    id SERIAL PRIMARY KEY,
    url text NOT NULL,
    videoID INTEGER NOT NULL,
    "order" INTEGER NOT NULL,
    thumbnail text NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (videoID) REFERENCES videos(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clips;
-- +goose StatementEnd
