package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func resetBookingsState() {
	bookings = []Booking{}
	nextID = 1
}

func TestHandleGetRooms(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/rooms", nil)
	rec := httptest.NewRecorder()

	HandleGetRooms(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var got []Room
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(got) != 3 {
		t.Errorf("got %d rooms, want 3", len(got))
	}

	wantNames := map[int]string{1: "Room A", 2: "Room B", 3: "Room C"}
	for _, room := range got {
		if wantNames[room.ID] != room.Name {
			t.Errorf("room %d: got name %q, want %q", room.ID, room.Name, wantNames[room.ID])
		}
	}
}

func TestHandleGetBookings(t *testing.T) {
	resetBookingsState()

	req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
	rec := httptest.NewRecorder()

	HandleGetBookings(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}

	var got []Booking
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(got) != 0 {
		t.Errorf("got %d bookings, want 0", len(got))
	}
}

func TestHandleGetBookingsByRoom(t *testing.T) {
	resetBookingsState()

	t.Run("missing roomId", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bookings/room", nil)
		rec := httptest.NewRecorder()
		HandleGetBookingsByRoom(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("invalid roomId", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bookings/room?roomId=abc", nil)
		rec := httptest.NewRecorder()
		HandleGetBookingsByRoom(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("room not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bookings/room?roomId=999", nil)
		rec := httptest.NewRecorder()
		HandleGetBookingsByRoom(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		bookings = []Booking{{ID: 1, RoomID: 1, Start: "2026-06-01T10:00:00Z", End: "2026-06-01T11:00:00Z", User: "alice"}}
		req := httptest.NewRequest(http.MethodGet, "/bookings/room?roomId=1", nil)
		rec := httptest.NewRecorder()
		HandleGetBookingsByRoom(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var got []Booking
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(got) != 1 || got[0].User != "alice" {
			t.Errorf("got %v, want one booking for alice", got)
		}
	})
}

func TestHandleCreateBooking(t *testing.T) {
	resetBookingsState()

	start := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	end := time.Now().Add(25 * time.Hour).Format(time.RFC3339)

	t.Run("invalid JSON", func(t *testing.T) {
		resetBookingsState()
		body := bytes.NewReader([]byte(`{invalid`))
		req := httptest.NewRequest(http.MethodPost, "/bookings", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		HandleCreateBooking(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		resetBookingsState()
		body := bytes.NewBufferString(`{"roomId":1}`)
		req := httptest.NewRequest(http.MethodPost, "/bookings", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		HandleCreateBooking(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		resetBookingsState()
		body := bytes.NewBufferString(`{"roomId":1,"start":"` + start + `","end":"` + end + `","user":"bob"}`)
		req := httptest.NewRequest(http.MethodPost, "/bookings", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		HandleCreateBooking(rec, req)
		if rec.Code != http.StatusCreated {
			t.Errorf("got status %d, want 201", rec.Code)
		}
		var got Booking
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if got.ID != 1 || got.RoomID != 1 || got.User != "bob" {
			t.Errorf("got booking %+v", got)
		}
		if len(bookings) != 1 {
			t.Errorf("bookings length %d, want 1", len(bookings))
		}
	})
}

func TestHandleDeleteBooking(t *testing.T) {
	resetBookingsState()

	t.Run("invalid id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/bookings?id=xyz", nil)
		rec := httptest.NewRecorder()
		HandleDeleteBooking(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want 400", rec.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		resetBookingsState()
		req := httptest.NewRequest(http.MethodDelete, "/bookings?id=1", nil)
		rec := httptest.NewRecorder()
		HandleDeleteBooking(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("got status %d, want 404", rec.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		resetBookingsState()
		bookings = []Booking{{ID: 1, RoomID: 1, Start: "2026-06-01T10:00:00Z", End: "2026-06-01T11:00:00Z", User: "alice"}}
		nextID = 2
		req := httptest.NewRequest(http.MethodDelete, "/bookings?id=1", nil)
		rec := httptest.NewRecorder()
		HandleDeleteBooking(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want 200", rec.Code)
		}
		var got map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if got["message"] != "Booking deleted successfully" {
			t.Errorf("got message %q", got["message"])
		}
		if len(bookings) != 0 {
			t.Errorf("bookings length %d, want 0", len(bookings))
		}
	})
}
