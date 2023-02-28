package learning

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestTransactionQueryAllInBlock(t *testing.T) {
	client, _, ganache, tearDown := testClient()
	defer tearDown()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	require.Greater(t, len(addresses), 3)

	err = generateTransactions(ganache,
		map[int][]TestTransaction{
			3: {
				{
					From:  addresses[0],
					To:    addresses[1],
					Value: big.NewInt(100),
				},
				{
					From:  addresses[2],
					To:    addresses[3],
					Value: big.NewInt(1000),
				},
			},
		})
	require.NoError(t, err)

	block, err := client.BlockByNumber(context.Background(), big.NewInt(3))
	require.NoError(t, err)

	count, err := client.TransactionCount(context.Background(), block.Hash())
	require.NoError(t, err)
	require.Equal(t, uint(2), count)

	// TODO validate transaction data
	for _, tx := range block.Transactions() {
		fmt.Println("@dd hash", tx.Hash().Hex())            // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
		fmt.Println("@dd value", tx.Value().String())       // 10000000000000000
		fmt.Println("@dd gas", tx.Gas())                    // 105000
		fmt.Println("@dd gasprice", tx.GasPrice().Uint64()) // 102000000000
		fmt.Println("@dd nonce", tx.Nonce())                // 110644
		fmt.Println("@dd data", tx.Data())                  // []
		fmt.Println("@dd to", tx.To().Hex())                // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e

		// read sender address
		chainID, err := client.NetworkID(context.Background())
		require.NoError(t, err)
		fmt.Println("@dd chainID", chainID)
		// The AsMessage method requires the EIP155 signer, which we derive the chain ID from the client.
		msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), nil)
		require.NoError(t, err)
		// TODO: check sender address
		fmt.Println("@dd from", msg.From().Hex()) // 0x4Bf2eDd5b99F84fe8406830DC0f8dDb09958C439

		// Each transaction has a receipt which contains the result of the execution of the transaction, such as any return values and logs, as well as the status which will be 1 (success) or 0 (fail).
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		require.NoError(t, err)
		// TODO
		fmt.Println("@dd receipt status", receipt.Status) // 1
		fmt.Println("@dd receipt logs", receipt.Logs)     // []
	}
}

func TestTransactionByHash(t *testing.T) {
	client, _, ganache, tearDown := testClient()
	defer tearDown()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	require.Greater(t, len(addresses), 3)

	err = generateTransactions(ganache,
		map[int][]TestTransaction{
			2: {
				{
					From:  addresses[0],
					To:    addresses[1],
					Value: big.NewInt(100),
				},
			},
		})
	require.NoError(t, err)

	block, err := client.BlockByNumber(context.Background(), big.NewInt(2))
	tx, isPending, err := client.TransactionByHash(context.Background(), block.Hash())
	require.NoError(t, err)
	require.False(t, isPending)
	// TODO: check tx to, frm, value
	fmt.Println("@dd tx to", tx.To().Hex()) // 0x55fE59D8Ad77035154dDd0AD0388D09Dd4047A8e
}

func TestTransactionCreate(t *testing.T) {
	client, testData, _, tearDown := testClient()
	defer tearDown()

	// load private key.
	privateKey, err := crypto.HexToECDSA(testData.PrivateKeys[0][2:])
	require.NoError(t, err)
	senderPublicKey := privateKey.Public()
	senderPublicKeyECDSA, ok := senderPublicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	fromAddress := crypto.PubkeyToAddress(*senderPublicKeyECDSA)

	// get the account nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	require.NoError(t, err)

	// get suggested gas price
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	fmt.Println("@dd gasPrice", gasPrice)

	// create transaction
	value := ethToWei(1)
	toAddress := common.HexToAddress(testData.Addresses[1])
	types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// sign the transaction with the private key of the sender
	// broadcast the transaction
}
