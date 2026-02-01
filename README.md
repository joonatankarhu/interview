# Bookings API (Go)

## How to Run Locally

1. **Install Go** (if not already): https://go.dev/dl/
2. Clone/download this repository.
3. Install dependencies:
   - `go mod tidy`
4. Run the server:
   - `go run cmd/main.go`
5. The server will start on `http://localhost:3000`

## API Endpoints

### Rooms

- **GET /rooms**
  - List all rooms.

### Bookings

- **GET /bookings**
  - List all bookings.
- **GET /bookings/room?roomId=ID**
  - List bookings for a specific room. Example: `/bookings/room?roomId=1`
- **POST /bookings**
  - Create a new booking.
  - Body (JSON):
    ```json
    {
      "roomId": 1,
      "start": "2026-02-01T15:00:00Z",
      "end": "2026-02-01T16:00:00Z",
      "user": "Alice"
    }
    ```
  - Content-Type: `application/json`
- **DELETE /bookings?id=ID**
  - Cancel a booking by ID. Example: `/bookings?id=2`
