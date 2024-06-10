package router

import (
	"github.com/gin-gonic/gin"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/LilLebowski/loyalty-system/internal/auth"
	"github.com/LilLebowski/loyalty-system/internal/handlers"
	"github.com/LilLebowski/loyalty-system/internal/storage"
)

func Init(s storage.Storage, cfg *config.Config) *gin.Engine {
	handlerWithStorage := handlers.Init(s, cfg)
	router := gin.Default()

	authRouter := router.Group("/api/user/auth")
	{
		authRouter.POST("/register", handlerWithStorage.Register)
		authRouter.POST("/login", handlerWithStorage.Login)
	}

	userRouter := router.Group("/api/user", auth.Authorization(cfg))
	{
		userRouter.POST("/orders", handlerWithStorage.AddOrder)
		userRouter.GET("/orders", handlerWithStorage.GetOrders)
		userRouter.GET("/balance", handlerWithStorage.GetBalance)
		userRouter.POST("/balance/withdraw", handlerWithStorage.AddWithdrawal)
		userRouter.GET("/withdrawals", handlerWithStorage.GetWithdrawals)
	}

	router.HandleMethodNotAllowed = true

	return router
}
