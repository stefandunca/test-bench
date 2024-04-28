package first

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ethTestData struct {
	Addresses   []string `json:"addresses"`
	PrivateKeys []string `json:"private_keys"`
}

func testClient() (client *ethclient.Client, testData *ethTestData, anvil *Anvil, tearDown func()) {
	anvil = StartAndConnect()
	addresses, err := anvil.AvailableAddresses()
	if err != nil {
		panic(err)
	}
	strAddresses := make([]string, 0, len(addresses))
	for _, address := range addresses {
		strAddresses = append(strAddresses, address.Hex())
	}
	return anvil.EthClient(), &ethTestData{
			Addresses: strAddresses,
			PrivateKeys: []string{
				"0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
				"0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
				"0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
			},
		},
		anvil, func() {
			anvil.Close()
		}
}

func testClientWithBlocks() (client *ethclient.Client, testData *ethTestData, anvil *Anvil, tearDown func()) {
	anvil = NewAnvilWithStandardBlocks(10)
	return anvil.EthClient(), testData, anvil, func() {
		anvil.Close()
	}
}

type TestTransaction struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

func generateTransactions(anvil *Anvil, blocksWithTransactions map[int][]TestTransaction) error {
	processedBlocks := 0
	err := MineBlocks(anvil, func(blockNo int) (blockInfo *BlockInfo, stop bool) {
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

func weiInEthAsFloat() *big.Float {
	return big.NewFloat(math.Pow10(18))
}

func weiInEth() *big.Int {
	res, accuracy := weiInEthAsFloat().Int(big.NewInt(1))
	if accuracy != big.Exact {
		panic("Unexpected loss of digits")
	}
	return res
}

func ethToWei(eth int) *big.Int {
	return new(big.Int).Mul(big.NewInt(int64(eth)), weiInEth())
}

func balanceToEther(balance *big.Int) float64 {
	res, _ := new(big.Float).Quo(new(big.Float).SetInt(balance), weiInEthAsFloat()).Float64()
	return res
}
