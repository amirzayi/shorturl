package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"
)

var tableName = os.Getenv("PREFIX") + "url_shortener"

func pingAndMigrate(db *sql.DB, query string) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}
	return nil
}

type sqlDB struct {
	db *sql.DB
}

func (s sqlDB) Get(ctx context.Context, key string) (string, error) {
	var url string
	var expiredAt sql.NullTime

	row := s.db.QueryRowContext(ctx, fmt.Sprintf("SELECT url, expired_at FROM %s WHERE short_key = ?", tableName), key)

	err := row.Scan(&url, &expiredAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	if expiredAt.Valid && !expiredAt.Time.After(time.Now()) {
		return "", ErrNotFound
	}
	return url, nil
}

func (s sqlDB) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	var expiredAt *time.Time
	if ttl > 0 {
		t := time.Now().Add(ttl)
		expiredAt = &t
	}
	res, err := s.db.ExecContext(ctx, fmt.Sprintf("INSERT INTO %s(short_key, url, expired_at) VALUES(?, ?, ?)", tableName),
		key, value, expiredAt)
	if err != nil {
		return fmt.Errorf("failed to insert on database: %v", err)
	}

	insertedID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get inserted record id: %v", err)
	}
	if insertedID == 0 {
		return errors.New("inserted id is zero")
	}
	return nil
}
