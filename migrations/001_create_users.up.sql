-- 001_create_users.sql
CREATE TABLE IF NOT EXISTS users (
                                     id          SERIAL PRIMARY KEY,
                                     google_id   TEXT   UNIQUE NOT NULL,
                                     email       TEXT   UNIQUE NOT NULL,
                                     created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
    );
