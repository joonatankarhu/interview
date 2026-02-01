# AI Prompts Used to Generate This Project

---

## Initial Setup

**Example prompt:**

> Create a Go HTTP API for room bookings. Use the Chi router. Put the main entrypoint in `cmd/main.go` and server logic in a `server` package. Use in-memory storage (no database). List the API on port 3000.

---

## Core API Endpoints

**Example prompt:**

> Add these endpoints:
>
> - GET /rooms — list all rooms
> - GET /bookings — list all bookings
> - GET /bookings/room?roomId=1 — list all bookings for a room
> - POST /bookings — create a new booking (JSON body: roomId, start, end, user)
> - DELETE /bookings?id=1 — cancel a booking by ID

---

## Data Models & Validation

**Example prompt:**

> Add Room and Booking structs. Rooms have id and name. Bookings have id, roomId, start, end, user. Validate that roomId exists, dates are RFC3339, start is before end, and no overlapping bookings for the same room. Reject bookings in the past.

---

## Error Handling

**Example prompt:**

> Return JSON error messages for invalid roomId, missing fields, invalid dates, and overlapping bookings. Use proper HTTP status codes (400, 404, 409).

---

## Delete Success Message

**Prompt:**

> After deleting a booking, give a success message.

---

## Code Formatting

**Prompt:**

> Can these be formatted like var ( variables, variables )? @server.go (10-16)

---

## Move Helper Functions

**Prompt:**

> Move the helper funcs to helpers.go @server.go (1-159)

---

## Testing

**Prompt:**

> Add testing, add test for GET /rooms endpoint.

I asked the AI to add testing to the project: a test for the GET /rooms endpoint. The AI added `server/server_test.go` with a test that checks the endpoint returns status 200 and the expected room list, and a Makefile with `make test` to run tests quickly.
