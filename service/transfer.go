package service

import (
	"math/big"

	"transferSrv/infra/library/log"
	"transferSrv/infra/library/solana"

	"github.com/shopspring/decimal"
)

func TransferSOL(wallet string, amount decimal.Decimal) (string, error) {

	if amount.LessThanOrEqual(decimal.NewFromBigInt(big.NewInt(0), 0)) {
		log.Errorf("TransferSOL amount LessThanOrEqual 0 :", amount)
	}

	sig, err := solana.SOLTransfer(wallet, amount.BigInt().Uint64())
	log.Info("TransferSOL signature : ", sig)

	return sig.String(), err
}

func TransferToken(wallet string, amount decimal.Decimal) (string, error) {

	if amount.LessThanOrEqual(decimal.NewFromBigInt(big.NewInt(0), 0)) {
		log.Errorf("TransferSOL amount LessThanOrEqual 0 :", amount)
	}

	sig, err := solana.TokenTransfer(wallet, amount.BigInt().Uint64())
	log.Info("TransferToken signature : ", sig)

	return sig.String(), err
}
