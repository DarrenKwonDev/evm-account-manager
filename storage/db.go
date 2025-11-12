package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Account struct {
	ID         int
	Address    string
	PrivateKey string
	Alias      string
	Chain      string
	Label      []string
	Memo       string
	totalValue float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DB struct {
	db *sql.DB
}

// 애플리케이션 코드는 스키마를 몰라야 함.
// 스키마 생성 및 수정 책임은 앱이 아닌 외부에 위치해야 한다
// 모든 스키마는 migration 내부 sql을 통해서만 진행되도록 한다
func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1) // sqlite single write

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}
