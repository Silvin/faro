package auth

import (
	"sync"
	"time"
)

// rateLimiter es un limitador de ventana fija en memoria. Suficiente para una
// sola instancia; en multi-instancia migrar a un store compartido (Redis) — deuda.
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string]*counter
	limit  int
	window time.Duration
}

type counter struct {
	count   int
	resetAt time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{hits: make(map[string]*counter), limit: limit, window: window}
}

// allow registra un intento para key y devuelve false si superó el límite en la ventana.
func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	c, ok := rl.hits[key]
	if !ok || now.After(c.resetAt) {
		rl.hits[key] = &counter{count: 1, resetAt: now.Add(rl.window)}
		return true
	}
	if c.count >= rl.limit {
		return false
	}
	c.count++
	return true
}
