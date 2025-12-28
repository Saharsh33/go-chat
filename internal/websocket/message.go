package websocket

type MessageType string

const (
	// client → server
	MsgJoinRoom    MessageType = "join"
	MsgLeaveRoom   MessageType = "leave"
	MsgCreateRoom  MessageType = "create"
	MsgRoomMessage MessageType = "messageRoom"
	MsgBroadcast   MessageType = "broadcast"

	// server → client
	MsgSystem MessageType = "system" //not added for now
)

type Message struct {
	Type    MessageType `json:"type"`
	User    string      `json:"user"`
	Room    string      `json:"room"`
	Content string      `json:"content"`
}
