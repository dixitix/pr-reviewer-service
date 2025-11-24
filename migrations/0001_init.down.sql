-- 0001_init.down.sql
-- Откат инициализации схемы БД для сервиса назначения ревьюверов.

DROP TABLE IF EXISTS pull_request_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
