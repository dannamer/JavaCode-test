package api_test

import (
	"bytes"
	"fmt"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dannamer/JavaCode-test/internal/api"
	"github.com/dannamer/JavaCode-test/internal/api/mock"
	"github.com/dannamer/JavaCode-test/internal/model"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestWalletOperation_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWalletService := mock.NewMockWalletService(ctrl)

	transaction := model.Transaction{
		WalletID:      uuid.New(),
		OperationType: model.Deposit,
		Amount:        decimal.NewFromInt32(100),
	}

	mockWalletService.EXPECT().WalletTransaction(gomock.Any(), transaction).Return(nil)

	handler := api.NewWalletHandler(mockWalletService)

	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handler.WalletOperation(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp model.Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, model.StatusTransactionSuccess, resp.Message)
}

func TestWalletOperation_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок WalletService
	mockWalletService := mock.NewMockWalletService(ctrl)

	// Подготавливаем транзакцию для теста
	transaction := model.Transaction{
		WalletID:      uuid.New(),
		OperationType: model.Withdraw,
		Amount:        decimal.NewFromInt32(1000),
	}

	// Мокаем ошибку "insufficient funds"
	mockWalletService.EXPECT().WalletTransaction(gomock.Any(), transaction).Return(fmt.Errorf("insufficient funds"))

	// Создаем обработчик
	handler := api.NewWalletHandler(mockWalletService)

	// Создаем тестовый HTTP-запрос
	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.WalletOperation(rr, req)

	// Проверка ответа
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	var resp model.Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, model.StatusInsufficientFunds, resp.Message)
}

func TestWalletOperation_WalletNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок WalletService
	mockWalletService := mock.NewMockWalletService(ctrl)

	// Подготавливаем транзакцию для теста
	transaction := model.Transaction{
		WalletID:      uuid.New(),
		OperationType: model.Deposit,
		Amount:        decimal.NewFromInt32(100),
	}

	// Мокаем ошибку "no rows in result set"
	mockWalletService.EXPECT().WalletTransaction(gomock.Any(), transaction).Return(fmt.Errorf("no rows in result set"))

	// Создаем обработчик
	handler := api.NewWalletHandler(mockWalletService)

	// Создаем тестовый HTTP-запрос
	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	handler.WalletOperation(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var resp model.Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(model.StatusWalletNotFound, transaction.WalletID), resp.Message)
}

func TestWalletOperation_WalletErrorDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок WalletService
	mockWalletService := mock.NewMockWalletService(ctrl)

	// Подготавливаем транзакцию для теста
	transaction := model.Transaction{
		WalletID:      uuid.New(),
		OperationType: model.Deposit,
		Amount:        decimal.NewFromInt32(100),
	}

	// Мокаем ошибку "no rows in result set"
	mockWalletService.EXPECT().WalletTransaction(gomock.Any(), transaction).Return(fmt.Errorf("error database"))

	// Создаем обработчик
	handler := api.NewWalletHandler(mockWalletService)

	// Создаем тестовый HTTP-запрос
	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.WalletOperation(rr, req)

	// Проверка ответа
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp model.Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, model.StatusInternalServerError, resp.Message)
}

func TestWalletOperation_WalletErrorValidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок WalletService
	mockWalletService := mock.NewMockWalletService(ctrl)

	// Подготавливаем транзакцию для теста
	transaction := model.Transaction{
		WalletID:      uuid.New(),
		OperationType: model.Deposit,
		Amount:        decimal.NewFromInt32(-100),
	}

	handler := api.NewWalletHandler(mockWalletService)

	// Создаем тестовый HTTP-запрос
	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.WalletOperation(rr, req)

	// Проверка ответа
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp model.Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, model.StatusInvalidRequestData, resp.Message)
}

func TestWalletOperation_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок WalletService
	mockWalletService := mock.NewMockWalletService(ctrl)

	// Создаем обработчик
	handler := api.NewWalletHandler(mockWalletService)

	// Создаем некорректный JSON в теле запроса (например, пропущенная закрывающая скобка)
	invalidJSON := `{"walletID": "123e4567-e89b-12d3-a456-426614174000", "amount": 100,` // Некорректный JSON

	// Создаем тестовый HTTP-запрос с некорректным телом
	req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(invalidJSON))
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.WalletOperation(rr, req)

	// Проверка, что код ответа - 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Проверка, что сообщение в ответе соответствует "Invalid request body"
	var resp model.Response
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, model.StatusInvalidRequestBody, resp.Message)
}

// func TestWalletHandlers_Wallet_InvalidUUIDFormat(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// Создаем мок WalletService
// 	mockWalletService := mock.NewMockWalletService(ctrl)

// 	// Создаем обработчик
// 	handler := api.NewWalletHandler(mockWalletService)

// 	// Некорректный UUID в URL
// 	req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/invalidUUID", nil)
// 	rr := httptest.NewRecorder()

// 	// Вызываем обработчик
// 	handler.Wallet(rr, req)

// 	// Проверка, что код ответа - 400 (Bad Request)
// 	assert.Equal(t, http.StatusBadRequest, rr.Code)

// 	// Проверка, что сообщение в ответе соответствует "Invalid wallet UUID format"
// 	var resp model.Response
// 	err := json.NewDecoder(rr.Body).Decode(&resp)
// 	assert.NoError(t, err)
// 	assert.Equal(t, model.StatusInvalidUUIDFormat, resp.Message)
// }

// func TestWalletHandlers_Wallet_WalletNotFound(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	// Создаем мок WalletService
// 	mockWalletService := mock.NewMockWalletService(ctrl)

// 	// Создаем UUID для теста
// 	walletUUID := uuid.New()

// 	// Мокаем вызов GetWalletBalance, чтобы он вернул ошибку "no rows in result set"
// 	mockWalletService.EXPECT().GetWalletBalance(context.Background(), walletUUID).Return(model.Wallet{}, fmt.Errorf("no rows in result set"))

// 	// Создаем обработчик
// 	handler := api.NewWalletHandler(mockWalletService)

// 	// Запрос с валидным UUID
// 	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallets/%s", walletUUID.String()), nil)
// 	rr := httptest.NewRecorder()

// 	// Вызываем обработчик
// 	handler.Wallet(rr, req)

// 	// Проверка, что код ответа - 404 (Not Found)
// 	assert.Equal(t, http.StatusNotFound, rr.Code)

// 	// Проверка, что сообщение в ответе соответствует "Wallet with UUID <UUID> not found"
// 	var resp model.Response
// 	err := json.NewDecoder(rr.Body).Decode(&resp)
// 	assert.NoError(t, err)
// 	assert.Equal(t, fmt.Sprintf(model.StatusWalletNotFound, walletUUID), resp.Message)
// }
