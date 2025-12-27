package websocket

type Message struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	Content string `json:"content"`
}