// Package categories implementa el CRUD de categorías de producto, acotado por
// negocio (tenant). Reutiliza la sesión y el tenant-scope del módulo auth.
package categories

import "time"

// Category es una categoría de producto de un negocio.
type Category struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	SortOrder int       `json:"sortOrder"`
	ImageURL  *string   `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
}
