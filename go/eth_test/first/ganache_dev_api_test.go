package first

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const debugTests = false

func setupTestingWithGanache(t *testing.T, ganache *Anvil) (*ethclient.Client, *Anvil, func()) {
	client := ethclient.NewClient(ganache.Client())
	return client, ganache, func() {
		ganache.Close()
	}
}

func setupTesting(t *testing.T) (*ethclient.Client, *Anvil, func()) {
	return setupTestingWithGanache(t, StartAndConnect())
}

func setupTestingWithBlocks(t *testing.T, duration time.Duration) (*ethclient.Client, *Anvil, func()) {
	return setupTestingWithGanache(t, NewGanacheWithStandardBlocks(int(duration/standardBlockDuration)))
}

// func debugTime(t *testing.T, client *ethclient.Client) {
// 	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
// 	if err == ethereum.NotFound {
// 		firstHeader, err = client.HeaderByNumber(context.Background(), big.NewInt(0))
// 		require.NoError(t, err)
// 	}
// 	lastHeader, err := client.HeaderByNumber(context.Background(), nil)
// 	require.NoError(t, err)
// 	fmt.Println("@dd RANGE [", firstHeader.Number, "](", time.Unix(int64(firstHeader.Time), 0), ") - [", lastHeader.Number, "](", time.Unix(int64(lastHeader.Time), 0), ")", (time.Duration(lastHeader.Time-firstHeader.Time) * time.Second).Seconds(), "seconds")
// }

func TestGanacheAPIIncreaseTimeAndBulkMineAllNewBlocksHaveSameTimestamp(t *testing.T) {
	client, ganache, tearDown := setupTesting(t)
	defer tearDown()

	blocksCount := 10
	blockchainDuration := time.Duration(blocksCount*12) * time.Second

	_, err := ganache.IncreaseTime(blockchainDuration)
	require.NoError(t, err)

	err = ganache.MineBlocks(blocksCount, blockchainDuration)
	require.NoError(t, err)

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), firstHeader.Number)

	lastHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(blocksCount)))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(int64(blocksCount)), lastHeader.Number)

	require.Greater(t, lastHeader.Time, firstHeader.Time)
	require.Equal(t, uint64((blocksCount-1)*int(blockchainDuration.Seconds())), lastHeader.Time-firstHeader.Time)

	middleHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(blocksCount/2)))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(int64(blocksCount/2)), middleHeader.Number)

	_, err = client.HeaderByNumber(context.Background(), big.NewInt(int64(blocksCount+1)))
	require.Error(t, err)
}

func TestGanacheAPICanControlBlockTimeByMiningBlockByBlock(t *testing.T) {
	client, ganache, tearDown := setupTesting(t)
	defer tearDown()

	blockTime := DefaultBlockTime
	_, err := ganache.IncreaseTime(blockTime)
	require.NoError(t, err)

	err = ganache.MineBlocks(1, blockTime)
	require.NoError(t, err)

	initialHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	// We expect that genesis block time is 0
	require.GreaterOrEqual(t, int(firstHeader.Time-initialHeader.Time), int(blockTime.Seconds()))

	delayBlocksMiningTime := time.Duration(15) * time.Second
	_, err = ganache.IncreaseTime(delayBlocksMiningTime)
	require.NoError(t, err)

	err = ganache.MineBlocks(1, blockTime)
	require.NoError(t, err)

	secondHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(2))
	require.NoError(t, err)
	require.Equal(t, int(delayBlocksMiningTime.Seconds()+blockTime.Seconds()), int(secondHeader.Time-firstHeader.Time))
}

func TestGanacheNewGanacheWithBlocks(t *testing.T) {
	blockchainDuration := 1 * time.Minute
	client, _, tearDown := setupTestingWithBlocks(t, blockchainDuration)
	defer tearDown()

	initialBlock, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)

	lastBlock, err := client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, blockchainDuration, time.Duration(lastBlock.Time-initialBlock.Time)*time.Second)
}

func TestGanacheNewGanacheWithBlocksBatchedDoNotWork(t *testing.T) {
	blockchainDuration := 1 * time.Minute
	blocksCount := int(blockchainDuration / standardBlockDuration)
	client, _, tearDown := setupTestingWithGanache(t, NewGanacheWithBlocks(func(blockNo int) (blockInfo *BlockInfo, stop bool) {
		if blockNo > blocksCount {
			return nil, true
		}
		return &BlockInfo{standardBlockDuration}, false
	}))
	defer tearDown()

	initialBlock, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)

	lastBlock, err := client.HeaderByNumber(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, initialBlock.Number, lastBlock.Number)
	require.Equal(t, time.Duration(0), time.Duration(lastBlock.Time-initialBlock.Time)*time.Second)
}
