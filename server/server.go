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

// roomExists returns true if a room with the given ID exists.
func roomExists(roomID int) bool {
	for _, room := range rooms {
		if room.ID == roomID {
			return true
		}
	}
	return false
}

// getBookingsByRoom returns all bookings for a room.
func getBookingsByRoom(roomID int) []Booking {
	result := []Booking{}
	for _, b := range bookings {
		if b.RoomID == roomID {
			result = append(result, b)
		}
	}
	return result
}

// findBookingIndex returns the index of a booking by ID, or -1 if not found.
func findBookingIndex(id int) int {
	for i, b := range bookings {
		if b.ID == id {
			return i
		}
	}
	return -1
}

// hasOverlappingBooking returns true if the room has a booking that overlaps with the given time range.
func hasOverlappingBooking(roomID int, start, end time.Time) bool {
	for _, b := range bookings {
		if b.RoomID != roomID {
			continue
		}
		bStart, _ := time.Parse(time.RFC3339, b.Start)
		bEnd, _ := time.Parse(time.RFC3339, b.End)
		if start.Before(bEnd) && end.After(bStart) {
			return true
		}
	}
	return false
}

// List bookings for a specific room by room ID
func HandleBookingsByRoom(w http.ResponseWriter, r *http.Request) {
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

	if !roomExists(roomID) {
		WriteError(w, "Room not found", http.StatusNotFound)
		return
	}

	WriteJSON(w, http.StatusOK, getBookingsByRoom(roomID))
}

// Create a new booking
func HandleCreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.RoomID == 0 || req.Start == "" || req.End == "" || req.User == "" {
		WriteError(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if !roomExists(req.RoomID) {
		WriteError(w, "Invalid roomId", http.StatusBadRequest)
		return
	}

	startDate, err1 := time.Parse(time.RFC3339, req.Start)
	endDate, err2 := time.Parse(time.RFC3339, req.End)
	if err1 != nil || err2 != nil {
		WriteError(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	if !startDate.Before(endDate) {
		WriteError(w, "Start time must be before end time", http.StatusBadRequest)
		return
	}

	if startDate.Before(time.Now()) {
		WriteError(w, "Cannot book in the past", http.StatusBadRequest)
		return
	}

	if hasOverlappingBooking(req.RoomID, startDate, endDate) {
		WriteError(w, "Room already booked for this time", http.StatusConflict)
		return
	}

	booking := Booking{
		ID:     nextID,
		RoomID: req.RoomID,
		Start:  req.Start,
		End:    req.End,
		User:   req.User,
	}
	nextID++
	bookings = append(bookings, booking)

	WriteJSON(w, http.StatusCreated, booking)
}

// Cancel a booking
func HandleDeleteBooking(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		WriteError(w, "Invalid id", http.StatusBadRequest)
		return
	}

	idx := findBookingIndex(id)
	if idx == -1 {
		WriteError(w, "Booking not found", http.StatusNotFound)
		return
	}

	bookings = append(bookings[:idx], bookings[idx+1:]...)
	WriteJSON(w, http.StatusOK, map[string]string{"message": "Booking deleted successfully"})
}
