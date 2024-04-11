package main

import (
	"reflect"
	"testing"
)

func Test_Store(t *testing.T) {
	testsStore := []struct {
		name string
		new  func(t *testing.T) store
	}{
		{
			name: "mem",
			new: func(t *testing.T) store {
				return newMem()
			},
		},
		{
			name: "badger",
			new: func(t *testing.T) store {
				s, err := NewBadger(t.TempDir())
				if err != nil {
					t.Fatal(err)
				}
				return s
			},
		},
	}
	for _, testStore := range testsStore {
		t.Run(testStore.name, func(t *testing.T) {
			tests := []struct {
				name    string
				year    *int
				toStore []Delegate
				want    []Delegate
			}{
				{
					name:    "already ordered",
					toStore: []Delegate{{Timestamp: "2020-01-01T01:01:01Z"}, {Timestamp: "2024-01-01T01:01:01Z"}},
					want:    []Delegate{{Timestamp: "2020-01-01T01:01:01Z"}, {Timestamp: "2024-01-01T01:01:01Z"}},
				},
				{
					name:    "sort",
					toStore: []Delegate{{Timestamp: "2024-01-01T00:00:00Z"}, {Timestamp: "2020-01-01T00:00:00Z"}},
					want:    []Delegate{{Timestamp: "2020-01-01T00:00:00Z"}, {Timestamp: "2024-01-01T00:00:00Z"}},
				},
				{
					name:    "filter by year",
					year:    address(2020),
					toStore: []Delegate{{Timestamp: "2020-01-01T01:01:01Z"}, {Timestamp: "2024-01-01T01:01:01Z"}},
					want:    []Delegate{{Timestamp: "2020-01-01T01:01:01Z"}},
				},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					store := testStore.new(t)
					if err := store.SaveDelegates(tt.toStore); err != nil {
						t.Fatal(err)
					}
					got, err := store.Delegates(tt.year)
					if err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("got = %v, want %v", got, tt.want)
					}

					expected := int64(42)
					if err := store.SaveLastBlock(expected); err != nil {
						t.Error(err)
					}
					lastBlock, err := store.LastBlock()
					if err != nil {
						t.Error(err)
					}
					if got, want := lastBlock, expected; got != want {
						t.Errorf("got = %v, want %v", got, want)
					}
				})
			}
		})
	}
}
