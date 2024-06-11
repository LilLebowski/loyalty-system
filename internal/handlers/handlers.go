package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/LilLebowski/loyalty-system/internal/auth"
	"github.com/LilLebowski/loyalty-system/internal/storage"
	"github.com/LilLebowski/loyalty-system/internal/utils"
)

type HandlerWithStorage struct {
	storage         storage.Storage
	config          *config.Config
	ordersToProcess chan string
}

func Init(storage storage.Storage, c *config.Config) *HandlerWithStorage {
	return &HandlerWithStorage{storage: storage, config: c, ordersToProcess: make(chan string, 10)}
}

func (strg *HandlerWithStorage) Register(ctx *gin.Context) {
	jsonBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		http.Error(ctx.Writer, "Error while reading body", http.StatusBadRequest)
		return
	}
	var authData storage.Auth
	err = json.Unmarshal(jsonBody, &authData)
	if err != nil {
		http.Error(ctx.Writer, "Error unmarshal body", http.StatusBadRequest)
		return
	}
	hash := sha256.New()
	hash.Write([]byte(authData.Password))
	passwordHash := hex.EncodeToString(hash.Sum(nil))
	userID, err := strg.storage.Register(ctx.Request.Context(), authData.Login, passwordHash)
	if err != nil {
		http.Error(ctx.Writer, "Error login is already in use", http.StatusConflict)
		return
	}
	token, err := auth.BuildJWTString(strg.config, userID)
	if err != nil {
		http.Error(ctx.Writer, "Server error", http.StatusInternalServerError)
		return
	}
	ctx.SetCookie(auth.CookieName, token, 3600, "/", "localhost", false, true)
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(make([]byte, 0))
}

func (strg *HandlerWithStorage) Login(ctx *gin.Context) {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		http.Error(ctx.Writer, "Error while reading body", http.StatusBadRequest)
		return
	}
	var authData storage.Auth
	err = json.Unmarshal(jsonData, &authData)
	if err != nil {
		http.Error(ctx.Writer, "Error unmarshal body", http.StatusBadRequest)
		return
	}
	userData, err := strg.storage.GetUserByLogin(ctx.Request.Context(), authData)
	if err != nil {
		http.Error(ctx.Writer, "Error get user by login", http.StatusUnauthorized)
		return
	}
	hash := sha256.New()
	hash.Write([]byte(authData.Password))
	passwordHash := hex.EncodeToString(hash.Sum(nil))
	if passwordHash == userData.PasswordHash {
		token, err := auth.BuildJWTString(strg.config, userData.ID)
		if err != nil {
			http.Error(ctx.Writer, "Error building token", http.StatusUnauthorized)
			return
		}
		ctx.SetCookie(auth.CookieName, token, 3600, "/", "localhost", false, true)
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(make([]byte, 0))
	} else {
		http.Error(ctx.Writer, "Error login-password pair", http.StatusUnauthorized)
	}
}

func (strg *HandlerWithStorage) AddOrder(ctx *gin.Context) {
	data, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		http.Error(ctx.Writer, "Error while reading body", http.StatusBadRequest)
		return
	}

	_, errCode, _ := utils.ValidateOrder(string(data))
	if errCode != http.StatusOK {
		http.Error(ctx.Writer, "Bad order number", errCode)
		return
	}
	userIDFromContext, _ := ctx.Get(auth.CookieName)
	userID, _ := userIDFromContext.(string)

	err = strg.storage.AddOrderForUser(ctx.Request.Context(), string(data), userID)
	var orderIsExistAnotherUserError *utils.OrderIsExistAnotherUserError
	var orderIsExistThisUserError *utils.OrderIsExistThisUserError
	if errors.As(err, &orderIsExistThisUserError) {
		ctx.Writer.WriteHeader(http.StatusOK)
		return
	}
	if errors.As(err, &orderIsExistAnotherUserError) {
		http.Error(ctx.Writer, "Error add order into db", http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(ctx.Writer, "Error add order into db", http.StatusInternalServerError)
		return
	}

	go func(orderNumber string) {
		strg.ordersToProcess <- orderNumber
	}(string(data))

	ctx.Writer.WriteHeader(http.StatusAccepted)
	ctx.Writer.Write(make([]byte, 0))
}

func (strg *HandlerWithStorage) GetOrders(ctx *gin.Context) {
	userIDFromContext, _ := ctx.Get(auth.CookieName)
	userID, _ := userIDFromContext.(string)
	orders, err := strg.storage.GetOrdersByUser(ctx.Request.Context(), userID)
	if err != nil {
		http.Error(ctx.Writer, "Server error", http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		http.Error(ctx.Writer, "No orders for this user", http.StatusNoContent)
		return
	}
	ordersMarshalled, err := json.Marshal(orders)
	if err != nil {
		http.Error(ctx.Writer, "Error while marshalling", http.StatusInternalServerError)
		return
	}
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(ordersMarshalled)
}

func (strg *HandlerWithStorage) GetBalance(ctx *gin.Context) {
	userIDFromContext, _ := ctx.Get(auth.CookieName)
	userID, _ := userIDFromContext.(string)
	userBalance, err := strg.storage.GetUserBalance(ctx.Request.Context(), userID)
	if err != nil {
		http.Error(ctx.Writer, "Error get user balance", http.StatusInternalServerError)
		return
	}
	userBalanceMarshalled, err := json.Marshal(userBalance)
	if err != nil {
		http.Error(ctx.Writer, "Error while marshalling", http.StatusInternalServerError)
		return
	}
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(userBalanceMarshalled)
}

func (strg *HandlerWithStorage) AddWithdrawal(ctx *gin.Context) {
	userIDFromContext, _ := ctx.Get(auth.CookieName)
	userID, _ := userIDFromContext.(string)
	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		http.Error(ctx.Writer, "Error while getting data", http.StatusInternalServerError)
		return
	}
	var withdrawal storage.Withdrawal
	err = json.Unmarshal(data, &withdrawal)
	if err != nil {
		http.Error(ctx.Writer, "Error while getting data", http.StatusInternalServerError)
		return
	}
	_, errCode, _ := utils.ValidateOrder(withdrawal.Order)
	if errCode != http.StatusOK {
		http.Error(ctx.Writer, "Error bad order number", errCode)
		return
	}
	err = strg.storage.AddWithdrawalForUser(ctx.Request.Context(), userID, withdrawal)
	var lessBonusErrorError *utils.LessBonusError
	if errors.As(err, &lessBonusErrorError) {
		http.Error(ctx.Writer, err.Error(), http.StatusPaymentRequired)
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(make([]byte, 0))
}

func (strg *HandlerWithStorage) GetWithdrawals(ctx *gin.Context) {
	userIDFromContext, _ := ctx.Get(auth.CookieName)
	userID, _ := userIDFromContext.(string)
	withdrawals, err := strg.storage.GetWithdrawalsForUser(ctx.Request.Context(), userID)
	if err != nil {
		http.Error(ctx.Writer, "Server error", http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		http.Error(ctx.Writer, "Error no withdrawals for this user", http.StatusNoContent)
		return
	}
	withdrawalsMarshalled, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(ctx.Writer, "Server error", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(withdrawalsMarshalled))
	ctx.Writer.Header().Set("Content-Type", "application/json")
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(withdrawalsMarshalled)
}
