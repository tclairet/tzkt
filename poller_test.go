package main

import (
	"testing"
	"time"
)

func Test_pollerStartStop(t *testing.T) {
	poller := newPoller(fakeClient{}, time.Second, 0, newMem())
	if err := poller.Start(); err != nil {
		t.Fatal(err)
	}
	if err := poller.Stop(); err != nil {
		t.Fatal(err)
	}
}

func Test_poller(t *testing.T) {
	tests := []struct {
		name        string
		delegates   map[int64][]Delegate
		blocksCount int64
	}{
		{
			name: "",
			delegates: map[int64][]Delegate{
				1: {{Delegator: "receiver"}},
			},
			blocksCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newMem()
			p := poller{
				client: fakeClient{
					delegates:   tt.delegates,
					blocksCount: tt.blocksCount,
				},
				startBlock: 0,
				store:      s,
			}
			if err := p.poll(); err != nil {
				t.Fatal(err)
			}
			delegates, err := s.Delegates(nil)
			if err != nil {
				t.Fatal(err)
			}

			for _, expectedDelegates := range tt.delegates {
				for _, expected := range expectedDelegates {
					found := false
					for _, d := range delegates {
						if d.Delegator == expected.Delegator {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("missing %s", expected.Delegator)
					}
				}
			}
		})
	}
}

type fakeClient struct {
	delegates   map[int64][]Delegate
	blocksCount int64
}

func (f fakeClient) GetDelegate(block int64) ([]Delegate, error) {
	return f.delegates[block], nil
}

func (f fakeClient) GetBlocksCount() (int64, error) {
	return f.blocksCount, nil
}
