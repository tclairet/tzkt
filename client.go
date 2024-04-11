package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type tzkt struct{}

func (sdk *tzkt) GetDelegate(block int64) ([]Delegate, error) {
	urlLevel := fmt.Sprintf("https://api.tzkt.io/v1/operations/delegations?level=%d", block)
	req, err := http.NewRequest(http.MethodGet, urlLevel, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var resp []GetDelegateResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	delegates := make([]Delegate, len(resp))
	for i := 0; i < len(delegates); i++ {
		delegates[i] = Delegate{
			Amount:    fmt.Sprintf("%d", resp[i].Amount),
			Delegator: resp[i].Sender.Address,
			Block:     fmt.Sprintf("%d", resp[i].Level),
			Timestamp: resp[i].Timestamp.Format(time.RFC3339),
		}
	}

	return delegates, nil
}

func (sdk *tzkt) GetBlocksCount() (int64, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.tzkt.io/v1/blocks/count", nil)
	if err != nil {
		return 0, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, err
	}
	return int64(num), nil
}

type Delegate struct {
	Amount    string `json:"amount"`
	Delegator string `json:"delegator"`
	Block     string `json:"block"`
	Timestamp string `json:"timestamp"`
}

type GetDelegateResponse struct {
	Type      string    `json:"type"`
	Id        int       `json:"id"`
	Level     int       `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Block     string    `json:"block"`
	Hash      string    `json:"hash"`
	Counter   int       `json:"counter"`
	Initiator struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"initiator"`
	Sender struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"sender"`
	SenderCodeHash       int `json:"senderCodeHash"`
	Nonce                int `json:"nonce"`
	GasLimit             int `json:"gasLimit"`
	GasUsed              int `json:"gasUsed"`
	StorageLimit         int `json:"storageLimit"`
	BakerFee             int `json:"bakerFee"`
	Amount               int `json:"amount"`
	UnstakedPseudotokens int `json:"unstakedPseudotokens"`
	UnstakedBalance      int `json:"unstakedBalance"`
	UnstakedRewards      int `json:"unstakedRewards"`
	PrevDelegate         struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"prevDelegate"`
	NewDelegate struct {
		Alias   string `json:"alias"`
		Address string `json:"address"`
	} `json:"newDelegate"`
	Status string `json:"status"`
	Errors []struct {
		Type string `json:"type"`
	} `json:"errors"`
	Quote struct {
		Btc int `json:"btc"`
		Eur int `json:"eur"`
		Usd int `json:"usd"`
		Cny int `json:"cny"`
		Jpy int `json:"jpy"`
		Krw int `json:"krw"`
		Eth int `json:"eth"`
		Gbp int `json:"gbp"`
	} `json:"quote"`
}
