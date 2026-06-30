package auth

import (
	"testing"
	"time"
)

func TestRateLimiterAllowsThenBlocks(t *testing.T) {
	rl := newRateLimiter(3, time.Minute)
	for i := 1; i <= 3; i++ {
		if !rl.allow("k") {
			t.Fatalf("el intento %d debería permitirse", i)
		}
	}
	if rl.allow("k") {
		t.Fatal("el 4º intento debería bloquearse")
	}
	if !rl.allow("otra-clave") {
		t.Fatal("otra clave no debería verse afectada")
	}
}

func TestRateLimiterResetsTrasVentana(t *testing.T) {
	rl := newRateLimiter(1, 10*time.Millisecond)
	if !rl.allow("k") {
		t.Fatal("primer intento permitido")
	}
	if rl.allow("k") {
		t.Fatal("segundo intento bloqueado dentro de la ventana")
	}
	time.Sleep(15 * time.Millisecond)
	if !rl.allow("k") {
		t.Fatal("tras la ventana debería permitirse de nuevo")
	}
}
