-- 0002_drop_pull_request_status_check.up.sql
-- Удаляет ограничение CHECK на статус pull-request.

ALTER TABLE pull_requests
    DROP CONSTRAINT IF EXISTS pull_requests_status_check;
