-- +goose Up
CREATE TABLE "users" (
   id TEXT PRIMARY KEY,
   username TEXT NOT NULL UNIQUE,
   password TEXT NOT NULL,
   name TEXT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE "users";
