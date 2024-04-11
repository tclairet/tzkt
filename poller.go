package main

import (
	"fmt"
	"time"
)

type client interface {
	GetDelegate(block int64) ([]Delegate, error)
	GetBlocksCount() (int64, error)
}

type store interface {
	SaveLastBlock(block int64) error
	LastBlock() (int64, error)
	SaveDelegates(delegates []Delegate) error
	Delegates(year *int) ([]Delegate, error)
}

type poller struct {
	client     client
	frequency  time.Duration
	startBlock int64
	store      store

	stop chan struct{}
}

func newPoller(client client, frequency time.Duration, startBlock int64, store store) *poller {
	return &poller{
		client:     client,
		frequency:  frequency,
		startBlock: startBlock,
		store:      store,
	}
}

func (p *poller) Start() error {
	if p.stop != nil {
		return fmt.Errorf("poller already started")
	}
	p.stop = make(chan struct{})
	go func() {
		ticker := time.NewTicker(p.frequency)
		for {
			select {
			case <-ticker.C:
				if err := p.poll(); err != nil {
					logger.Error("cannot poll", "error", err.Error())
				}
			case <-p.stop:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

func (p *poller) Stop() error {
	if p.stop == nil {
		return fmt.Errorf("cannot stop, poller not started")
	}
	close(p.stop)
	p.stop = nil
	return nil
}

func (p *poller) lastBlock() (int64, error) {
	stored, err := p.store.LastBlock()
	if err != nil {
		return 0, err
	}
	if stored > p.startBlock {
		return stored, nil
	}
	return p.startBlock, nil
}

func (p *poller) poll() error {
	lastBlockStored, err := p.lastBlock()
	if err != nil {
		return err
	}

	lastBlock, err := p.client.GetBlocksCount()
	if err != nil {
		return err
	}

	if lastBlockStored >= lastBlock {
		return nil
	}

	if lastBlockStored == 0 {
		lastBlockStored = lastBlock - 1
	}

	logger.Info("polling", "start", lastBlockStored+1, "end", lastBlock)

	var delegates []Delegate
	for i := lastBlockStored + 1; i <= lastBlock; i++ {
		d, err := p.client.GetDelegate(i)
		if err != nil {
			return err
		}
		delegates = append(delegates, d...)
	}
	if err := p.store.SaveDelegates(delegates); err != nil {
		return err
	}
	if err := p.store.SaveLastBlock(lastBlock); err != nil {
		return err
	}

	if len(delegates) != 0 {
		logger.Info("new delegates", "count", len(delegates))
	}

	return err
}
