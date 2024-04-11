package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

const delegationsRoute = "/xtz/delegations"

type delegatesStore interface {
	SaveDelegates(delegates []Delegate) error
	Delegates(year *int) ([]Delegate, error)
}

type API struct {
	store delegatesStore
}

func (api API) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get(delegationsRoute, api.delegationsHandler)
	return r
}

func (api API) delegationsHandler(w http.ResponseWriter, r *http.Request) {
	var year *int
	yearStr := r.URL.Query().Get("year")
	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, err)
			return
		}
		year = &y
	}

	delegations, err := api.delegations(year)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	RespondWithJSON(w, http.StatusOK, DelegationsResponse{delegations})
}

type DelegationsResponse struct {
	Data []Delegate `json:"data"`
}

func (api API) delegations(year *int) ([]Delegate, error) {
	delegates, err := api.store.Delegates(year)
	if err != nil {
		return nil, err
	}
	return delegates, nil
}

func RespondWithError(w http.ResponseWriter, code int, msg interface{}) {
	var message string
	switch m := msg.(type) {
	case error:
		message = m.Error()
	case string:
		message = m
	}
	RespondWithJSON(w, code, JSONError{Error: message})
}

type JSONError struct {
	Error string `json:"error,omitempty"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
