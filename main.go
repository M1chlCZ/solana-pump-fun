package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"golang.org/x/time/rate"
	"log"
	"strings"
	"time"
)

var (
	signatures = make(map[string]bool)
	rpcClient  = rpc.NewWithCustomRPCClient(rpc.NewWithLimiter(
		"https://solana.api.onfinality.io/public",
		rate.Every(time.Second), // time frame
		5,                       // limit of requests per time frame
	))
	pfProgramKey = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	currentBlock = 0
	syncBlock    = 0
)

func main() {
	fmt.Println("Listening for logs...")
	//go logSubscribe()
	ctx := context.Background()
	client, err := ws.Connect(context.Background(), rpc.MainNetBeta.WS)
	if err != nil {
		panic(err)
	}
	program := solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P") // serum
	{
		// Subscribe to log events that mention the provided pubkey:
		sub, err := client.LogsSubscribeMentions(
			program,
			rpc.CommitmentRecent,
		)
		if err != nil {
			fmt.Println("err:", err)
		}
		defer sub.Unsubscribe()

		for {
			got, err := sub.Recv(ctx)
			if err != nil {
				fmt.Println("err:", err)
				continue
			}
			go getTx(got.Value.Signature)
			//spew.Dump(got)
		}
	}
}

func getTx(signature solana.Signature) {
	maxVersion := uint64(0)
	opts := rpc.GetTransactionOpts{
		MaxSupportedTransactionVersion: &maxVersion,
		Commitment:                     rpc.CommitmentFinalized,
	}
	tx, err := rpcClient.GetTransaction(context.Background(), signature, &opts)
	if err != nil {
		log.Printf("Failed to get tx %s: %v", signature.String(), err)
		return
	}

	transaction, err := tx.Transaction.GetTransaction()
	if err != nil {
		fmt.Printf("Failed to parse transaction: %v \n", err)
		return
	}
	sl, err := transaction.Message.AccountMetaList()
	if err != nil {
		fmt.Printf("Failed to get account meta list: %v \n", err)
		return
	}

	decodeTransaction(tx.Meta.LogMessages, transaction.Message.Header.NumRequiredSignatures, sl, tx.Meta.PreTokenBalances, tx.Meta.PostTokenBalances)
}

func decodeTransaction(tx []string, numRequiredSignatures uint8, accountMeta solana.AccountMetaSlice, pre, post []rpc.TokenBalance) {

	fmt.Println("Decoding transaction:")
	spew.Dump(tx)
	for _, log := range tx {
		hasBuy := strings.Contains(log, "Instruction: Buy")
		hasSell := strings.Contains(log, "Instruction: Sell")
		hasCreate := strings.Contains(log, "Create")
		signerCount := numRequiredSignatures
		signers := accountMeta[0:signerCount]
		//const signers = message.accountKeys.slice(0, signerCount);
		if hasBuy && !hasSell && !hasCreate {
			for i, signer := range signers {
				preBalance := pre[i].UiTokenAmount.UiAmount
				postBalance := post[i].UiTokenAmount.UiAmount
				if preBalance != nil && postBalance != nil {
					substraction := *preBalance - *postBalance
					fmt.Printf("buyer: %s, buy amount: %d, failed signatures: %d \n", signer.PublicKey.String(), substraction/1000000000, len(signatures))
				}
			}
		}
	}
}
