package websocket

type Hub struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:

			h.Clients[client] = true

		case client := <-h.Unregister:

			_, exists := h.Clients[client] //break to key,value pair
			if exists {                    //if client exists then remove from map and close channel
				delete(h.Clients, client)
				close(client.Send)
			}

		case message := <-h.Broadcast: //broadcast msg is sent

			for client := range h.Clients { //for all clients
				select {
				case client.Send <- message:
				default:
					close(client.Send) //if client is slow or dead
					delete(h.Clients, client)
				}
			}
		}
	}
}
