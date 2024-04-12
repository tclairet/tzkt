package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type mem map[string][]byte

func newMem() mem {
	return map[string][]byte{}
}

func (m mem) SaveLastBlock(block int64) error {
	m["block"] = []byte(fmt.Sprintf("%d", block))
	return nil
}

func (m mem) LastBlock() (int64, error) {
	if len(m["block"]) == 0 {
		return 0, nil
	}
	block, err := strconv.Atoi(string(m["block"]))
	if err != nil {
		return 0, err
	}
	return int64(block), nil
}

func (m mem) SaveDelegates(delegates []Delegate) error {
	stored, err := m.Delegates(nil)
	if err != nil {
		return err
	}
	stored = append(stored, delegates...)
	b, err := json.Marshal(stored)
	if err != nil {
		return err
	}
	m["delegates"] = b
	return nil
}

func (m mem) Delegates(year *int) ([]Delegate, error) {
	if len(m["delegates"]) == 0 {
		return []Delegate{}, nil
	}
	var delegates []Delegate
	if err := json.Unmarshal(m["delegates"], &delegates); err != nil {
		return nil, err
	}
	if year != nil {
		delegates = slices.DeleteFunc(delegates, func(delegate Delegate) bool {
			return *year != delegate.Timestamp.Year()
		})
	}
	slices.SortFunc(delegates, func(a, b Delegate) int {
		return a.Timestamp.Compare(b.Timestamp)
	})
	return delegates, nil
}

type Badger struct {
	db *badger.DB
}

func NewBadger(path string) (*Badger, error) {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	return &Badger{
		db: db,
	}, nil
}

func (s *Badger) SaveLastBlock(block int64) error {
	return s.persist("lastBlock", []byte(fmt.Sprintf("%d", block)))
}

func (s *Badger) LastBlock() (int64, error) {
	b, err := s.get("lastBlock")
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		return 0, nil
	}
	lastBlock, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, err
	}
	return int64(lastBlock), nil
}

func (s *Badger) SaveDelegates(delegates []Delegate) error {
	for _, delegate := range delegates {
		b, err := json.Marshal(delegate)
		if err != nil {
			return err
		}
		if err := s.persist(key(delegate), b); err != nil {
			return err
		}
	}
	return nil
}

func (s *Badger) Delegates(year *int) ([]Delegate, error) {
	yearStr := ""
	if year != nil {
		yearStr = fmt.Sprintf("%d", *year)
	}
	res, err := s.read("delegates" + yearStr)
	if err != nil {
		return nil, err
	}

	var delegates []Delegate
	for _, b := range res {
		var d Delegate
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, err
		}
		delegates = append(delegates, d)
	}

	slices.SortFunc(delegates, func(a, b Delegate) int {
		return a.Timestamp.Compare(b.Timestamp)
	})
	return delegates, nil
}

func (s *Badger) persist(key string, value []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), value)
		return err
	})
	return err
}

func (s *Badger) get(key string) ([]byte, error) {
	var valCopy []byte
	if err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})

		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return valCopy, nil
}

func (s *Badger) read(prefix string) (map[string][]byte, error) {
	result := make(map[string][]byte)
	if err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		p := []byte(prefix)
		for it.Seek(p); it.ValidForPrefix(p); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				result[strings.TrimPrefix(string(k), prefix)] = v
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func key(d Delegate) string {
	return "delegates" + d.Timestamp.Format(time.RFC3339)
}
