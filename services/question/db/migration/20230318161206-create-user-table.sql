
-- +migrate Up
CREATE TABLE IF NOT EXISTS "user" (
    "id" uuid PRIMARY KEY,
    "name" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "age" INTEGER,
    "country" TEXT NOT NULL DEFAULT '',
    "avatar" TEXT NOT NULL,
    "is_active" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_deleted" BOOLEAN NOT NULL DEFAULT FALSE,
    "need_update_password" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamptz DEFAULT NULL,
    "last_active_at" timestamptz,
);


-- +migrate Down
DROP TABLE IF EXISTS "user";
