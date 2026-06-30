// Package products implementa el CRUD del catálogo de productos por negocio,
// reutilizando la sesión/tenant-scope de auth y las categorías (M2).
package products

import "time"

// Product es un producto del catálogo de un negocio. price_cents es entero.
type Product struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenantId"`
	CategoryID   *string   `json:"categoryId"`
	CategoryName *string   `json:"categoryName"`
	Name         string    `json:"name"`
	PriceCents   int       `json:"priceCents"`
	Status       string    `json:"status"`
	ImageURL     *string   `json:"imageUrl"`
	CreatedAt    time.Time `json:"createdAt"`
}
