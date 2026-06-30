package api

import (
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

func New(mux *http.ServeMux) huma.API {
	config := huma.DefaultConfig("Invite QR API", "1.0.0")
	config.Info.Description = "API for managing wedding invitations, QR codes, and guest attendance"

	baseURL := os.Getenv("BASE_WEB_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	config.Servers = []*huma.Server{{URL: baseURL}}

	api := humago.New(mux, config)
	return api
}
