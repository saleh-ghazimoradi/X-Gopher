package middlewares

import (
	"fmt"
	"github.com/saleh-ghazimoradi/X-Gopher/config"
	"github.com/saleh-ghazimoradi/X-Gopher/internal/helper"
	"github.com/saleh-ghazimoradi/X-Gopher/utils"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Middleware struct {
	config *config.Config
	logger *slog.Logger
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Info("Incoming request: ", "method", r.Method, "path", r.URL.Path, "protocol", r.Proto, "remote_addr", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				helper.InternalServerError(w, "panic recovery failed", fmt.Errorf("%v", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	if !m.config.RateLimiter.Enabled {
		return next
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(m.config.RateLimiter.RPS), m.config.RateLimiter.Burst),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			helper.RateLimitExceededResponse(w, "Rate limit exceeded")
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			helper.UnauthorizedResponse(w, "Authorization header missing")
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			helper.UnauthorizedResponse(w, "Invalid authorization header")
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], m.config.JWT.Secret)
		if err != nil {
			helper.UnauthorizedResponse(w, "Invalid token")
			return
		}

		ctx := r.Context()
		ctx = utils.WithUserId(ctx, claims.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func NewMiddleware(config *config.Config, logger *slog.Logger) *Middleware {
	return &Middleware{
		config: config,
		logger: logger,
	}
}
