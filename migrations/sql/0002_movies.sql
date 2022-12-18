-- +goose Up
CREATE TABLE "movies" (
   id TEXT PRIMARY KEY,
   name TEXT NOT NULL,
   description TEXT NOT NULL,
   link TEXT NOT NULL,
   thumbnail TEXT NOT NULL,
   shared_by TEXT NOT NULL REFERENCES users(id),
   shared_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMPTZ
);

-- +goose Down
DROP TABLE "movies";