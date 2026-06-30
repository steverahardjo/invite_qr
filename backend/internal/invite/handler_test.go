package invite

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockInviteService struct {
	bulkErr error
	sendErr error
}

func (m *mockInviteService) BulkSendInvite(_ context.Context, _ string) error {
	return m.bulkErr
}

func (m *mockInviteService) SendInviteOnetime(_ context.Context, _ int32, _, _, _, _ string) error {
	return m.sendErr
}

func TestHandleBulkInvite_Success(t *testing.T) {
	mock := &mockInviteService{}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/bulk-invite", nil)
	rec := httptest.NewRecorder()

	h.HandleBulkInvite("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandleBulkInvite_ServiceError(t *testing.T) {
	mock := &mockInviteService{bulkErr: errors.New("send failed")}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/bulk-invite", nil)
	rec := httptest.NewRecorder()

	h.HandleBulkInvite("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestHandleSendInviteOnetime_Success(t *testing.T) {
	mock := &mockInviteService{}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/send-invite?guest_id=1&email=a@b.com&name=Alice", nil)
	rec := httptest.NewRecorder()

	h.HandleSendInviteOnetime("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandleSendInviteOnetime_SuccessWithWA(t *testing.T) {
	mock := &mockInviteService{}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/send-invite?guest_id=1&wa_number=123&name=Bob", nil)
	rec := httptest.NewRecorder()

	h.HandleSendInviteOnetime("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandleSendInviteOnetime_MissingGuestID(t *testing.T) {
	h := NewHandler(&mockInviteService{})

	req := httptest.NewRequest(http.MethodGet, "/send-invite", nil)
	rec := httptest.NewRecorder()

	h.HandleSendInviteOnetime("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleSendInviteOnetime_InvalidGuestID(t *testing.T) {
	h := NewHandler(&mockInviteService{})

	req := httptest.NewRequest(http.MethodGet, "/send-invite?guest_id=abc", nil)
	rec := httptest.NewRecorder()

	h.HandleSendInviteOnetime("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleSendInviteOnetime_MissingContact(t *testing.T) {
	h := NewHandler(&mockInviteService{})

	req := httptest.NewRequest(http.MethodGet, "/send-invite?guest_id=1&name=Charlie", nil)
	rec := httptest.NewRecorder()

	h.HandleSendInviteOnetime("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandleSendInviteOnetime_ServiceError(t *testing.T) {
	mock := &mockInviteService{sendErr: errors.New("send failed")}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/send-invite?guest_id=1&email=a@b.com", nil)
	rec := httptest.NewRecorder()

	h.HandleSendInviteOnetime("My Event").ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
