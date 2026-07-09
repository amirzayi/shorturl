package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var createURLShortenerPostgresTableQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
    id BIGSERIAL PRIMARY KEY,
    short_key VARCHAR(20) NOT NULL UNIQUE,
    url VARCHAR(255) NOT NULL,
    expired_at TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_url_shortener_key ON %s(short_key);
CREATE INDEX idx_url_shortener_expired_at ON %s(expired_at);`, tableName, tableName, tableName)

func NewPostgresShortener(dsn string) (Store, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	err = pingAndMigrate(db, createURLShortenerPostgresTableQuery)
	if err != nil {
		return nil, err
	}
	return sqlDB{db: db}, nil
}
