package first

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func setupTestingWithAnvil(anvil *Anvil) (*ethclient.Client, *Anvil, func()) {
	client := ethclient.NewClient(anvil.Client())
	return client, anvil, func() {
		anvil.Close()
	}
}

func setupTesting(t *testing.T) (*ethclient.Client, *Anvil, func()) {
	anvil := StartAndConnect()
	return setupTestingWithAnvil(anvil)
}

func TestAnvilAPIIncreaseTimeAndBulkMineAllNewBlocksHaveSameTimestamp(t *testing.T) {
	client, anvil, tearDown := setupTesting(t)
	defer tearDown()

	blocksCount := 10
	blockchainDuration := time.Duration(blocksCount*12) * time.Second

	_, err := anvil.IncreaseTime(blockchainDuration)
	require.NoError(t, err)

	err = anvil.MineBlocks(blocksCount, blockchainDuration)
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

func TestAnvilAPICanControlBlockTimeByMiningBlockByBlock(t *testing.T) {
	client, anvil, tearDown := setupTesting(t)
	defer tearDown()

	blockTime := DefaultBlockTime
	_, err := anvil.IncreaseTime(blockTime)
	require.NoError(t, err)

	err = anvil.MineBlocks(1, blockTime)
	require.NoError(t, err)

	initialHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)

	firstHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
	// We expect that genesis block time is 0
	require.GreaterOrEqual(t, int(firstHeader.Time-initialHeader.Time), int(blockTime.Seconds()))

	delayBlocksMiningTime := time.Duration(15) * time.Second
	_, err = anvil.IncreaseTime(delayBlocksMiningTime)
	require.NoError(t, err)

	err = anvil.MineBlocks(1, blockTime)
	require.NoError(t, err)

	secondHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(2))
	require.NoError(t, err)
	require.Equal(t, int(delayBlocksMiningTime.Seconds()+blockTime.Seconds()), int(secondHeader.Time-firstHeader.Time))
}
