package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"eth-pulse/backend/internal/supastore"
	"eth-pulse/backend/internal/worker"
)

func main() {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	alchemyURL := mustGetEnv("ALCHEMY_WSS_URL")
	supabaseURL := mustGetEnv("SUPABASE_URL")
	supabaseServiceRole := mustGetEnv("SUPABASE_SERVICE_ROLE_KEY")
	minETH := getEnvFloat("WHALE_MIN_ETH", 5.0)

	supabaseClient, err := supastore.New(supabaseURL, supabaseServiceRole)
	if err != nil {
		log.Fatalf("failed to init supabase: %v", err)
	}

	w := worker.NewAlchemyWorker(alchemyURL, supabaseClient, minETH)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			if err := w.Start(ctx); err != nil {
				log.Printf("worker disconnected, retrying in 2s: %v", err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):
				}
				continue
			}
			return
		}
	}()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"ok": true})
	})
	e.GET("/metrics", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"latest_gas_gwei": w.LatestGasGwei(),
			"min_eth":         minETH,
		})
	})

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("echo start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}

func mustGetEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("missing required env var: %s", k)
	}
	return v
}

func getEnv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func getEnvFloat(k string, fallback float64) float64 {
	raw := os.Getenv(k)
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return fallback
	}
	return parsed
}
