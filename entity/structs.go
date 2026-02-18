package entity

import "time"

type ParamCreate struct {
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamCreateOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamUpdate struct {
	WalletID  string
	Balance   float64
	UpdatedAT time.Time
}

type ParamUpdateOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamGet struct {
	WalletID string
}

type ParamGetOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}

type ParamGetByAccount struct {
	AccountID string
}

type ParamGetByAccountOutput struct {
	WalletID  string
	AccountID string
	Balance   float64
	CreatedAT time.Time
	UpdatedAT time.Time
}
