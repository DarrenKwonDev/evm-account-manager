package service

import (
	"log"
	"strings"

	"drkup/account-tracker/db"
	"drkup/account-tracker/onchain"
)

// AccountService 계정 관리 서비스
type AccountService struct {
	db *db.DB
}

// 전역 인스턴스
var AccountServiceInstance *AccountService

// InitAccountService 서비스 초기화
func InitAccountService(db *db.DB) {
	AccountServiceInstance = &AccountService{
		db: db,
	}
}

// GetAccountService 전역 인스턴스 가져오기
func GetAccountService() *AccountService {
	return AccountServiceInstance
}

// NewAccountService creates a new account service
func NewAccountService(db *db.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

// CreateAccount 계정 생성 및 DB 저장
func (s *AccountService) CreateAccount(alias, chain, label, memo string) (*db.Account, error) {
	// 1. 온체인에서 계정 생성
	address, privateKey, err := onchain.CreateAccount()
	if err != nil {
		return nil, err
	}

	log.Printf("Generated account - Address: %s, PrivateKey: %s", address, privateKey)

	// 2. DB 저장을 위한 Account 객체 생성
	account := db.NewAccount()
	account.Address = address
	account.PrivateKey = privateKey
	account.Alias = alias
	account.Chain = chain
	account.Label = s.parseLabels(label)
	account.Memo = memo
	account.TotalValue = 0.0

	// 3. DB에 저장
	err = s.db.SaveAccount(account)
	if err != nil {
		log.Printf("Failed to save account to DB: %v", err)
		return nil, err
	}

	log.Printf("Account saved successfully - Alias: %s, Address: %s", alias, address)

	return account, nil
}

// parseLabels 레이블 파싱 (콤마로 구분)
func (s *AccountService) parseLabels(labelStr string) []string {
	if strings.TrimSpace(labelStr) == "" {
		return []string{}
	}

	labels := strings.Split(labelStr, ",")
	for i, label := range labels {
		labels[i] = strings.TrimSpace(label)
	}

	return labels
}
