package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var createURLShortenerMysqlTableQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
id BIGINT PRIMARY KEY AUTO_INCREMENT,
short_key VARCHAR(20) NOT NULL UNIQUE,
url VARCHAR(255) NOT NULL,
expired_at DATETIME,
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
INDEX idx_key (short_key),
INDEX idx_expired_at (expired_at)
)ENGINE=MyISAM;`, tableName)

func NewMysqlShortener(dsn string) (Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = pingAndMigrate(db, createURLShortenerMysqlTableQuery)
	if err != nil {
		return nil, err
	}
	return sqlDB{db: db}, nil
}
