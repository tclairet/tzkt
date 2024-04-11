package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAPI_delegations(t *testing.T) {
	tests := []struct {
		name    string
		year    string
		toStore []Delegate
		want    DelegationsResponse
	}{
		{
			name:    "retrieve all",
			toStore: []Delegate{{Timestamp: "2020-01-01T01:01:01Z"}, {Timestamp: "2024-01-01T01:01:01Z"}},
			want:    DelegationsResponse{[]Delegate{{Timestamp: "2020-01-01T01:01:01Z"}, {Timestamp: "2024-01-01T01:01:01Z"}}},
		},
		{
			name:    "filter by year",
			year:    "2020",
			toStore: []Delegate{{Timestamp: "2020-01-01T01:01:01Z"}, {Timestamp: "2024-01-01T01:01:01Z"}},
			want:    DelegationsResponse{[]Delegate{{Timestamp: "2020-01-01T01:01:01Z"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newMem()
			_ = store.SaveDelegates(tt.toStore)
			api := API{
				store: store,
			}
			server := httptest.NewServer(api.Routes())
			t.Cleanup(func() { server.Close() })

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s?year=%s", server.URL, delegationsRoute, tt.year), nil)
			if err != nil {
				t.Fatal(err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			var got DelegationsResponse
			if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("delegations() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func address(in int) *int {
	return &in
}
