package learning

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stretchr/testify/require"
)

type ethTestData struct {
	Addresses   []string `json:"addresses"`
	PrivateKeys []string `json:"private_keys"`
}

func testClient() (client *ethclient.Client, testData *ethTestData, ganache *Ganache) {
	var err error
	client, err = ethclient.Dial("http://localhost:8545")
	if err != nil {
		panic(err)
	}

	optDataFilePath := os.Getenv("TEST_DATA_FILE")
	if optDataFilePath != "" {
		testData = readTestData(optDataFilePath)
	}

	ganache = NewGanache(false)
	return client, testData, ganache
}

type TestTransaction struct {
	From    string
	To      string
	Balance *big.Int
}

// func generateTransactions() []TestTransaction {
// 	return []TestTransaction{
// 		{
// 			From:    "",
// 			To:      "0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
// 			Balance: big.NewInt(100),
// 		}}
// }

// func testClientWithTransactions(transactions []TestTransaction) (*ethclient.Client, *ethTestData, *Ganache) {
// 	client, data, ganache := testClient()

// 	for _, tx := range transactions {
// 		// TODO: transactions
// 		_ = tx
// 	}
// 	return client, data, ganache
// }

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
	client, _, ganache := testClient()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], nil)
	require.NoError(t, err)

	oneThousandEthAsWei, ok := new(big.Int).SetString("1000000000000000000000", 0)
	require.True(t, ok)
	require.Equal(t, 0, oneThousandEthAsWei.Cmp(balance), "balance should be 1000 ETH")
}

func TestGetBalanceFirstBlock(t *testing.T) {
	client, _, ganache := testClient()

	addresses, err := ganache.AvailableAddresses()
	require.NoError(t, err)
	balance, err := client.BalanceAt(context.Background(), addresses[0], big.NewInt(0))
	require.NoError(t, err)

	require.Equalf(t, 1000.0, balanceToEther(balance), "balance should be 1000 ETH")
}

func TestHeaderByNumberLast(t *testing.T) {
	client, _, _ := testClient()

	lastHeader, err := client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, lastHeader)
	require.Greater(t, 1, lastHeader.Number.Cmp(big.NewInt(0)), "last block should be greater")

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)
	require.NotNil(t, firstHeader)
	require.Equal(t, int64(0), firstHeader.Number.Int64())

	fmt.Print(time.Unix(int64(firstHeader.Time), 0).Sub(time.Unix(int64(firstHeader.Time), 0)))
}
