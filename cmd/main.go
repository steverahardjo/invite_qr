// Command invite_qr is the entry point for the QR Invite backend server.
package main

import (
	"bufio"
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"invite_qr/db/db_gen"
	"invite_qr/internal/admin"
	"invite_qr/internal/auth"
	"invite_qr/internal/config"
	"invite_qr/internal/invite"
	"invite_qr/internal/public"
	"invite_qr/internal/server"

	"go.uber.org/zap"
)

func main() {
	loadEnv()

	logger := SetLogger()
	defer logger.Sync()

	ctx := SetContext(logger)
	cfg := SetConfig()

	dbConn := SetDatabase(ctx, cfg, logger)
	defer dbConn.Close(logger)

	queries := db.New(dbConn.Conn)

	inviteSvc := SetInviteService(dbConn.Conn)
	jwtSvc := auth.NewJwtService(
		os.Getenv("JWT_SECRET"),
		24*time.Hour,
		"invite-qr",
	)

	if pass := os.Getenv("ADMIN_PASSWORD"); pass != "" {
		jwtSvc.SetPasswordHashEnv(pass)
	}

	authH := auth.NewJwtHandler(jwtSvc)
	adminH := admin.NewHandler(queries)
	publicH := public.NewHandler(public.NewService(queries))
	inviteH := invite.NewHandler(inviteSvc)

	mux := http.NewServeMux()
	SetRoutes(mux, ctx, authH, adminH, publicH, inviteH, jwtSvc.SecretKey)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ws, err := server.New(port)
	if err != nil {
		logger.Fatal("failed to create server", zap.Error(err))
	}

	logger.Info("starting server", zap.String("addr", ws.Addr()))
	if err := ws.ServeHTTPHandler(ctx, mux); err != nil {
		logger.Fatal("server error", zap.Error(err))
	}
}

// loadEnv reads .env file and sets environment variables.
// Skips empty lines and comments (#).
func loadEnv() {
	f, err := os.Open(".env")
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.IndexByte(line, '='); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			os.Setenv(key, val)
		}
	}
}

// SetLogger initializes the zap logger in JSON format for production use.
func SetLogger() *zap.Logger {
	logger, err := server.NewLogger(false, "json")
	if err != nil {
		panic(err)
	}
	return logger
}

// SetContext returns a context that is cancelled on SIGINT or SIGTERM.
func SetContext(logger *zap.Logger) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))
		cancel()
	}()

	return server.WithLogger(ctx, logger)
}

// SetConfig loads the application configuration from environment variables.
func SetConfig() *config.Config {
	return &config.Config{
		Name:              os.Getenv("DB_NAME"),
		User:              os.Getenv("DB_USER"),
		Host:              os.Getenv("DB_HOST"),
		Port:              os.Getenv("DB_PORT"),
		SSLMode:           os.Getenv("DB_SSLMODE"),
		ConnectionTimeout: 10,
		Password:          os.Getenv("DB_PASSWORD"),
		PoolMinConnections: 5,
		PoolMaxConnections: 25,
	}
}

// SetDatabase connects to PostgreSQL and verifies connectivity.
func SetDatabase(ctx context.Context, cfg *config.Config, logger *zap.Logger) *config.DB {
	dbConn, err := config.NewDBFromEnv(ctx, cfg, logger)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	return dbConn
}

// SetInviteService initializes the invite service with WhatsApp and email
// senders if their respective API keys are configured, and a default send
// window extending to the year 2099.
func SetInviteService(dbConn *sql.DB) *invite.Service {
	waSender := invite.InitWhatsappSender(os.Getenv("WA_PHONE"))
	emailSender := invite.InitEmailSender(os.Getenv("RESEND_EMAIL"))

	baseURL := os.Getenv("BASE_WEB_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080/"
	}
	parsedURL, _ := url.Parse(baseURL)

	return invite.NewService(
		dbConn,
		waSender,
		emailSender,
		time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC),
		parsedURL,
	)
}

// SetRoutes registers all HTTP handlers on the provided mux.
// Admin routes are protected by JWT auth middleware.
func SetRoutes(
	mux *http.ServeMux,
	ctx context.Context,
	authH *auth.JwtHandler,
	adminH *admin.Handler,
	publicH *public.Handler,
	inviteH *invite.Handler,
	secretKey []byte,
) {
	// public
	mux.HandleFunc("GET /api/invite/{token}", publicH.HandleGetInvite())
	mux.HandleFunc("GET /api/user", publicH.GetUserDetails())
	mux.HandleFunc("GET /api/qr", publicH.SendQRCode())

	// auth
	mux.HandleFunc("POST /api/admin/login", authH.LoginAdmin())

	// admin (protected)
	authMw := auth.AuthMiddleware(secretKey)
	mux.Handle("GET /api/admin/participants", authMw(adminH.ListParticipants()))
	mux.Handle("POST /api/admin/participants", authMw(adminH.AddParticipant()))
	mux.Handle("POST /api/admin/attendance", authMw(adminH.MarkAttendance()))

	// invite
	mux.HandleFunc("POST /api/bulk-invite", inviteH.HandleBulkInvite(ctx))
	mux.HandleFunc("GET /api/send-invite", inviteH.HandleSendInviteOnetime("My Event"))
}
