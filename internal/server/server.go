// Package server arma el router HTTP y el middleware base del backend.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"faro/internal/auth"
	"faro/internal/categories"
	"faro/internal/products"
	"faro/internal/sales"
)

// New construye el handler HTTP raíz. Los módulos (auth, products, …) montarán
// aquí sus sub-routers en incrementos siguientes. corsOrigin es el origen del
// frontend (faro-ui) autorizado a consumir la API con credenciales.
func New(pool *pgxpool.Pool, corsOrigin string, authSvc *auth.Service, catSvc *categories.Service, prodSvc *products.Service, salesSvc *sales.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS para el frontend (faro-ui), que vive en otro origen.
	// AllowCredentials=true para que viaje la cookie de sesión httpOnly.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{corsOrigin},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Liveness: el proceso está vivo.
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Readiness: además la base de datos responde.
	r.Get("/ready", func(w http.ResponseWriter, req *http.Request) {
		if err := pool.Ping(req.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "db_unavailable"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	// Módulo auth (login): /auth/login, /auth/logout, /auth/me
	r.Mount("/auth", authSvc.Routes())
	// Provisión: alta de negocios (super admin) y de usuarios (acotado al negocio).
	r.Mount("/tenants", authSvc.TenantRoutes())
	r.Mount("/users", authSvc.UserRoutes())
	// Módulo categorías (M2): CRUD acotado al negocio, protegido por sesión.
	r.Mount("/categories", catSvc.Routes(authSvc.RequireSession))
	// Módulo productos (M3): CRUD acotado al negocio, protegido por sesión.
	r.Mount("/products", prodSvc.Routes(authSvc.RequireSession))
	// Módulo POS (M4): ventas, total calculado en servidor.
	r.Mount("/sales", salesSvc.Routes(authSvc.RequireSession))

	return r
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
