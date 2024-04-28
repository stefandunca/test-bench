package first

import (
	"errors"
	"fmt"
	"math/big"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	DefaultBlockTime = 12 * time.Second
)

const (
	anvilPort = 8545
)

type Anvil struct {
	c                 *rpc.Client
	eth               *ethclient.Client
	initialSnapshotId int
	cmd               *exec.Cmd
}

func isAnvilWorking() bool {
	client, err := rpc.Dial(fmt.Sprintf("http://localhost:%d", anvilPort))
	if err != nil {
		return false
	}
	defer client.Close()

	var accounts []string
	err = client.Call(&accounts, "eth_accounts")
	return err == nil && len(accounts) == 10
}

func waitForAnvilToStart(client *rpc.Client, timeout time.Duration) (err error) {
	start := time.Now()
	for time.Since(start) < timeout {
		if isAnvilWorking() {
			err = client.Call(nil, "eth_chainId")
			if err == nil {
				return nil
			}
			time.Sleep(5 * time.Millisecond)
		}
	}
	return fmt.Errorf("anvil did not start in %s; last error: %w", timeout, err)
}

func startAnvil() (*exec.Cmd, error) {
	cmd := exec.Command("anvil", "--port", strconv.Itoa(anvilPort), "--timestamp", "1713900000", "--no-mining", "--silent")
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func stopAnvil() error {
	cmd := exec.Command("killall", "anvil")
	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func StartAndConnect() *Anvil {
	var cmd *exec.Cmd
	var err error
	if isAnvilWorking() {
		stopAnvil()
	}

	cmd, err = startAnvil()
	panicOnError(err)

	client, err := rpc.Dial(fmt.Sprintf("http://localhost:%d", anvilPort))
	panicOnError(err)
	err = waitForAnvilToStart(client, 1*time.Second)
	panicOnError(err)
	anvil := &Anvil{
		c:   client,
		eth: ethclient.NewClient(client),
		cmd: cmd,
	}

	time.Sleep(100 * time.Millisecond)
	anvil.initialSnapshotId, err = anvil.TakeSnapshot()
	panicOnError(err)

	return anvil
}

func (g *Anvil) StopAllInstances() {
	if g.cmd != nil {
		err := stopAnvil()
		panicOnError(err)
		err = g.cmd.Wait()
		panicOnError(err)
	}
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

func NewAnvilWithStandardBlocks(blocksCount int) *Anvil {
	return NewAnvilWithBlocks(func(blockNo int) (blockInfo *BlockInfo, stop bool) {
		if blockNo > blocksCount {
			return nil, true
		}
		return &BlockInfo{standardBlockDuration}, false
	})
}

// getBlockInfoCallback should return stop == true and blockInfo == nil if should stop otherwise return blockInfo
type getBlockInfoCallback func(blockNo int) (blockInfo *BlockInfo, stop bool)

// NewAnvilWithBlocks will mine blocks based on information returned by blockInfo function
func NewAnvilWithBlocks(blockInfo getBlockInfoCallback) *Anvil {
	anvil := StartAndConnect()

	err := MineBlocks(anvil, blockInfo)
	panicOnError(err)

	return anvil
}

// MineBlocks will mine blocks based on information returned by getBlockInfo function
func MineBlocks(anvil *Anvil, getBlockInfo getBlockInfoCallback) error {
	blockNo := 1
	blockCountInSlice := 0
	prevBlockTime := DefaultBlockTime
	stop := false
	var blockInfo *BlockInfo
	calls := make([]rpc.BatchElem, 0)
	for !stop {
		blockInfo, stop = getBlockInfo(blockNo)
		if !stop {
			blockCountInSlice++
		}

		// If stop then blockInfo is nil
		if stop || prevBlockTime != blockInfo.blockDuration {
			if blockCountInSlice > 0 {
				calls = append(calls, prepareMineCall(blockCountInSlice, prevBlockTime))
				if !stop {
					prevBlockTime = blockInfo.blockDuration
					blockCountInSlice = 0
				}
			}
		}
		blockNo++
	}

	err := anvil.c.BatchCall(calls)
	if err != nil {
		return err
	}

	for _, call := range calls {
		if call.Error != nil {
			return fmt.Errorf("error mining blocks; first error: %w", call.Error)
		}
	}
	return nil
}

func (g *Anvil) Close() {
	err := g.RevertSnapshot(g.initialSnapshotId)
	panicOnError(err)
	g.StopAllInstances()
	g.c.Close()
}

func (g *Anvil) Client() *rpc.Client {
	return g.c
}

func (g *Anvil) EthClient() *ethclient.Client {
	return g.eth
}

func prepareCall(method string, args []any, result interface{}) rpc.BatchElem {
	return rpc.BatchElem{
		Method: method,
		Args:   args,
		Result: result,
	}
}

func (g *Anvil) TakeSnapshot() (snapshotId int, err error) {
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

func (g *Anvil) RevertSnapshot(snapshotId int) error {
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

// `evm_setTime` didn't seem to be working as expected with Anvil
// func prepareSetTimeCall(newTime time.Time) rpc.BatchElem {
// 	return prepareCall("evm_setTime", []any{newTime.UnixMilli()}, new(float64))
// }
//
// SetTime set the blockchain timestamp to a specific time
// Returns duration between the given timestamp and the current time.
// func (g *Anvil) SetTime(newTime time.Time) (offset time.Duration, err error) {
// 	call := prepareSetTimeCall(newTime)
// 	response := new(interface{})
// 	err = g.c.Call(response, call.Method, call.Args...)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return time.Duration(int64((*response).(float64))) * time.Second, nil
// }

func prepareIncreaseTimeCall(duration time.Duration) rpc.BatchElem {
	return prepareCall("evm_increaseTime", []any{"0x" + strconv.FormatInt(int64(duration.Seconds()), 16)}, new(float64))
}

// Increase the blockchain current timestamp by the specified amount of time in seconds
func (g *Anvil) IncreaseTime(duration time.Duration) (adjustedTime time.Duration, err error) {
	call := prepareIncreaseTimeCall(duration)
	response := new(interface{})
	err = g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return 0, err
	}
	return time.Duration(int64((*response).(float64))) * time.Second, nil
}

func prepareMineCall(blockCount int, blockLength time.Duration) rpc.BatchElem {
	return prepareCall("anvil_mine", []any{big.NewInt(int64(blockCount)), big.NewInt(int64(blockLength.Seconds()))}, new(string))
}

func (g *Anvil) MineBlocks(blockCount int, blockTime time.Duration) error {
	call := prepareMineCall(blockCount, blockTime)
	response := new(interface{})
	err := g.c.Call(response, call.Method, call.Args...)
	if err != nil {
		return err
	}
	if (*response) != nil {
		return errors.New("unexpected response")
	}
	return nil
}

func (g *Anvil) AvailableAddresses() ([]common.Address, error) {
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
