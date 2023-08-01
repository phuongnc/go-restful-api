
-- +migrate Up
CREATE TABLE IF NOT EXISTS "user" (
    "id" uuid PRIMARY KEY,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "age" INTEGER,
    "country" TEXT NOT NULL DEFAULT '',
    "avatar" TEXT NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT NULL,
    "last_active_at" timestamptz DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_email ON "user" ("email");
-- +migrate Down
DROP TABLE IF EXISTS "user";
