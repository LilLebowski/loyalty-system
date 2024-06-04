package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/LilLebowski/loyalty-system/cmd/gophermart/config"
	"github.com/LilLebowski/loyalty-system/internal/router"
	"github.com/LilLebowski/loyalty-system/internal/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := config.Init()

	storageInstance := storage.Init(cfg.DBPath)

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

	err := server.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}
