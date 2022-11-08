-- +migrate Up

-- +migrate StatementBegin
CREATE TABLE "users" (
   id TEXT PRIMARY KEY,
   username TEXT NOT NULL UNIQUE,
   password TEXT NOT NULL,
   name TEXT NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +migrate StatementEnd

-- +migrate Down

DROP TABLE "users";