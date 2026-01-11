package websocket

type Room struct {
	name string
}
type RoomOps struct {
	clientDetails *Client
	roomDetails   int
	roomName      string
}
