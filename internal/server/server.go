// Package server provides application-level infrastructure including
// structured logging via zap, context-based logger injection, and a
// WebServer wrapper for graceful HTTP server lifecycle management.
package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is an unexported type used for context value keys to avoid collisions.
type contextKey string

// loggerKey is the context key used to store the zap logger.
const loggerKey = contextKey("logger")

// NewLogger builds a zap.Logger with ISO8601 timestamps, configurable
// development mode, and encoding (json or console). Output goes to stdout
// and stderr.
func NewLogger(dev bool, encoding string) (*zap.Logger, error) {
	encoderCfg := zap.NewDevelopmentConfig().EncoderConfig
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Encoding:         encoding,
		Development:      dev,
	}

	return config.Build()
}

// WithLogger returns a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// LoggerFromContext retrieves the zap.Logger from the context. If none is
// found, it returns the global default logger via zap.L().
func LoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return zap.L()
	}
	return logger
}

// WebServer wraps a net.Listener and provides methods to serve HTTP traffic
// with graceful shutdown on context cancellation.
type WebServer struct {
	ip       string
	port     string
	listener net.Listener
}

// New creates a WebServer by listening on the given port. The actual port
// number is captured from the listener, so passing ":0" gives a free port.
func New(port string) (*WebServer, error) {
	addr := ":" + port

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener on %s: %w", addr, err)
	}

	tcpAddr := listener.Addr().(*net.TCPAddr)

	return &WebServer{
		ip:       tcpAddr.IP.String(),
		port:     fmt.Sprintf("%d", tcpAddr.Port),
		listener: listener,
	}, nil
}

// ServeHTTP starts the given http.Server on the WebServer's listener. It runs
// a goroutine that waits for the context to be cancelled, then issues a graceful
// shutdown with a 5-second timeout. Blocks until the server stops or fails.
func (s *WebServer) ServeHTTP(
	ctx context.Context,
	srv *http.Server,
) error {
	logger := LoggerFromContext(ctx)

	errCh := make(chan error, 1)

	go func() {
		<-ctx.Done()

		logger.Info("server context cancelled")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		logger.Info("shutting down http server")

		errCh <- srv.Shutdown(shutdownCtx)
	}()

	logger.Info(
		"starting http server",
		zap.String("addr", s.Addr()),
	)

	if err := srv.Serve(s.listener); err != nil &&
		err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve: %w", err)
	}

	logger.Info("http server stopped")

	return <-errCh
}

// ServeHTTPHandler is a convenience wrapper that creates an http.Server with
// a 10-second ReadHeaderTimeout and delegates to ServeHTTP.
func (s *WebServer) ServeHTTPHandler(
	ctx context.Context,
	handler http.Handler,
) error {
	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handler,
	}

	return s.ServeHTTP(ctx, srv)
}

// Addr returns the listener address in host:port format.
func (s *WebServer) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

// IP returns the listener IP address.
func (s *WebServer) IP() string {
	return s.ip
}

// Port returns the listener port number.
func (s *WebServer) Port() string {
	return s.port
}
