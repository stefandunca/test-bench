package learning

import (
	"context"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

const debugTests = false

func setupTestingWithGanache(t *testing.T, ganache *Ganache) (*ethclient.Client, *Ganache, func()) {
	client := ethclient.NewClient(ganache.Client())
	return client, ganache, func() {
		ganache.Close()
	}
}

func setupTesting(t *testing.T) (*ethclient.Client, *Ganache, func()) {
	return setupTestingWithGanache(t, NewGanache())
}

func setupTestingWithBlocks(t *testing.T, duration time.Duration) (*ethclient.Client, *Ganache, func()) {
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

	_, err = ganache.MineBlocks(blocksCount)
	require.NoError(t, err)

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), firstHeader.Number)

	lastHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(blocksCount)))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(int64(blocksCount)), lastHeader.Number)

	// All bulk mined blocks have the same timestamp
	require.Equal(t, firstHeader.Time, lastHeader.Time)

	middleHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(blocksCount/2)))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(int64(blocksCount/2)), middleHeader.Number)

	_, err = client.HeaderByNumber(context.Background(), big.NewInt(int64(blocksCount+1)))
	require.Error(t, err)
}

func TestGanacheAPICanControlBlockTimeByMiningBlockByBlock(t *testing.T) {
	client, ganache, tearDown := setupTesting(t)
	defer tearDown()

	_, err := ganache.IncreaseTime(time.Duration(12) * time.Second)
	require.NoError(t, err)

	_, err = ganache.MineBlocks(1)
	require.NoError(t, err)

	initialHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	require.Greater(t, int(firstHeader.Time-initialHeader.Time), 12)

	_, err = ganache.IncreaseTime(time.Duration(15) * time.Second)
	require.NoError(t, err)

	_, err = ganache.MineBlocks(1)
	require.NoError(t, err)

	secondHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(2))
	require.NoError(t, err)
	require.Equal(t, 15, int(secondHeader.Time-firstHeader.Time))
}

func TestGanacheCantUseAPIToGetTheGenesisBlockBackInTime(t *testing.T) {
	client, ganache, tearDown := setupTesting(t)
	defer tearDown()

	initialHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)

	// Go back 2 minutes
	offset, err := ganache.SetTime(time.Now().Add(time.Duration(-2) * time.Minute))
	require.NoError(t, err)
	require.Less(t, math.Abs(offset.Minutes()+2.0), 0.1, "We expect API confirm a negative offset of 2 minutes to account for error due to RPC call")

	// This doesn't change the initial header time
	initialHeaderAfterSetTime, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)

	require.Equal(t, initialHeader.Time, initialHeaderAfterSetTime.Time)

	_, err = ganache.MineBlocks(1)
	require.NoError(t, err)

	// API can't reset the genesis block time, only command line should work
	minedHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	require.Less(t, time.Duration(initialHeader.Time-minedHeader.Time)*time.Second, 2*time.Minute, "We expect a negative offset for the mined block equal to the passed block from where it started plus 2 minutes (due to SetTime call)")

	_, err = ganache.MineBlocks(1)
	require.NoError(t, err)

	secondBlockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(2))
	require.NoError(t, err)

	_, err = ganache.SetTime(time.Now().Add(-3 * time.Minute))
	require.NoError(t, err)

	_, err = ganache.MineBlocks(1)
	require.NoError(t, err)

	thirdAfterBackInTimeBlockHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(3))
	require.NoError(t, err)

	require.Equal(t, math.Abs((time.Duration(secondBlockHeader.Time-thirdAfterBackInTimeBlockHeader.Time) * time.Second).Seconds()), (1 * time.Minute).Seconds(), "This is absolute time so we got approximately one second over the previous offset")
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
	client, _, tearDown := setupTestingWithGanache(t, NewGanacheWithBlocksBatch(func(blockNo int) (blockInfo *BlockInfo, stop bool) {
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
