// Command api es el punto de entrada del backend de Faro (monolito modular).
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"faro/internal/auth"
	"faro/internal/categories"
	"faro/internal/config"
	"faro/internal/db"
	"faro/internal/products"
	"faro/internal/sales"
	"faro/internal/server"
)

func main() {
	cfg := config.Load()
	if cfg.UsesDefaultJWTSecret() {
		log.Println("ADVERTENCIA: usando JWT_SECRET por defecto (inseguro). Define JWT_SECRET en producción.")
	}

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	authSvc := auth.NewService(pool, cfg.JWTSecret, 8*time.Hour, cfg.CookieSecure)

	// Seed idempotente del super admin global (T2).
	if cfg.SuperAdminEmail != "" && cfg.SuperAdminPassword != "" {
		seedCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		created, err := authSvc.SeedSuperAdmin(seedCtx, cfg.SuperAdminEmail, cfg.SuperAdminPassword)
		cancel()
		switch {
		case err != nil:
			log.Printf("seed super admin: %v", err)
		case created:
			log.Printf("super admin global sembrado: %s", cfg.SuperAdminEmail)
		}
	}

	catSvc := categories.NewService(pool)
	prodSvc := products.NewService(pool)
	salesSvc := sales.NewService(pool)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           server.New(pool, cfg.CORSOrigin, authSvc, catSvc, prodSvc, salesSvc),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("faro api escuchando en :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	// Apagado ordenado (graceful shutdown).
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
