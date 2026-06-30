// Package config carga la configuración del backend desde variables de entorno.
package config

import "os"

// Config agrupa la configuración de runtime del servicio.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Load lee la configuración del entorno con valores por defecto para desarrollo.
func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://faro:faro@localhost:5432/faro?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET", "dev-insecure-secret-change-me"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
