-- 004_create_rate_limit_rules.sql
CREATE TABLE IF NOT EXISTS rate_limit_rules (
                                                id             SERIAL PRIMARY KEY,
                                                project_id     INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
                                                endpoint       TEXT    NOT NULL,
                                                strategy       TEXT    NOT NULL,
                                                key_by         TEXT    NOT NULL,
                                                limit_count    INTEGER NOT NULL,
                                                window_seconds INTEGER NOT NULL,
                                                created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
                                                UNIQUE(project_id, endpoint)
);
