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
	TotalValue float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewAccount 새 Account 생성 (기본값 설정)
func NewAccount() *Account {
	return &Account{
		TotalValue: 0.0,
		Label:      []string{},
	}
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

// SaveAccount 계정 저장
func (d *DB) SaveAccount(account *Account) error {
	query := `
		INSERT INTO accounts (address, private_key, alias, chain, label, memo, total_value, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`

	// Label slice를 JSON으로 변환 (sqlite에는 배열 타입이 없으므로)
	labelsJSON := ""
	if len(account.Label) > 0 {
		// 간단한 콤마 구분 문자열로 저장
		for i, label := range account.Label {
			if i > 0 {
				labelsJSON += ","
			}
			labelsJSON += label
		}
	}

	_, err := d.db.Exec(query,
		account.Address,
		account.PrivateKey,
		account.Alias,
		account.Chain,
		labelsJSON,
		account.Memo,
		account.TotalValue,
	)

	return err
}
