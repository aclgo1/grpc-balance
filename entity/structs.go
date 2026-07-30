package entity

import "time"

type ParamCreate struct {
	AccountID string
	Balance   int64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamCreateOutput struct {
	WalletID  string
	AccountID string
	Balance   int64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamUpdate struct {
	WalletID  string
	Balance   int64
	UpdatedAT time.Time
}

type ParamUpdateOutput struct {
	WalletID  string
	AccountID string
	Balance   int64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamGet struct {
	WalletID string
}

type ParamGetOutput struct {
	WalletID  string
	AccountID string
	Balance   int64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamGetByAccount struct {
	AccountID string
}

type ParamGetByAccountOutput struct {
	WalletID  string
	AccountID string
	Balance   int64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamRegisterTransaction struct {
	ReferenceId string
	CreatedAt   time.Time
}

type ParamLedgerEntry struct{
	WalletId string
	ReferenceId string
	Amount int64
	CreatedAt time.Time
}

type ParamAuditWallet struct {
	WalletId string
}

type AuditReport struct {
	WalletID      string 
	WalletBalance int64  
	LedgerTotal   int64  
	Difference    int64  
	IsConsistent  bool  
}

type ParamReconcile struct {
	WalletId string
	Reason string
	Operator string
}

