package learning

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Ganache struct {
	address *url.URL
	debug   bool
}

func NewGanache(debug bool) *Ganache {
	url, err := url.Parse("http://localhost:8545")
	if err != nil {
		panic(err)
	}
	return &Ganache{url, debug}
}

// func (g *Ganache) NewGanacheWithBlocksToTime(blocksCount int, startBlockTime time.Time, blockInfo func(blockNo big.Int) int) (*ethclient.Client, *Ganache, func()) {
// 	ganache := NewGanache()
// 	lastBlockTime = time.Now()
// 	ganache.IncreaseTime(startBlockTime.)
// 	for i := 0; i < blocksCount; i++ {
// 		blocDuration = blockInfo
// 		_, err := ganache.IncreaseTime(blockchainDuration)
// 		require.NoError(t, err)
// 		_, err = ganache.MineTo(1)
// 		require.NoError(t, err)
// 	}
// 	return client, ganache, tearDown
// }

type ganachePayload struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

func (g *Ganache) sendRequest(payload ganachePayload) (interface{}, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", g.address.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	if g.debug {
		fmt.Printf("@dd sent: %s; rec: %#v\n", payloadBytes, response["result"])
	}
	return response["result"], nil
}

func generatePayload(method string, params []any) ganachePayload {
	return ganachePayload{
		ID:      1337,
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

func (g *Ganache) TakeSnapshot() (snapshotId int, err error) {
	payload := generatePayload("evm_snapshot", []any{})
	response, err := g.sendRequest(payload)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseInt(strings.Replace(strings.ToLower(response.(string)), "0x", "", -1), 16, 16)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func (g *Ganache) RevertSnapshot(snapshotId int) error {
	payload := generatePayload("evm_revert", []any{snapshotId})
	response, err := g.sendRequest(payload)
	if err != nil {
		return err
	}
	if !response.(bool) {
		return errors.New("failed to revert snapshot")
	}
	return nil
}

// SetTime set the blockchain timestamp to a specific time
// Returns duration between the given timestamp and the current time.
// Beware - use this method cautiously as it allows you to move backwards in time, which may cause new blocks to appear
// to be mined before older blocks, thereby invalidating the blockchain state.
func (g *Ganache) SetTime(newTime time.Time) (offset time.Duration, err error) {
	payload := generatePayload("evm_setTime", []any{newTime.UnixMilli()})
	response, err := g.sendRequest(payload)
	if err != nil {
		return 0, err
	}
	return time.Duration(int64(response.(float64))) * time.Second, nil
}

func (g *Ganache) IncreaseTime(duration time.Duration) (adjustedTime time.Duration, err error) {
	payload := generatePayload("evm_increaseTime", []any{"0x" + strconv.FormatInt(int64(duration.Seconds()), 16)})
	response, err := g.sendRequest(payload)
	if err != nil {
		return 0, err
	}
	return time.Duration(int64(response.(float64))) * time.Second, nil
}

func (g *Ganache) MineBlocks(blockCount int) (int64, error) {
	payload := generatePayload("evm_mine", []any{map[string]any{"blocks": blockCount}})
	response, err := g.sendRequest(payload)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseInt(strings.Replace(strings.ToLower(response.(string)), "0x", "", -1), 16, 16)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (g *Ganache) AvailableAddresses() ([]common.Address, error) {
	payload := generatePayload("eth_accounts", []any{})
	response, err := g.sendRequest(payload)
	if err != nil {
		return nil, err
	}
	addresses := make([]common.Address, 0, len(response.([]interface{})))
	for _, account := range response.([]interface{}) {
		addresses = append(addresses, common.HexToAddress(account.(string)))
	}
	return addresses, nil
}
