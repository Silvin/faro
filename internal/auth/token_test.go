package auth

import (
	"testing"
	"time"
)

func TestTokenIssueAndParse(t *testing.T) {
	tm := newTokenManager("secret", time.Hour)
	tid := "tenant-1"

	tok, exp, err := tm.issue(Claims{UserID: "u1", TenantID: &tid, IsSuperAdmin: false})
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if !exp.After(time.Now()) {
		t.Fatal("la expiración debería estar en el futuro")
	}

	c, err := tm.parse(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if c.UserID != "u1" {
		t.Fatalf("userID = %q, esperaba u1", c.UserID)
	}
	if c.TenantID == nil || *c.TenantID != "tenant-1" {
		t.Fatal("el tenantID no coincide")
	}
	if c.IsSuperAdmin {
		t.Fatal("isSuperAdmin debería ser false")
	}
}

func TestTokenSuperAdminNoTenant(t *testing.T) {
	tm := newTokenManager("secret", time.Hour)
	tok, _, err := tm.issue(Claims{UserID: "root", TenantID: nil, IsSuperAdmin: true})
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	c, err := tm.parse(tok)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if c.TenantID != nil {
		t.Fatal("el super admin global no debe tener tenant")
	}
	if !c.IsSuperAdmin {
		t.Fatal("isSuperAdmin debería ser true")
	}
}

func TestTokenRejectsWrongSecret(t *testing.T) {
	tok, _, _ := newTokenManager("secret", time.Hour).issue(Claims{UserID: "u1"})
	if _, err := newTokenManager("otro-secreto", time.Hour).parse(tok); err == nil {
		t.Fatal("debería rechazar un token firmado con otro secreto")
	}
}
