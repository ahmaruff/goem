// Package middleware holds the app's Echo middlewares.
package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"maka-go/internal/logger"

	"github.com/labstack/echo/v5"
)

// LoggerConfig tunes request logging.
type LoggerConfig struct {
	// QuietPaths are path prefixes that are NOT logged when the request
	// succeeds (status < 400), for example "/api/v1/health" also covers
	// "/api/v1/health/live". This is the main noise remover: probes and other
	// periodic traffic. Failures on these paths are still logged.
	QuietPaths []string
}

// Logger emits exactly one structured line per request and injects a
// request-scoped logger into the request context, reachable in handlers with:
//
//	logger.Ctx(c.Request().Context()).Info().Msg("...")
//
// Levels are picked so noise stays low: 5xx -> error, 4xx -> warn,
// everything else -> info.
func Logger(cfg LoggerConfig) echo.MiddlewareFunc {
	quiet := make(map[string]struct{}, len(cfg.QuietPaths))
	for _, p := range cfg.QuietPaths {
		quiet[p] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()
			req := c.Request()

			path := c.Path()
			if path == "" {
				path = req.URL.Path
			}

			l := logger.L().With().
				Str("request_id", requestID(req)).
				Str("method", req.Method).
				Str("path", path).
				Str("remote_ip", c.RealIP()).
				Logger()

			// Enrich the request context so handlers share these fields.
			c.SetRequest(req.WithContext(logger.WithContext(req.Context(), l)))

			err := next(c)

			_, status := echo.ResolveResponseStatus(c.Response(), err)

			if isQuiet(quiet, path) && status < http.StatusBadRequest {
				return err
			}

			ev := l.Info()
			switch {
			case status >= http.StatusInternalServerError:
				ev = l.Error()
			case status >= http.StatusBadRequest:
				ev = l.Warn()
			}

			ev.Int("status", status).
				Dur("latency_ms", time.Since(start)).
				Msg("http request")

			return err
		}
	}
}

// isQuiet reports whether path matches a configured quiet path, exactly or as
// a sub-path of it (prefix match).
func isQuiet(quiet map[string]struct{}, path string) bool {
	if _, ok := quiet[path]; ok {
		return true
	}
	for p := range quiet {
		if p != "/" && strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// requestID reuses an upstream id (X-Request-ID) or mints a short new one so
// logs from one request can be stitched together.
func requestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}

	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return ""
	}
	return hex.EncodeToString(b[:])
}
