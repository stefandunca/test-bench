package learning

import (
	"context"
	"encoding/json"
	"math"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stretchr/testify/require"
)

type ethTestData struct {
	Addresses   []string `json:"addresses"`
	PrivateKeys []string `json:"private_keys"`
}

func testClient() (client *ethclient.Client, testData *ethTestData, ganache *Ganache, tearDown func()) {
	optDataFilePath := os.Getenv("TEST_DATA_FILE")
	if optDataFilePath != "" {
		testData = readTestData(optDataFilePath)
	}

	ganache = NewGanache()
	return ganache.EthClient(), testData, ganache, func() {
		ganache.Close()
	}
}

func testClientWithBlocks() (client *ethclient.Client, testData *ethTestData, ganache *Ganache, tearDown func()) {
	optDataFilePath := os.Getenv("TEST_DATA_FILE")
	if optDataFilePath != "" {
		testData = readTestData(optDataFilePath)
	}

	ganache = NewGanacheWithStandardBlocks(10)
	return ganache.EthClient(), testData, ganache, func() {
		ganache.Close()
	}
}

type TestTransaction struct {
	From    common.Address
	To      common.Address
	Balance *big.Int
}

func generateTransactions(ganache *Ganache, blocksWithTransactions map[int][]TestTransaction) error {
	processedBlocks := 0
	err := MineBlocks(ganache, func(blockNo int) (blockInfo *BlockInfo, stop bool) {
		if processedBlocks == len(blocksWithTransactions) {
			return nil, true
		}
		if transactions, found := blocksWithTransactions[blockNo]; found {
			for _, _ /*tx*/ = range transactions {
				// TODO: send transaction
			}
			processedBlocks++
		}
		return &BlockInfo{standardBlockDuration}, false
	})

	return err
}

func balanceToEther(balance *big.Int) float64 {
	res, _ := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(math.Pow10(18))).Float64()
	return res
}

// Parse a json string list
func readTestData(json_file string) *ethTestData {
	// read json file
	file, err := os.Open(json_file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data *ethTestData
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
	return data
}

func TestEthClientDial(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8545")
	require.NoError(t, err)

	_ = client
}

func TestConvertAddresses(t *testing.T) {
	strAddress := "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"

	address := common.HexToAddress(strAddress)
	require.Equal(t, strAddress, address.Hex())
}

func TestGetBalanceLastBlock(t *testing.T) {
	client, _, ganache, tearDown := testClientWithBlocks()
	defer tearDown()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], nil)
	require.NoError(t, err)

	oneThousandEthAsWei, ok := new(big.Int).SetString("1000000000000000000000", 0)
	require.True(t, ok)
	require.Equal(t, 0, oneThousandEthAsWei.Cmp(balance), "balance should be 1000 ETH")
}

func TestGetBalanceFirstBlock(t *testing.T) {
	client, _, ganache, tearDown := testClientWithBlocks()
	defer tearDown()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], big.NewInt(0))
	require.NoError(t, err)

	require.Equalf(t, 1000.0, balanceToEther(balance), "balance should be 1000 ETH")
}

func TestHeaderByNumberLast(t *testing.T) {
	client, _, _, tearDown := testClientWithBlocks()
	defer tearDown()

	lastHeader, err := client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, lastHeader)
	require.Greater(t, lastHeader.Number.Cmp(big.NewInt(0)), 0, "last block should be greater")

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)
	require.NotNil(t, firstHeader)
	require.Equal(t, int64(0), firstHeader.Number.Int64())
	require.Greater(t, lastHeader.Time, firstHeader.Time, "last timestamp should be greater than")
}

func TestBlockByNumber(t *testing.T) {
	client, _, _, tearDown := testClientWithBlocks()
	defer tearDown()

	lastBlock, err := client.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, lastBlock)
	require.Greater(t, lastBlock.Number().Cmp(big.NewInt(0)), 0, "last block should be greater")

	firstBlock, err := client.BlockByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)
	require.NotNil(t, firstBlock)
	require.Equal(t, int64(0), firstBlock.Number().Int64())

	require.Greater(t, lastBlock.Time(), firstBlock.Time(), "last timestamp should be greater than")
	require.Equal(t, 0, len(lastBlock.Transactions()), "no transactions expected, mined empty blocks")
	require.True(t, strings.Contains(strings.ToLower(lastBlock.Hash().Hex()), "0x"))
	require.Equal(t, uint64(1), lastBlock.Difficulty().Uint64())
}

// TODO: enable after finishing generateTransactions
// func TestTransactionCount(t *testing.T) {
// 	client, _, ganache, tearDown := testClient()
// 	defer tearDown()

// 	addresses, err := ganache.AvailableAddresses()
// 	require.NoError(t, err)
// 	require.Greater(t, len(addresses), 3)

// 	err = generateTransactions(ganache,
// 		map[int][]TestTransaction{
// 			3: {
// 				{
// 					From:    addresses[0],
// 					To:      addresses[1],
// 					Balance: big.NewInt(100),
// 				},
// 				{
// 					From:    addresses[2],
// 					To:      addresses[3],
// 					Balance: big.NewInt(1000),
// 				},
// 			},
// 		})
// 	require.NoError(t, err)

// 	block, err := client.BlockByNumber(context.Background(), big.NewInt(3))
// 	require.NoError(t, err)

// 	count, err := client.TransactionCount(context.Background(), block.Hash())
// 	require.NoError(t, err)
// 	require.Equal(t, uint(2), count)
// }
