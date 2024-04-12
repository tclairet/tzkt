package main

import (
	"fmt"
	"testing"
	"time"
)

func Test_tzkt_GetDelegate(t *testing.T) {
	tests := []struct {
		name          string
		block         int64
		wantLen       int
		wantSender    string
		wantTimestamp string
		wantAmount    string
	}{
		{
			// https://tzstats.com/5409254#delegations
			name:          "block 5409254",
			block:         5409254,
			wantLen:       1,
			wantSender:    "tz1M2AoHajDRaczmD6bH1xe1KMTfAvGNwKAc",
			wantTimestamp: "2024-04-10T10:57:29Z",
			wantAmount:    "7245529",
		},
	}

	sdk := &tzkt{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delegates, err := sdk.GetDelegate(tt.block)
			if err != nil {
				t.Fatal(err)
			}
			if got, want := len(delegates), tt.wantLen; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := delegates[0].Delegator, tt.wantSender; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := delegates[0].Timestamp.Format(time.RFC3339), tt.wantTimestamp; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := delegates[0].Amount, tt.wantAmount; got != want {
				t.Errorf("got %v, want %v", got, want)
			}
			if got, want := delegates[0].Block, fmt.Sprintf("%d", tt.block); got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func Test_tzkt_GetBlocksCount(t *testing.T) {
	sdk := &tzkt{}

	got, err := sdk.GetBlocksCount()
	if err != nil {
		t.Fatal(err)
	}
	if got == 0 {
		t.Error("block count is 0")
	}
}
