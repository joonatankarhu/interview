package main

import (
	"bookings/server"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/bookings/room", server.HandleBookingsByRoom)
	r.Post("/bookings", server.HandleCreateBooking)
	r.Delete("/bookings", server.HandleDeleteBooking)

	log.Println("Server started on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
