package handler

import (
	"sync"
	"time"
)

type Visitor struct {
	rate       float64
	cap        float64
	tokens     float64
	lastAccess time.Time
	mu         sync.Mutex
}

var (
	visitors = make(map[string]*Visitor)
	muMap    sync.Mutex
)

func newVisitor(rate, cap float64) *Visitor {
	return &Visitor{
		rate:   rate,
		cap:    cap,
		tokens: cap,
	}
}

func (v *Visitor) Take(tokens float64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.tokens += v.rate * time.Since(v.lastAccess).Seconds()

	if v.tokens > v.cap {
		v.tokens = v.cap
	}

	if tokens <= v.tokens {
		v.tokens -= tokens
		v.lastAccess = time.Now()
		return true
	}
	return false

}

func getVisitor(ip string, rate, cap float64) *Visitor {
	muMap.Lock()
	defer muMap.Unlock()

	v, ok := visitors[ip]
	if !ok {
		limiter := newVisitor(rate, cap)
		visitors[ip] = limiter
		return limiter
	}

	return v
}

func addBlockList(ip string) {
	muMap.Lock()
	defer muMap.Unlock()

	if _, ok := visitors[ip]; ok {
		limiter := newVisitor(-1, -1)
		visitors[ip] = limiter
	}
}

func CleanupVisitors() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		muMap.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastAccess) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		muMap.Unlock()
	}
}
