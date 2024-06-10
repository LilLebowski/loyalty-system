package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/LilLebowski/loyalty-system/internal/clients"
	"github.com/LilLebowski/loyalty-system/internal/db"
	"github.com/LilLebowski/loyalty-system/internal/router"
	"github.com/LilLebowski/loyalty-system/internal/services"
	"github.com/LilLebowski/loyalty-system/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.Init()

	err := db.RunMigrations(cfg)
	if err != nil {
		panic(err)
	}

	storageInstance := storage.Init(cfg.DBPath)

	service := services.OrderInit(storageInstance)
	client := resty.New()
	accrual := clients.AccrualInit(client, cfg.AccrualSysAddr)
	ticker := time.NewTicker(5 * time.Second)
	worker := services.NewPoolWorker(accrual, service)
	go func() {
		worker.StarIntegration(5, ticker)
	}()

	routerInstance := router.Init(storageInstance, cfg)

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: routerInstance,
	}

	go func() {
		log.Println(server.ListenAndServe())
		cancel()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
	}

	err = server.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}
