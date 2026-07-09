package store

import (
	"database/sql"
	"fmt"

	_ "github.com/ncruces/go-sqlite3/driver"
)

var createURLShortenerSqliteTableQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
id INTEGER primary key AUTOINCREMENT,
short_key TEXT,
url TEXT,
expired_at TEXT,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`, tableName)

func NewSqliteShortener(dsn string) (Store, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	err = pingAndMigrate(db, createURLShortenerSqliteTableQuery)
	if err != nil {
		return nil, err
	}
	return sqlDB{db: db}, nil
}
