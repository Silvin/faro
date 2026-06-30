// Package config carga la configuración del backend desde variables de entorno.
package config

import "os"

// defaultJWTSecret es un valor inseguro solo para desarrollo. En producción
// DEBE definirse JWT_SECRET con un secreto fuerte.
const defaultJWTSecret = "dev-insecure-secret-change-me"

// Config agrupa la configuración de runtime del servicio.
type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	CORSOrigin         string // origen del frontend (faro-ui) permitido por CORS
	CookieSecure       bool   // true en prod (HTTPS); false en dev local (HTTP)
	SuperAdminEmail    string // seed del super admin global
	SuperAdminPassword string
	UploadDir          string // directorio local para imágenes subidas
}

// Load lee la configuración del entorno con valores por defecto para desarrollo.
func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://faro:faro@localhost:5432/faro?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET", defaultJWTSecret),
		CORSOrigin:  getenv("CORS_ORIGIN", "http://localhost:3000"),

		CookieSecure:       getbool("COOKIE_SECURE", false),
		SuperAdminEmail:    getenv("FARO_SUPERADMIN_EMAIL", ""),
		SuperAdminPassword: getenv("FARO_SUPERADMIN_PASSWORD", ""),
		UploadDir:          getenv("UPLOAD_DIR", "./uploads"),
	}
}

// UsesDefaultJWTSecret indica si se está usando el secreto inseguro por defecto.
func (c Config) UsesDefaultJWTSecret() bool {
	return c.JWTSecret == defaultJWTSecret
}

func getbool(key string, fallback bool) bool {
	switch os.Getenv(key) {
	case "1", "true", "TRUE", "yes":
		return true
	case "0", "false", "FALSE", "no":
		return false
	default:
		return fallback
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
