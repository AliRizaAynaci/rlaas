-- 003_add_user_id_to_projects.sql
ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE projects
    DROP CONSTRAINT IF EXISTS projects_name_unique,
    ADD CONSTRAINT projects_user_name_unique UNIQUE(user_id, name);
