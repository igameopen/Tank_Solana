package solana

import (
	"context"
	"fmt"
	"strconv"
	"transferSrv/infra/config"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	solanaapi "github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

func TokenTransfer(wallet string, amount uint64) (solana.Signature, error) {

	sender, err := solanaapi.PrivateKeyFromBase58(config.GetCnf().Web3Cnf.FromPrivateKey)
	receiver := solanaapi.MustPublicKeyFromBase58(wallet)
	tokenMint := solanaapi.MustPublicKeyFromBase58(config.GetCnf().Web3Cnf.TokenMint)

	httpClient := rpc.New(config.GetCnf().Web3Cnf.SolanaAPIDomain)
	ixs := []solana.Instruction{}
	tokenAta := solana.PublicKey{}

	//get list of ata associated to the token mint address
	atas, err := httpClient.GetTokenAccountsByOwner(
		context.TODO(),
		receiver,
		&rpc.GetTokenAccountsConfig{
			Mint: &tokenMint,
		},
		&rpc.GetTokenAccountsOpts{})

	//create new ata if none found
	if err != nil || len(atas.Value) == 0 {
		if err != nil {
			return [64]byte{}, err
		}

		ix, err := associatedtokenaccount.NewCreateInstruction(
			sender.PublicKey(),
			receiver,
			tokenMint).ValidateAndBuild()
		if err != nil {
			return [64]byte{}, err
		}
		ixs = append(ixs, ix)
		tokenAta = ix.Accounts()[1].PublicKey
	} else {
		//if token atas found i just get the first one
		tokenAta = atas.Value[0].Pubkey
	}

	//transfer instruction
	// amount := uint64(10 * math.Pow10(9)) //Sending n10 tokens of a token that has 9 decimals
	// amount = uint64(float64(amount) * math.Pow10(9)) //Sending n10 tokens of a token that has 9 decimals
	fromAccountAtas, err := httpClient.GetTokenAccountsByOwner(
		context.TODO(),
		sender.PublicKey(),
		&rpc.GetTokenAccountsConfig{
			Mint: &tokenMint,
		},
		&rpc.GetTokenAccountsOpts{})
	if err != nil {
		return [64]byte{}, err
	}
	if len(fromAccountAtas.Value) == 0 {
		return [64]byte{}, fmt.Errorf("sender does not have any tokens")
	}

	//get sender wallet token ata + check if sender have enough tokens
	fromAccountAta := solana.PublicKey{}
	hasEnoughTokens := true
	for _, ata := range fromAccountAtas.Value {
		hasEnoughTokens = true
		balance, err := httpClient.GetTokenAccountBalance(context.TODO(), ata.Pubkey, rpc.CommitmentFinalized)
		if err != nil {
			continue
		}
		a, err := strconv.ParseUint(balance.Value.Amount, 10, 64)
		if err != nil {
			continue
		}
		if a < amount {
			hasEnoughTokens = false
			continue
		}
		fromAccountAta = ata.Pubkey
		break
	}

	if fromAccountAta == solana.SystemProgramID {
		if !hasEnoughTokens {
			err = fmt.Errorf("sender does not have enough tokens")
		}
		return [64]byte{}, err
	}

	//transfer instruction
	ix := token.NewTransferInstruction(
		amount,
		fromAccountAta,
		tokenAta,
		sender.PublicKey(),
		[]solana.PublicKey{}).Build()
	ixs = append(ixs, ix)

	recentBlockHash, err := httpClient.GetRecentBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return [64]byte{}, err
	}
	tx, err := solana.NewTransaction(
		ixs,
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(sender.PublicKey()),
	)
	if err != nil {
		return [64]byte{}, err
	}

	//sign transaction
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if sender.PublicKey().Equals(key) {
				return &sender
			}
			return nil
		},
	)
	if err != nil {
		return [64]byte{}, err
	}

	//send transaction
	sig, err := httpClient.SendTransaction(
		context.TODO(),
		tx,
	)
	if err != nil {
		return [64]byte{}, err
	}

	spew.Dump(sig)

	// send transaction and wait for confirmation
	// import confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	// sig, err := confirm.SendAndConfirmTransaction(
	// 	context.TODO(),
	// 	httpClient,
	// 	wsClient,
	// 	tx,
	//   )
	return sig, nil
}
