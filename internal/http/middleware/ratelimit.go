package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	count int
	reset int64
}

var rl = struct {
	sync.Mutex
	data map[string]*bucket
}{data: map[string]*bucket{}}

func RateLimitLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		now := time.Now().Unix()

		rl.Lock()
		b, ok := rl.data[ip]
		if !ok || now > b.reset {
			rl.data[ip] = &bucket{count: 1, reset: now + 300}
			rl.Unlock()
		} else {
			if b.count >= 5 {
				rl.Unlock()
				http.Error(w, "rate_limit", 429)
				return
			}
			b.count++
			rl.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
