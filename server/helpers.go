package server

import (
	"encoding/json"
	"net/http"
	"time"
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

// WriteError sends a JSON error response with the given message and status code.
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// WriteJSON sends a JSON response with the given status code and data.
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
