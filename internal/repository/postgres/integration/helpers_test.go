// Package postgres_test содержит интеграционные тесты репозиториев.
package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const defaultTestDSN = "postgres://pr_reviewer:pr_reviewer@localhost:5432/pr_reviewer_test?sslmode=disable"

// openTestDB открывает соединение с тестовой БД.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = defaultTestDSN
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		t.Fatalf("ping test db: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

// truncateAllTables очищает все таблицы домена.
func truncateAllTables(t *testing.T, db *sql.DB) {
	t.Helper()

	const query = `
		TRUNCATE TABLE
			pull_request_reviewers,
			pull_requests,
			users,
			teams
		RESTART IDENTITY CASCADE;
	`

	if _, err := db.Exec(query); err != nil {
		t.Fatalf("truncate tables: %v", err)
	}
}

// insertTeam добавляет тестовую команду в БД.
func insertTeam(t *testing.T, db *sql.DB, name string) {
	t.Helper()

	const query = `INSERT INTO teams (name) VALUES ($1);`

	if _, err := db.Exec(query, name); err != nil {
		t.Fatalf("insert team %s: %v", name, err)
	}
}

// insertUser добавляет тестового пользователя в БД.
func insertUser(
	t *testing.T,
	db *sql.DB,
	id, username, teamName string,
	isActive bool,
) {
	t.Helper()

	const query = `
		INSERT INTO users (id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4);
	`

	if _, err := db.Exec(query, id, username, teamName, isActive); err != nil {
		t.Fatalf("insert user %s: %v", id, err)
	}
}
