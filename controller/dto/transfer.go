package dto

import "github.com/shopspring/decimal"

//转换
type TransferReq struct {
	Wallet string          `json:"wallet" binding:"required"`
	Amount decimal.Decimal `json:"amount" binding:"required"`
}
