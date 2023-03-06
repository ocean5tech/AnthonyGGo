package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/ocean5tech/pricefetcher/types"
)

type APIFunc func(context.Context, http.ResponseWriter, *http.Request) error

type JSONAPIServer struct {
	listenAddr string
	svc        PriceFetcher
}

func NewJSONAPIServer(listenAddr string, svc PriceFetcher) *JSONAPIServer {
	return &JSONAPIServer{
		listenAddr: listenAddr,
		svc:        svc,
	}
}

func (s *JSONAPIServer) Run() {
	http.HandleFunc("/", MakeHTTPHandlerFunc(s.HandleFetchPrice))
	http.ListenAndServe(s.listenAddr, nil)
}

func MakeHTTPHandlerFunc(apiFn APIFunc) http.HandlerFunc {
	ctx := context.Background()

	return func(w http.ResponseWriter, r *http.Request) {
		ctx = context.WithValue(ctx, "requestID", rand.Intn(10000000))
		if err := apiFn(ctx, w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err})
		}
	}
}

func (s *JSONAPIServer) HandleFetchPrice(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ticker := r.URL.Query().Get("ticker")
	price, err := s.svc.FetchPrice(ctx, ticker)
	if err != nil {
		return err
	}

	priceResp := types.PriceResponse{
		Price:  price,
		Ticker: ticker,
	}
	return writeJSON(w, http.StatusOK, &priceResp)
}

func writeJSON(w http.ResponseWriter, s int, v any) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(v)
}
