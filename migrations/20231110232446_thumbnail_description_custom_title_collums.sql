-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS videos ADD COLUMN thumbnail text NOT NULL DEFAULT '';
ALTER TABLE IF EXISTS videos ADD COLUMN description text NOT NULL DEFAULT '';
ALTER TABLE IF EXISTS videos ADD COLUMN custom_title text NOT NULL DEFAULT '';
ALTER TABLE IF EXISTS videos ALTER COLUMN url TYPE text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Undo the changes for column 'thumbnail'
ALTER TABLE IF EXISTS videos DROP COLUMN IF EXISTS thumbnail;

-- Undo the changes for column 'description'
ALTER TABLE IF EXISTS videos DROP COLUMN IF EXISTS description;

-- Undo the changes for column 'custom_title'
ALTER TABLE IF EXISTS videos DROP COLUMN IF EXISTS custom_title;

-- Undo the changes for column 'url'
-- Change the data type of the column back to the original type
ALTER TABLE videos ALTER COLUMN column_name TYPE VARCHAR(255);


-- +goose StatementEnd
