package admin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	db "invite_qr/db/db_gen"

	"github.com/google/uuid"
)

type mockAdminStore struct {
	participants []db.Participant
	inserted     *db.Participant
	updated      *db.Participant
	err          error
}

func (m *mockAdminStore) ListParticipants(_ context.Context, _ db.ListParticipantsParams) ([]db.Participant, error) {
	return m.participants, m.err
}

func (m *mockAdminStore) InsertParticipant(_ context.Context, _ db.InsertParticipantParams) (db.Participant, error) {
	if m.inserted != nil {
		return *m.inserted, m.err
	}
	return db.Participant{}, m.err
}

func (m *mockAdminStore) UpdateParticipantAccessedByExternalID(_ context.Context, _ uuid.UUID) (db.Participant, error) {
	if m.updated != nil {
		return *m.updated, m.err
	}
	return db.Participant{}, m.err
}

func TestListParticipants_Success(t *testing.T) {
	mock := &mockAdminStore{
		participants: []db.Participant{
			{ID: 1, Name: "Alice", Email: "alice@test.com", WaNumber: "123"},
		},
	}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/participants", nil)
	rec := httptest.NewRecorder()

	h.ListParticipants().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got []db.Participant
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(got) != 1 || got[0].Name != "Alice" {
		t.Fatalf("expected 1 participant named Alice, got %+v", got)
	}
}

func TestListParticipants_StoreError(t *testing.T) {
	mock := &mockAdminStore{err: errors.New("db error")}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/participants", nil)
	rec := httptest.NewRecorder()

	h.ListParticipants().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestAddParticipant_Success(t *testing.T) {
	mock := &mockAdminStore{
		inserted: &db.Participant{ID: 1, Name: "Bob", Email: "bob@test.com", WaNumber: "456"},
	}
	h := NewHandler(mock)

	body := `{"name":"Bob","email":"bob@test.com","wa_number":"456"}`
	req := httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.AddParticipant().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestAddParticipant_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockAdminStore{})

	req := httptest.NewRequest(http.MethodPost, "/participants", strings.NewReader(`{bad}`))
	rec := httptest.NewRecorder()

	h.AddParticipant().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMarkAttendance_Success(t *testing.T) {
	mock := &mockAdminStore{
		updated: &db.Participant{ID: 1, Accessed: true},
	}
	h := NewHandler(mock)

	id := uuid.New()
	body := `{"participant_id":"` + id.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/attendance", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.MarkAttendance().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestMarkAttendance_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockAdminStore{})

	req := httptest.NewRequest(http.MethodPost, "/admin/attendance", strings.NewReader(`{bad}`))
	rec := httptest.NewRecorder()

	h.MarkAttendance().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMarkAttendance_MissingID(t *testing.T) {
	h := NewHandler(&mockAdminStore{})

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/admin/attendance", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.MarkAttendance().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMarkAttendance_InvalidUUID(t *testing.T) {
	h := NewHandler(&mockAdminStore{})

	body := `{"participant_id":"not-a-uuid"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/attendance", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.MarkAttendance().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestMarkAttendance_StoreError(t *testing.T) {
	mock := &mockAdminStore{err: errors.New("db error")}
	h := NewHandler(mock)

	id := uuid.New()
	body := `{"participant_id":"` + id.String() + `"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/attendance", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.MarkAttendance().ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}
