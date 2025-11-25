#!/usr/bin/env bash
set -euo pipefail

echo "Truncating tables in Postgres..."

docker compose exec db psql \
  -U pr_reviewer \
  -d pr_reviewer \
  -c "TRUNCATE TABLE pull_request_reviewers, pull_requests, users, teams RESTART IDENTITY CASCADE;"

echo "Done."
