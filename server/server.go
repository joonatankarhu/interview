package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

var (
	rooms = []Room{
		{ID: 1, Name: "Room A"},
		{ID: 2, Name: "Room B"},
		{ID: 3, Name: "Room C"},
	}
	bookings = []Booking{}
	nextID   = 1
)

// List bookings for a specific room by room ID
func HandleBookingsByRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roomIDStr := r.URL.Query().Get("roomId")
	if roomIDStr == "" {
		WriteError(w, "roomId required as query param", http.StatusBadRequest)
		return
	}
	roomID, err := strconv.Atoi(roomIDStr)
	if err != nil {
		WriteError(w, "Invalid roomId", http.StatusBadRequest)
		return
	}
	found := false
	for _, room := range rooms {
		if room.ID == roomID {
			found = true
			break
		}
	}
	if !found {
		WriteError(w, "Room not found", http.StatusNotFound)
		return
	}
	filtered := []Booking{}
	for _, b := range bookings {
		if b.RoomID == roomID {
			filtered = append(filtered, b)
		}
	}
	json.NewEncoder(w).Encode(filtered)
}

// Create a new booking
func HandleCreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct {
		RoomID int    `json:"roomId"`
		Start  string `json:"start"`
		End    string `json:"end"`
		User   string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.RoomID == 0 || req.Start == "" || req.End == "" || req.User == "" {
		WriteError(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	// Check room exists
	found := false
	for _, room := range rooms {
		if room.ID == req.RoomID {
			found = true
			break
		}
	}
	if !found {
		WriteError(w, "Invalid roomId", http.StatusBadRequest)
		return
	}
	// Parse dates
	startDate, err1 := time.Parse(time.RFC3339, req.Start)
	endDate, err2 := time.Parse(time.RFC3339, req.End)
	now := time.Now()
	if err1 != nil || err2 != nil {
		WriteError(w, "Invalid date format", http.StatusBadRequest)
		return
	}
	if !startDate.Before(endDate) {
		WriteError(w, "Start time must be before end time", http.StatusBadRequest)
		return
	}
	if startDate.Before(now) {
		WriteError(w, "Cannot book in the past", http.StatusBadRequest)
		return
	}
	// Check for overlapping bookings
	for _, b := range bookings {
		if b.RoomID != req.RoomID {
			continue
		}
		bStart, _ := time.Parse(time.RFC3339, b.Start)
		bEnd, _ := time.Parse(time.RFC3339, b.End)
		if startDate.Before(bEnd) && endDate.After(bStart) {
			WriteError(w, "Room already booked for this time", http.StatusConflict)
			return
		}
	}
	// Create and store booking
	booking := Booking{
		ID:     nextID,
		RoomID: req.RoomID,
		Start:  req.Start,
		End:    req.End,
		User:   req.User,
	}
	nextID++
	bookings = append(bookings, booking)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

// Cancel a booking
func HandleDeleteBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, "Invalid id", http.StatusBadRequest)
		return
	}
	idx := -1
	for i, b := range bookings {
		if b.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		WriteError(w, "Booking not found", http.StatusNotFound)
		return
	}
	bookings = append(bookings[:idx], bookings[idx+1:]...)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking deleted successfully"})
}