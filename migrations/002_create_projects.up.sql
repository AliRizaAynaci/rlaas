-- 002_create_projects.sql
CREATE TABLE IF NOT EXISTS projects (
                                        id          SERIAL PRIMARY KEY,
                                        name        TEXT   NOT NULL,
                                        api_key     TEXT   NOT NULL,
                                        created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- unique constraint for api_key and project name (global or later per-user)
ALTER TABLE projects
    ADD CONSTRAINT projects_api_key_unique UNIQUE(api_key),
    ADD CONSTRAINT projects_name_unique    UNIQUE(name);
