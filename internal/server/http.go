package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPServer struct {
	Host         string
	Port         string
	Handler      http.Handler
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ErrorLog     *log.Logger
	logger       *slog.Logger
}

type Options func(*HTTPServer)

func WithHost(host string) Options {
	return func(h *HTTPServer) {
		h.Host = host
	}
}

func WithPort(port string) Options {
	return func(h *HTTPServer) {
		h.Port = port
	}
}

func WithHandler(handler http.Handler) Options {
	return func(h *HTTPServer) {
		h.Handler = handler
	}
}

func WithIdleTimeout(idleTimeout time.Duration) Options {
	return func(h *HTTPServer) {
		h.IdleTimeout = idleTimeout
	}
}

func WithReadTimeout(readTimeout time.Duration) Options {
	return func(h *HTTPServer) {
		h.ReadTimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout time.Duration) Options {
	return func(h *HTTPServer) {
		h.WriteTimeout = writeTimeout
	}
}

func WithErrorLog(errorLog *log.Logger) Options {
	return func(s *HTTPServer) {
		s.ErrorLog = errorLog
	}
}

func WithLogger(logger *slog.Logger) Options {
	return func(h *HTTPServer) {
		h.logger = logger
	}
}

func (h *HTTPServer) Connect() error {
	addr := fmt.Sprintf("%s:%s", h.Host, h.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      h.Handler,
		IdleTimeout:  h.IdleTimeout,
		ReadTimeout:  h.ReadTimeout,
		WriteTimeout: h.WriteTimeout,
		ErrorLog:     h.ErrorLog,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		se := <-quit
		h.logger.Info("caught signal", "signal", se.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		h.logger.Info("completing background tasks", "addr", server.Addr)
		shutdownError <- nil
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdownError; err != nil {
		return err
	}

	h.logger.Info("stopped server", "addr", server.Addr)

	return nil
}

func NewHTTPServer(opts ...Options) *HTTPServer {
	httpServer := &HTTPServer{}
	for _, opt := range opts {
		opt(httpServer)
	}
	return httpServer
}
