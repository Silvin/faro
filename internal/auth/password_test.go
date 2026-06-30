package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	h, err := hashPassword("s3cret-pass")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if h == "s3cret-pass" {
		t.Fatal("la contraseña no debe almacenarse en claro")
	}
	if !checkPassword(h, "s3cret-pass") {
		t.Fatal("debería validar la contraseña correcta")
	}
	if checkPassword(h, "otra-cosa") {
		t.Fatal("no debería validar una contraseña incorrecta")
	}
}
