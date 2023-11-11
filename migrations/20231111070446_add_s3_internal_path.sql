-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS videos ADD COLUMN s3_int_path text NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS videos DROP COLUMN s3_int_path;
-- +goose StatementEnd
