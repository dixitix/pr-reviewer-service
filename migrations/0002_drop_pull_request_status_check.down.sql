-- 0002_drop_pull_request_status_check.down.sql
-- Возвращает ограничение CHECK на статус pull-request.

ALTER TABLE pull_requests
    ADD CONSTRAINT pull_requests_status_check CHECK (status IN ('OPEN', 'MERGED'));
