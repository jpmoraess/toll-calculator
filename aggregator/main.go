package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/jpmoraess/toll-calculator/common"
	"google.golang.org/grpc"
)

func main() {
	httpAddr := flag.String("httpAddr", ":3000", "the listen address of the HTTP server")
	grpcAddr := flag.String("grpcAddr", ":50051", "the listen address of the GRPC server")
	flag.Parse()

	var (
		store   Storer
		service Aggregator
	)

	store = NewMemoryStore()
	service = NewInvoiceAggregator(store)
	service = NewLogMiddleware(service)

	go makeGRPCTransport(*grpcAddr, service)
	makeHTTPTransport(*httpAddr, service)

}

func makeHTTPTransport(listenAddr string, service Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("POST /aggregate", handleAggregate(service))
	http.HandleFunc("GET /invoice", handleGetInvoice(service))
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		fmt.Printf("error server HTTP: %s\n", err)
	}
}

func makeGRPCTransport(listenAddr string, service Aggregator) error {
	fmt.Println("GRPC transport running on port", listenAddr)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer lis.Close()
	server := grpc.NewServer([]grpc.ServerOption{}...)
	common.RegisterAggregatorServer(server, NewGRPCServer(service))
	return server.Serve(lis)
}

func handleAggregate(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance common.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := service.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(service Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuID := r.URL.Query().Get("obuID")
		if obuID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query param obuID is required"})
			return
		}
		id, err := strconv.Atoi(obuID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obuID"})
			return
		}
		invoice, err := service.GetInvoice(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
