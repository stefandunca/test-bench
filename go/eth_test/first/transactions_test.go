package first

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestTransactionQueryAllInBlock(t *testing.T) {
	t.Skip("TODO: implement generateTransactions")

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
	t.Skip("TODO: implement generateTransactions")
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

// See https://geth.ethereum.org/docs/developers/dapp-developer/native
func TestTransactionCreate(t *testing.T) {
	client, testData, _, tearDown := testClient()
	defer tearDown()

	ctx := context.Background()

	getEventsCount := func() int {
		blockNo, err := client.BlockNumber(ctx)
		require.NoError(t, err)
		lastNo := big.NewInt(int64(blockNo))

		q := ethereum.FilterQuery{
			FromBlock: lastNo,
			ToBlock:   lastNo,
			Topics:    [][]common.Hash{{common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")}},
		}
		logs, err := client.FilterLogs(context.Background(), q)
		require.NoError(t, err)
		return len(logs)
	}

	eventsCount := getEventsCount()
	require.Equal(t, 0, eventsCount)

	// load private key.
	privateKey, err := crypto.HexToECDSA(testData.PrivateKeys[0][2:])
	require.NoError(t, err)
	senderPublicKey := privateKey.Public()
	senderPublicKeyECDSA, ok := senderPublicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	fromAddress := crypto.PubkeyToAddress(*senderPublicKeyECDSA)

	// get the account nonce
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	require.NoError(t, err)

	gasLimit := uint64(21000)

	tipCap, err := client.SuggestGasTipCap(ctx)
	require.NoError(t, err)

	feeCap, err := client.SuggestGasPrice(ctx)
	require.NoError(t, err)

	// create transaction
	value := new(big.Int).Mul(big.NewInt(1), big.NewInt(params.Ether))
	toAddress := common.HexToAddress(testData.Addresses[1])
	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)

	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        &toAddress,
		Value:     value,
		Data:      nil,
	})
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	require.NoError(t, err)

	// broadcast transaction
	err = client.SendTransaction(ctx, signedTx)
	require.NoError(t, err)

	eventsCount = getEventsCount()
	require.Equal(t, 0, eventsCount)
}
