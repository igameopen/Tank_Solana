package solana

import (
	"context"
	"fmt"
	"transferSrv/infra/config"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func SOLTransfer(wallet string, amount uint64) (solana.Signature, error) {
	// Create a new RPC client:
	rpcClient := rpc.New(config.GetCnf().Web3Cnf.SolanaAPIDomain)

	// Create a new WS client (used for confirming transactions)
	wsClient, err := ws.Connect(context.Background(), config.GetCnf().Web3Cnf.SolanaWSDomain)
	if err != nil {
		return [64]byte{}, err
	}

	// Load the account that you will send funds FROM:
	accountFrom, err := solana.PrivateKeyFromBase58(config.GetCnf().Web3Cnf.FromPrivateKey)
	if err != nil {
		return [64]byte{}, err
	}
	// fmt.Println("accountFrom private key:", accountFrom)
	// fmt.Println("accountFrom public key:", accountFrom.PublicKey())

	// The public key of the account that you will send sol TO:
	accountTo := solana.MustPublicKeyFromBase58(wallet)
	// The amount to send (in lamports);
	// 1 sol = 1000000000 lamports
	// amount := uint64(500000000)

	recent, err := rpcClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return [64]byte{}, err
	}

	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(
				amount,
				accountFrom.PublicKey(),
				accountTo,
			).Build(),
		},
		recent.Value.Blockhash,
		solana.TransactionPayer(accountFrom.PublicKey()),
	)
	if err != nil {
		return [64]byte{}, err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if accountFrom.PublicKey().Equals(key) {
				return &accountFrom
			}
			return nil
		},
	)
	if err != nil {
		return [64]byte{}, fmt.Errorf("unable to sign transaction: %w", err)
	}
	// spew.Dump(tx)
	// Pretty print the transaction:
	// tx.EncodeTree(text.NewTreeEncoder(os.Stdout, "Transfer SOL"))

	// Send transaction, and wait for confirmation:
	sig, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		rpcClient,
		wsClient,
		tx,
	)
	if err != nil {
		return [64]byte{}, err
	}
	spew.Dump(sig)

	return sig, nil
}
