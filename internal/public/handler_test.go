package public

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	db_gen "invite_qr/db/db_gen"

	"github.com/google/uuid"
)

type mockPublicService struct {
	participant *db_gen.Participant
	err         error
}

func (m *mockPublicService) GetParticipantByExternalID(_ context.Context, _ string) (*db_gen.Participant, error) {
	return m.participant, m.err
}

func TestHandleGetInvite_Success(t *testing.T) {
	id := uuid.New()
	mock := &mockPublicService{
		participant: &db_gen.Participant{
			ExternalID: id,
			Name:       "Alice",
			Email:      "alice@test.com",
			WaNumber:   "123",
		},
	}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/invite/"+id.String(), nil)
	req.SetPathValue("token", id.String())
	rec := httptest.NewRecorder()

	h.HandleGetInvite().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got db_gen.Participant
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Name != "Alice" {
		t.Fatalf("expected name Alice, got %s", got.Name)
	}
}

func TestHandleGetInvite_MissingToken(t *testing.T) {
	h := NewHandler(&mockPublicService{})

	req := httptest.NewRequest(http.MethodGet, "/invite/", nil)
	rec := httptest.NewRecorder()

	h.HandleGetInvite().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleGetInvite_ServiceError(t *testing.T) {
	mock := &mockPublicService{err: errors.New("not found")}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/invite/some-token", nil)
	req.SetPathValue("token", "some-token")
	rec := httptest.NewRecorder()

	h.HandleGetInvite().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestGetUserDetails_Success(t *testing.T) {
	id := uuid.New()
	mock := &mockPublicService{
		participant: &db_gen.Participant{
			ExternalID: id,
			Name:       "Bob",
		},
	}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/user?id="+id.String(), nil)
	rec := httptest.NewRecorder()

	h.GetUserDetails().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got db_gen.Participant
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Name != "Bob" {
		t.Fatalf("expected name Bob, got %s", got.Name)
	}
}

func TestGetUserDetails_MissingID(t *testing.T) {
	h := NewHandler(&mockPublicService{})

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	rec := httptest.NewRecorder()

	h.GetUserDetails().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestSendQRCode_Success(t *testing.T) {
	h := NewHandler(&mockPublicService{})

	req := httptest.NewRequest(http.MethodGet, "/send-qr?participant_id=some-uuid", nil)
	rec := httptest.NewRecorder()

	h.SendQRCode().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "image/png" {
		t.Fatalf("expected Content-Type image/png, got %s", ct)
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty response body")
	}
}

func TestSendQRCode_MissingID(t *testing.T) {
	h := NewHandler(&mockPublicService{})

	req := httptest.NewRequest(http.MethodGet, "/send-qr", nil)
	rec := httptest.NewRecorder()

	h.SendQRCode().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
