package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

// fullRouter monta auth + provisión como en el servidor real.
func fullRouter(svc *Service) http.Handler {
	r := chi.NewRouter()
	r.Mount("/auth", svc.Routes())
	r.Mount("/tenants", svc.TenantRoutes())
	r.Mount("/users", svc.UserRoutes())
	return r
}

func jarClient(t *testing.T) *http.Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("cookiejar: %v", err)
	}
	return &http.Client{Jar: jar}
}

func post(t *testing.T, c *http.Client, url string, body any) *http.Response {
	t.Helper()
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	return resp
}

func login(t *testing.T, c *http.Client, base, email, pass string) {
	t.Helper()
	if resp := post(t, c, base+"/auth/login", map[string]string{"email": email, "password": pass}); resp.StatusCode != http.StatusOK {
		t.Fatalf("login %s: esperaba 200, obtuvo %d", email, resp.StatusCode)
	}
}

func usersCount(t *testing.T, c *http.Client, base string) int {
	t.Helper()
	resp, err := c.Get(base + "/users")
	if err != nil {
		t.Fatalf("GET /users: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /users: esperaba 200, obtuvo %d", resp.StatusCode)
	}
	var body struct {
		Items []map[string]any `json:"items"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return len(body.Items)
}

func TestLogoutInvalidatesSession(t *testing.T) {
	svc, pool := testService(t)
	defer pool.Close()
	if _, err := svc.SeedSuperAdmin(context.Background(), "admin@faro.test", "secret123"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	srv := httptest.NewServer(fullRouter(svc))
	defer srv.Close()
	c := jarClient(t)

	login(t, c, srv.URL, "admin@faro.test", "secret123")

	// Con sesión: /me = 200.
	if resp, err := c.Get(srv.URL + "/auth/me"); err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("/me con sesión: esperaba 200, obtuvo %v", resp.StatusCode)
	}

	// Logout = 204 y limpia la cookie.
	if resp := post(t, c, srv.URL+"/auth/logout", map[string]string{}); resp.StatusCode != http.StatusNoContent {
		t.Fatalf("logout: esperaba 204, obtuvo %d", resp.StatusCode)
	}

	// Tras logout: /me = 401.
	if resp, err := c.Get(srv.URL + "/auth/me"); err != nil || resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("/me tras logout: esperaba 401, obtuvo %v", resp.StatusCode)
	}
}

func TestProvisioningAndTenantIsolation(t *testing.T) {
	svc, pool := testService(t)
	defer pool.Close()
	if _, err := svc.SeedSuperAdmin(context.Background(), "root@faro.test", "secret123"); err != nil {
		t.Fatalf("seed: %v", err)
	}

	srv := httptest.NewServer(fullRouter(svc))
	defer srv.Close()

	// Super admin crea dos negocios con sus dueños.
	root := jarClient(t)
	login(t, root, srv.URL, "root@faro.test", "secret123")
	if resp := post(t, root, srv.URL+"/tenants", map[string]string{
		"name": "Cafe Uno", "ownerEmail": "owner@uno.test", "ownerPassword": "secret123", "ownerName": "Dueño Uno",
	}); resp.StatusCode != http.StatusCreated {
		t.Fatalf("crear negocio 1: esperaba 201, obtuvo %d", resp.StatusCode)
	}
	if resp := post(t, root, srv.URL+"/tenants", map[string]string{
		"name": "Cafe Dos", "ownerEmail": "owner@dos.test", "ownerPassword": "secret123", "ownerName": "Dueño Dos",
	}); resp.StatusCode != http.StatusCreated {
		t.Fatalf("crear negocio 2: esperaba 201, obtuvo %d", resp.StatusCode)
	}

	// Dueño 1 inicia sesión: solo ve su negocio (1 usuario).
	owner1 := jarClient(t)
	login(t, owner1, srv.URL, "owner@uno.test", "secret123")
	if n := usersCount(t, owner1, srv.URL); n != 1 {
		t.Fatalf("negocio 1 esperaba 1 usuario, obtuvo %d", n)
	}

	// Crea un barista en su negocio.
	if resp := post(t, owner1, srv.URL+"/users", map[string]string{
		"email": "barista@uno.test", "password": "secret123", "name": "Barista Uno",
	}); resp.StatusCode != http.StatusCreated {
		t.Fatalf("crear usuario: esperaba 201, obtuvo %d", resp.StatusCode)
	}

	// Aislamiento: el negocio 1 ahora tiene 2 usuarios; el negocio 2, solo su dueño.
	if n := usersCount(t, owner1, srv.URL); n != 2 {
		t.Fatalf("negocio 1 esperaba 2 usuarios, obtuvo %d", n)
	}
	owner2 := jarClient(t)
	login(t, owner2, srv.URL, "owner@dos.test", "secret123")
	if n := usersCount(t, owner2, srv.URL); n != 1 {
		t.Fatalf("negocio 2 esperaba 1 usuario (aislado), obtuvo %d", n)
	}

	// Un dueño (no super admin) no puede crear negocios -> 403.
	if resp := post(t, owner1, srv.URL+"/tenants", map[string]string{
		"name": "X", "ownerEmail": "x@x.test", "ownerPassword": "secret123", "ownerName": "X",
	}); resp.StatusCode != http.StatusForbidden {
		t.Fatalf("dueño creando negocio: esperaba 403, obtuvo %d", resp.StatusCode)
	}

	// Email duplicado -> 409.
	if resp := post(t, owner1, srv.URL+"/users", map[string]string{
		"email": "barista@uno.test", "password": "secret123", "name": "Repetido",
	}); resp.StatusCode != http.StatusConflict {
		t.Fatalf("email duplicado: esperaba 409, obtuvo %d", resp.StatusCode)
	}
}
