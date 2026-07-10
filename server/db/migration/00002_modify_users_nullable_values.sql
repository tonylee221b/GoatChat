-- +goose Up
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_not_blank;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_shape;
DROP INDEX IF EXISTS users_email_unique_active;

