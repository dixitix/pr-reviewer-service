-- 0001_init.sql
-- Инициализация схемы БД для сервиса назначения ревьюверов.

CREATE TABLE teams (
    name text PRIMARY KEY,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id text PRIMARY KEY,
    username text NOT NULL,
    team_name text NOT NULL REFERENCES teams(name) ON DELETE RESTRICT,
    is_active boolean NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_team_active
    ON users (team_name, is_active);

CREATE TABLE pull_requests (
    id text PRIMARY KEY,
    name text NOT NULL,
    author_id text NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status text NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at timestamptz NOT NULL DEFAULT now(),
    merged_at timestamptz
);

CREATE INDEX idx_pull_requests_author
    ON pull_requests (author_id);

CREATE INDEX idx_pull_requests_status
    ON pull_requests (status);

CREATE TABLE pull_request_reviewers (
    pull_request_id text NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id text NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX idx_pr_reviewers_reviewer
    ON pull_request_reviewers (reviewer_id);
