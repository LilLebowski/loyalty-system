package router

import (
	"github.com/gin-gonic/gin"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/LilLebowski/loyalty-system/internal/handlers"
	"github.com/LilLebowski/loyalty-system/internal/storage"
)

func Init(s *storage.Storage, cfg *config.Config) *gin.Engine {
	handlerWithStorage := handlers.Init(s, cfg)
	router := gin.Default()

	router.POST("/api/user/register", handlerWithStorage.Register)
	router.POST("/api/user/login", handlerWithStorage.Login)
	router.POST("/api/user/orders", handlerWithStorage.AddOrder)
	router.GET("/api/user/orders", handlerWithStorage.GetOrders)
	router.GET("/api/user/balance", handlerWithStorage.GetBalance)
	router.POST("/api/user/balance/withdraw", handlerWithStorage.AddWithdrawal)
	router.GET("/api/user/withdrawals", handlerWithStorage.GetWithdrawals)

	router.HandleMethodNotAllowed = true

	return router
}
