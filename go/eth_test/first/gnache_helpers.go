package learning

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Ganache struct {
	c                 *rpc.Client
	eth               *ethclient.Client
	initialSnapshotId int
}

func NewGanache() *Ganache {
	client, err := rpc.Dial("http://localhost:8545")
	panicOnError(err)
	ganache := &Ganache{client, ethclient.NewClient(client), 0}
	ganache.initialSnapshotId, err = ganache.TakeSnapshot()
	panicOnError(err)
	return ganache
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

type BlockInfo struct {
	blockDuration time.Duration
}

const standardBlockDuration = time.Duration(12 * time.Second)

func NewGanacheWithStandardBlocks(blocksCount int) *Ganache {
	return NewGanacheWithBlocks(func(blockNo int) (blockInfo *BlockInfo, stop bool) {
		if blockNo > blocksCount {
			return nil, true
		}
		return &BlockInfo{standardBlockDuration}, false
	})
}

type getBlockInfoCallback func(blockNo int) (blockInfo *BlockInfo, stop bool)

// NewGanacheWithBlocks will mine blocks based on information returned by blockInfo function
// TODO: start ganache from here with the genesis block time set. Use the command line to set the genesis block time then call this with current time for a proper approach
func NewGanacheWithBlocks(blockInfo getBlockInfoCallback) *Ganache {
	ganache := NewGanache()

	err := MineBlocks(ganache, blockInfo)
	panicOnError(err)

	return ganache
}

// NewGanacheWithBlocks will mine blocks based on information returned by blockInfo function
// TODO: start ganache from here with the genesis block time set. Use the command line to set the genesis block time then call this with current time for a proper approach
func MineBlocks(ganache *Ganache, blockInfo getBlockInfoCallback) error {
	lastHeader, err := ganache.eth.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	currentTime := time.Unix(int64(lastHeader.Time), 0)
	_, err = ganache.SetTime(currentTime)
	if err != nil {
		return err
	}

	blockNo := 1
	for {
		blockInfo, stop := blockInfo(blockNo)
		if stop {
			break
		}
		_, err := ganache.IncreaseTime(blockInfo.blockDuration)
		if err != nil {
			return err
		}

		_, err = ganache.MineBlocks(1)
		if err != nil {
			return err
		}

		currentTime = currentTime.Add(blockInfo.blockDuration)
		blockNo++
	}

	return nil
}

// NewGanacheWithBlocksBatch same as NewGanacheWithBlocks just that calls are batched
func NewGanacheWithBlocksBatch(blockInfo getBlockInfoCallback) *Ganache {
	ganache := NewGanache()

	header, err := ganache.eth.HeaderByNumber(context.Background(), big.NewInt(0))
	panicOnError(err)

	calls := make([]rpc.BatchElem, 0)

	currentTime := time.Unix(int64(header.Time), 0)
	calls = append(calls, prepareSetTimeCall(currentTime))
	blockNo := 1
	for {
		blockInfo, stop := blockInfo(blockNo)
		if stop {
			break
		}
		calls = append(calls, prepareIncreaseTimeCall(blockInfo.blockDuration))
		calls = append(calls, prepareMineCall(1))

		currentTime = currentTime.Add(blockInfo.blockDuration)
		blockNo++
	}

	err = ganache.c.BatchCall(calls)
	for _, call := range calls {
		if call.Error != nil {
			panic(err)
		}
	}
	return ganache
}

func (g *Ganache) Close() {
	g.RevertSnapshot(g.initialSnapshotId)
	g.c.Close()
}

func (g *Ganache) Client() *rpc.Client {
	return g.c
}

func (g *Ganache) EthClient() *ethclient.Client {
	return g.eth
}

func prepareCall(method string, params []any, result interface{}) rpc.BatchElem {
	return rpc.BatchElem{
		Method: method,
		Args:   params,
		Result: result,
	}
}

func (g *Ganache) TakeSnapshot() (snapshotId int, err error) {
	call := prepareCall("evm_snapshot", []any{}, new(string))
	response := new(interface{})
	err = g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseInt(strings.Replace(strings.ToLower((*response).(string)), "0x", "", -1), 16, 16)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (g *Ganache) RevertSnapshot(snapshotId int) error {
	call := prepareCall("evm_revert", []any{snapshotId}, new(bool))
	response := new(interface{})
	err := g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return err
	}
	if !(*response).(bool) {
		return errors.New("failed to revert snapshot")
	}
	return nil
}

func prepareSetTimeCall(newTime time.Time) rpc.BatchElem {
	return prepareCall("evm_setTime", []any{newTime.UnixMilli()}, new(float64))
}

// SetTime set the blockchain timestamp to a specific time
// Returns duration between the given timestamp and the current time.
// Beware - use this method cautiously as it allows you to move backwards in time, which may cause new blocks to appear
// to be mined before older blocks, thereby invalidating the blockchain state.
func (g *Ganache) SetTime(newTime time.Time) (offset time.Duration, err error) {
	call := prepareSetTimeCall(newTime)
	response := new(interface{})
	err = g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return 0, err
	}
	return time.Duration(int64((*response).(float64))) * time.Second, nil
}

func prepareIncreaseTimeCall(duration time.Duration) rpc.BatchElem {
	return prepareCall("evm_increaseTime", []any{"0x" + strconv.FormatInt(int64(duration.Seconds()), 16)}, new(float64))
}

// Increase the blockchain current timestamp by the specified amount of time in seconds
func (g *Ganache) IncreaseTime(duration time.Duration) (adjustedTime time.Duration, err error) {
	call := prepareIncreaseTimeCall(duration)
	response := new(interface{})
	err = g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return 0, err
	}
	return time.Duration(int64((*response).(float64))) * time.Second, nil
}

func prepareMineCall(blockCount int) rpc.BatchElem {
	return prepareCall("evm_mine", []any{map[string]any{"blocks": blockCount}}, new(string))
}

func (g *Ganache) MineBlocks(blockCount int) (int64, error) {
	call := prepareMineCall(blockCount)
	response := new(interface{})
	err := g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseInt(strings.Replace(strings.ToLower((*response).(string)), "0x", "", -1), 16, 16)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (g *Ganache) AvailableAddresses() ([]common.Address, error) {
	call := prepareCall("eth_accounts", []any{}, new(string))
	response := new(interface{})
	err := g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return nil, err
	}

	addresses := make([]common.Address, 0, len((*response).([]interface{})))
	for _, account := range (*response).([]interface{}) {
		addresses = append(addresses, common.HexToAddress(account.(string)))
	}
	return addresses, nil
}

// TODO add https://trufflesuite.github.io/ganache/#evm_addAccount endpoint support
