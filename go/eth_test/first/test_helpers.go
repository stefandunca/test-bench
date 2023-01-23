package learning

import (
	"encoding/json"
	"math"
	"math/big"
	"os"
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
