package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/LilLebowski/loyalty-system/internal/storage"
)

type HandlerWithStorage struct {
	storage *storage.Storage
	config  *config.Config
}

func Init(storage *storage.Storage, c *config.Config) *HandlerWithStorage {
	return &HandlerWithStorage{storage: storage, config: c}
}

func (strg *HandlerWithStorage) Register(ctx *gin.Context) {

}

func (strg *HandlerWithStorage) Login(ctx *gin.Context) {

}

func (strg *HandlerWithStorage) AddOrder(ctx *gin.Context) {

}

func (strg *HandlerWithStorage) GetOrders(ctx *gin.Context) {

}

func (strg *HandlerWithStorage) GetBalance(ctx *gin.Context) {

}

func (strg *HandlerWithStorage) AddWithdrawal(ctx *gin.Context) {

}

func (strg *HandlerWithStorage) GetWithdrawals(ctx *gin.Context) {

}
