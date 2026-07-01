package reports

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"faro/internal/auth"
)

// Routes se monta en /reports. Requiere sesión.
func (svc *Service) Routes(requireSession func(http.Handler) http.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(requireSession)
	r.Get("/sales", svc.handleSalesReport)
	return r
}

func (svc *Service) handleSalesReport(w http.ResponseWriter, r *http.Request) {
	u, ok := auth.UserFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "No autenticado")
		return
	}
	if u.TenantID == nil {
		writeError(w, http.StatusBadRequest, "tenant_required", "Esta operación requiere un negocio")
		return
	}

	// Rango: por defecto hoy (hora local del servidor); sobreescribible por from/to.
	now := time.Now()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	to := from.AddDate(0, 0, 1)
	if f := parseTime(r.URL.Query().Get("from")); f != nil {
		from = *f
	}
	if t := parseTime(r.URL.Query().Get("to")); t != nil {
		to = *t
	}
	tz, _ := strconv.Atoi(r.URL.Query().Get("tz"))

	rep, err := svc.SalesReport(r.Context(), *u.TenantID, from, to, tz)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "No se pudo generar el reporte")
		return
	}
	writeJSON(w, http.StatusOK, rep)
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}
