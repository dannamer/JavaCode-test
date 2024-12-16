package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dannamer/JavaCode-test/internal/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type WalletService interface {
	WalletTransaction(ctx context.Context, transaction model.Transaction) error
	GetWalletBalance(ctx context.Context, UUID uuid.UUID) (model.Wallet, error)
}

type WalletHandlers struct {
	WalletService
}

func NewWalletHandler(WalletService WalletService) WalletHandlers {
	return WalletHandlers{WalletService: WalletService}
}

func sendResponse(w http.ResponseWriter, r *http.Request, resp model.Response) {
	log.Printf("Handled request to %s, Status: %d, Message: %s", r.URL.Path, resp.Status, resp.Message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	json.NewEncoder(w).Encode(resp)
}

func (h *WalletHandlers) WalletOperation(w http.ResponseWriter, r *http.Request) {
	var response model.Transaction

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		sendResponse(w, r, model.Response{
			Status:  http.StatusBadRequest,
			Message: model.StatusInvalidRequestBody,
		})
		return
	}

	if !response.Validate() {
		sendResponse(w, r, model.Response{
			Status:  http.StatusBadRequest,
			Message: model.StatusInvalidRequestData,
		})
		return
	}

	err := h.WalletTransaction(context.Background(), response)
	if err != nil {
		if err.Error() == "insufficient funds" {
			sendResponse(w, r, model.Response{
				Status:  http.StatusUnprocessableEntity,
				Message: model.StatusInsufficientFunds,
			})
			return
		}
		if err.Error() == "no rows in result set" {
			sendResponse(w, r, model.Response{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf(model.StatusWalletNotFound, response.WalletID),
			})
			return
		}
		sendResponse(w, r, model.Response{
			Status:  http.StatusInternalServerError,
			Message: model.StatusInternalServerError,
		})
		return
	}

	sendResponse(w, r, model.Response{
		Status:  http.StatusOK,
		Message: model.StatusTransactionSuccess,
	})
}

func (h *WalletHandlers) Wallet(w http.ResponseWriter, r *http.Request) {
	walletUUIDStr := mux.Vars(r)["WALLET_UUID"]

	walletUUID, err := uuid.Parse(walletUUIDStr)
	if err != nil {
		sendResponse(w, r, model.Response{
			Status:  http.StatusBadRequest,
			Message: model.StatusInvalidUUIDFormat,
		})
		return
	}

	wallet, err := h.GetWalletBalance(context.Background(), walletUUID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			sendResponse(w, r, model.Response{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf(model.StatusWalletNotFound, walletUUID),
			})
			return
		}
		sendResponse(w, r, model.Response{
			Status:  http.StatusInternalServerError,
			Message: model.StatusInternalServerError,
		})
		return
	}
	sendResponse(w, r, model.Response{
		Status:  http.StatusOK,
		Message: model.StatusWalletBalanceSuccess,
		Data:    wallet,
	})
}
