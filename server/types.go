package server

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Booking struct {
	ID     int    `json:"id"`
	RoomID int    `json:"roomId"`
	Start  string `json:"start"`
	End    string `json:"end"`
	User   string `json:"user"`
}