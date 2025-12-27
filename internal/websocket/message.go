package websocket

type Message struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Room    string `json:"room"`
	Content []byte `json:"content"`
}
