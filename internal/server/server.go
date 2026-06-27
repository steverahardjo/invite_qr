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

type contextKey string

const loggerKey = contextKey("logger")

func NewLogger(dev bool, encoding string) (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
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

func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return zap.L()
	}
	return logger
}

type WebServer struct {
	ip       string
	port     string
	listener net.Listener
}

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

func (s *WebServer) Addr() string {
	return net.JoinHostPort(s.ip, s.port)
}

func (s *WebServer) IP() string {
	return s.ip
}

func (s *WebServer) Port() string {
	return s.port
}
