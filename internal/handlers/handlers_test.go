package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/LilLebowski/loyalty-system/internal/auth"
	"github.com/LilLebowski/loyalty-system/internal/mock_storage"
	"github.com/LilLebowski/loyalty-system/internal/storage"
	"github.com/LilLebowski/loyalty-system/internal/utils"
)

type wantResponse struct {
	code            int
	headerContent   string
	responseContent string
}

func SetUpRouter(s storage.Storage, cfg *config.Config) *gin.Engine {
	handlerWithStorage := Init(s, cfg)
	router := gin.Default()
	authRouter := router.Group("/api/user/auth")
	{
		authRouter.POST("/register", handlerWithStorage.Register)
		authRouter.POST("/login", handlerWithStorage.Login)
	}
	userRouter := router.Group("/api/user", auth.Authorization(cfg))
	{
		userRouter.POST("/register", handlerWithStorage.Register)
		userRouter.POST("/login", handlerWithStorage.Login)
		userRouter.POST("/orders", handlerWithStorage.AddOrder)
		userRouter.GET("/orders", handlerWithStorage.GetOrders)
		userRouter.GET("/balance", handlerWithStorage.GetBalance)
		userRouter.POST("/balance/withdraw", handlerWithStorage.AddWithdrawal)
		userRouter.GET("/withdrawals", handlerWithStorage.GetWithdrawals)
	}
	return router
}

func TestRegisterHandler(t *testing.T) {
	tt := []struct {
		name            string
		want            wantResponse
		registerData    storage.Auth
		mockResponseID  string
		mockResponseErr error
	}{
		{
			"User register 200",
			wantResponse{
				http.StatusOK,
				"",
				``,
			},
			storage.Auth{Login: "test", Password: "test"},
			"bdf8817b-3225-4a46-9358-aa091b3cb478",
			nil,
		},
		{
			"User register 409",
			wantResponse{
				http.StatusConflict,
				"text/plain; charset=utf-8",
				"Error login is already in use\n",
			},
			storage.Auth{Login: "test", Password: "test"},
			"",
			fmt.Errorf("user with login test exist"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Init()
			marshalledData, _ := json.Marshal(tc.registerData)
			request := httptest.NewRequest(http.MethodPost, "/api/user/register", bytes.NewBuffer(marshalledData))
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			strg := mock_storage.NewMockStorage(ctrl)
			hash := sha256.New()
			hash.Write([]byte(tc.registerData.Password))
			passwordHash := hex.EncodeToString(hash.Sum(nil))
			strg.
				EXPECT().
				Register(ctx, tc.registerData.Login, passwordHash).
				Return(tc.mockResponseID, tc.mockResponseErr).
				AnyTimes()

			routerInstance := SetUpRouter(strg, cfg)
			routerInstance.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				cookies := result.Cookies()
				for _, cookie := range cookies {
					if cookie.Name == auth.CookieName {
						token, _ := auth.BuildJWTString(cfg, tc.mockResponseID)
						assert.Equal(t, token, cookie.Value)
						break
					}
					assert.Fail(t, "get no cookies for UserID")
				}
			}
			responseBody, err := io.ReadAll(result.Body)
			assert.Nil(t, err)
			assert.Equal(t, tc.want.headerContent, result.Header.Get("Content-Type"))
			assert.Equal(t, tc.want.responseContent, string(responseBody))
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tt := []struct {
		name            string
		want            wantResponse
		loginData       storage.Auth
		mockResponse    storage.User
		mockResponseErr error
	}{
		{
			"User login 200",
			wantResponse{
				http.StatusOK,
				"",
				``,
			},
			storage.Auth{Login: "test", Password: "test"},
			storage.User{
				Login:        "test",
				PasswordHash: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
				ID:           "bdf8817b-3225-4a46-9358-aa091b3cb478",
			},
			nil,
		},
		{
			"User login 401",
			wantResponse{
				http.StatusUnauthorized,
				"text/plain; charset=utf-8",
				"Error login-password pair\n",
			},
			storage.Auth{Login: "test", Password: "test"},
			storage.User{Login: "test", PasswordHash: "test", ID: "bdf8817b-3225-4a46-9358-aa091b3cb478"},
			nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Init()
			marshalledData, _ := json.Marshal(tc.loginData)
			request := httptest.NewRequest(http.MethodPost, "/api/user/login", bytes.NewBuffer(marshalledData))
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			strg := mock_storage.NewMockStorage(ctrl)
			strg.
				EXPECT().
				GetUserByLogin(ctx, tc.loginData).
				Return(tc.mockResponse, tc.mockResponseErr).
				AnyTimes()

			routerInstance := SetUpRouter(strg, cfg)
			routerInstance.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
			if result.StatusCode == http.StatusOK {
				cookies := result.Cookies()
				for _, cookie := range cookies {
					if cookie.Name == auth.CookieName {
						token, _ := auth.BuildJWTString(cfg, tc.mockResponse.ID)
						assert.Equal(t, token, cookie.Value)
						break
					}
					assert.Fail(t, "get no cookies for UserID")
				}
			}
			responseBody, err := io.ReadAll(result.Body)
			assert.Nil(t, err)
			assert.Equal(t, tc.want.headerContent, result.Header.Get("Content-Type"))
			assert.Equal(t, tc.want.responseContent, string(responseBody))
		})
	}
}

func TestAddOrderHandler(t *testing.T) {
	tt := []struct {
		name            string
		want            wantResponse
		externalOrderID string
		userID          string
		mockResponseErr error
	}{
		{
			"Add order 200",
			wantResponse{
				http.StatusOK,
				"",
				"",
			},
			"12345678903",
			"bdf8817b-3225-4a46-9358-aa091b3cb478",
			utils.NewOrderIsExistThisUserError("this order is exist the user", nil),
		},
		{
			"Add order 202",
			wantResponse{
				http.StatusAccepted,
				"",
				"",
			},
			"12345678903",
			"bdf8817b-3225-4a46-9358-aa091b3cb478",
			nil,
		},
		{
			"Add order 409",
			wantResponse{
				http.StatusConflict,
				"",
				"",
			},
			"12345678903",
			"bdf8817b-3225-4a46-9358-aa091b3cb478",
			utils.NewOrderIsExistAnotherUserError("this order is exist another user", nil),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Init()
			token, _ := auth.BuildJWTString(cfg, tc.userID)
			param := strings.NewReader(tc.externalOrderID)
			request := httptest.NewRequest(http.MethodPost, "/api/user/orders", param)
			newCookie := http.Cookie{Name: auth.CookieName, Value: token}
			request.Header.Add("Cookie", newCookie.String())
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			strg := mock_storage.NewMockStorage(ctrl)
			strg.
				EXPECT().
				AddOrderForUser(ctx, tc.externalOrderID, tc.userID).
				Return(tc.mockResponseErr).
				AnyTimes()

			routerInstance := SetUpRouter(strg, cfg)
			routerInstance.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
		})
	}
}

func TestGetOrdersHandler(t *testing.T) {
	tt := []struct {
		name            string
		want            wantResponse
		userID          string
		mockResponse    []storage.Order
		mockResponseErr error
	}{
		{
			"Get orders 200",
			wantResponse{
				http.StatusOK,
				"",
				"",
			},
			"12345678903",
			[]storage.Order{
				{
					Number:     "test",
					Status:     "NEW",
					Accrual:    200,
					UploadedAt: time.Time{},
				},
			},
			nil,
		},
		{
			"Get orders 204",
			wantResponse{
				http.StatusNoContent,
				"",
				"",
			},
			"12345678903",
			[]storage.Order{},
			nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Init()
			token, _ := auth.BuildJWTString(cfg, tc.userID)
			jsonBody, _ := json.Marshal(tc.mockResponse)
			param := strings.NewReader(string(jsonBody))
			request := httptest.NewRequest(http.MethodGet, "/api/user/orders", param)
			newCookie := http.Cookie{Name: auth.CookieName, Value: token}
			request.Header.Add("Cookie", newCookie.String())
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			strg := mock_storage.NewMockStorage(ctrl)
			strg.
				EXPECT().
				GetOrdersByUser(ctx, tc.userID).
				Return(tc.mockResponse, tc.mockResponseErr).
				AnyTimes()

			routerInstance := SetUpRouter(strg, cfg)
			routerInstance.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
		})
	}
}

func TestGetUserBalanceHandler(t *testing.T) {
	tt := []struct {
		name            string
		want            wantResponse
		userID          string
		mockResponse    storage.UserBalance
		mockResponseErr error
	}{
		{
			"Get orders 200",
			wantResponse{
				http.StatusOK,
				"",
				"",
			},
			"12345678903",
			storage.UserBalance{
				Current:   100,
				Withdrawn: 100,
			},
			nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Init()
			token, _ := auth.BuildJWTString(cfg, tc.userID)
			request := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			newCookie := http.Cookie{Name: auth.CookieName, Value: token}
			request.Header.Add("Cookie", newCookie.String())
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			strg := mock_storage.NewMockStorage(ctrl)
			strg.
				EXPECT().
				GetUserBalance(ctx, tc.userID).
				Return(tc.mockResponse, tc.mockResponseErr).
				AnyTimes()

			routerInstance := SetUpRouter(strg, cfg)
			routerInstance.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
		})
	}
}

func TestAddWithdrawalHandler(t *testing.T) {
	tt := []struct {
		name            string
		want            wantResponse
		userID          string
		withdrawal      storage.Withdrawal
		mockResponseErr error
	}{
		{
			"Add withdrawal 200",
			wantResponse{
				http.StatusOK,
				"",
				"",
			},
			"12345678903",
			storage.Withdrawal{
				ExternalOrderID: "12345678903",
				Sum:             100,
			},
			nil,
		},
		{
			"Add withdrawal 402",
			wantResponse{
				http.StatusPaymentRequired,
				"",
				"",
			},
			"12345678903",
			storage.Withdrawal{
				ExternalOrderID: "12345678903",
				Sum:             100,
			},
			utils.NewLessBonusErrorError("Got less bonus points than expected", nil),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Init()
			token, _ := auth.BuildJWTString(cfg, tc.userID)
			jsonBody, _ := json.Marshal(tc.withdrawal)
			param := strings.NewReader(string(jsonBody))
			request := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", param)
			newCookie := http.Cookie{Name: auth.CookieName, Value: token}
			request.Header.Add("Cookie", newCookie.String())
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			strg := mock_storage.NewMockStorage(ctrl)
			strg.
				EXPECT().
				AddWithdrawalForUser(ctx, tc.userID, tc.withdrawal).
				Return(tc.mockResponseErr).
				AnyTimes()

			routerInstance := SetUpRouter(strg, cfg)
			routerInstance.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
		})
	}
}

func TestGetWithdrawalsHandler(t *testing.T) {
	tt := []struct {
		name            string
		want            wantResponse
		userID          string
		mockResponse    []storage.Withdrawal
		mockResponseErr error
	}{
		{
			"Get withdrawals 200",
			wantResponse{
				http.StatusOK,
				"",
				"",
			},
			"12345678903",
			[]storage.Withdrawal{
				{
					ExternalOrderID: "12345678903",
					Sum:             200,
					ProcessedAt:     time.Time{},
				},
			},
			nil,
		},
		{
			"Get withdrawals 204",
			wantResponse{
				http.StatusNoContent,
				"",
				"",
			},
			"12345678903",
			[]storage.Withdrawal{},
			nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			cfg := config.Init()
			token, _ := auth.BuildJWTString(cfg, tc.userID)
			request := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
			newCookie := http.Cookie{Name: auth.CookieName, Value: token}
			request.Header.Add("Cookie", newCookie.String())
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(t)
			strg := mock_storage.NewMockStorage(ctrl)
			strg.
				EXPECT().
				GetWithdrawalsForUser(ctx, tc.userID).
				Return(tc.mockResponse, tc.mockResponseErr).
				AnyTimes()

			routerInstance := SetUpRouter(strg, cfg)
			routerInstance.ServeHTTP(w, request)

			result := w.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.want.code, result.StatusCode)
		})
	}
}
