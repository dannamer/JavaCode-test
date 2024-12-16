package api

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func (h *WalletHandlers) RunServer() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/wallets/{WALLET_UUID}", h.Wallet).Methods("GET")
	r.HandleFunc("/api/v1/wallet", h.WalletOperation).Methods("POST")

	log.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
