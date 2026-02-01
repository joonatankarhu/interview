package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
