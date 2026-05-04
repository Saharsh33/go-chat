package websocket

type MessageType string

const (
	// client → server
	MsgJoinRoom           MessageType = "join"
	MsgLeaveRoom          MessageType = "leave"
	MsgCreateRoom         MessageType = "create"
	MsgRoomMessage        MessageType = "messageRoom"
	MsgDirectMessage      MessageType = "messageDirect"
	MsgBroadcast          MessageType = "broadcast"
	MsgNextRoomMessages   MessageType = "nextRoomMsgs"
	MsgNextDirectMessages MessageType = "nextDirectMsgs"
	MsgGetConversations   MessageType = "getConversations"
	// server → client
	MsgSystem             MessageType = "system"
	MsgConversationsList  MessageType = "conversationsList"
)

type Message struct {
	Type           MessageType `json:"type"`
	User           string      `json:"user,omitempty"`
	Room           int         `json:"room,omitempty"`
	Content        string      `json:"content,omitempty"`
	Receiver       string      `json:"receiver,omitempty"`
	ConversationID int         `json:"conversationId,omitempty"`
}
