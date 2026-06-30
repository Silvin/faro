// Package auth implementa la autenticación de Faro: usuarios, negocios (tenants),
// login con email+password, sesión por JWT en cookie httpOnly y scoping por tenant.
package auth

import "time"

// User representa a un usuario del sistema. password_hash nunca se expone.
type User struct {
	ID           string    `json:"id"`
	TenantID     *string   `json:"tenantId"` // nil solo para el super admin global
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	IsSuperAdmin bool      `json:"isSuperAdmin"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}
