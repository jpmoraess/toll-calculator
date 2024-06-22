package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/jpmoraess/toll-calculator/common"
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "the listen address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	service := NewInvoiceAggregator(store)

	makeHTTPTransport(*listenAddr, service)
}

func makeHTTPTransport(listenAddr string, service Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(service))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance common.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
