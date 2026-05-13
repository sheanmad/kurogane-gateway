package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kurogane/gateway/internal/config"
	"github.com/kurogane/gateway/internal/proxy"
	"github.com/kurogane/gateway/internal/router"
)

func main() {
	cfg := config.Load()

	os.Setenv("GIN_MODE", cfg.GinMode)

	p := proxy.New(cfg)
	r := router.Setup(p)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	go func() {
		log.Printf("gateway starting on :%s, proxying to %s", cfg.Port, cfg.EngineURL)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("gateway forced shutdown: %v", err)
	}

	fmt.Println("gateway stopped")
}