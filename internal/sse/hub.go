package sse

type Hub struct {
	Rooms      map[string]map[chan []byte]bool
	Broadcast  chan RoomMessage
	Register   chan Subscription
	Unregister chan Subscription
}

type RoomMessage struct {
	Room    string
	Message []byte
}

type Subscription struct {
	Room   string
	Client chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[chan []byte]bool),
		Broadcast:  make(chan RoomMessage),
		Register:   make(chan Subscription),
		Unregister: make(chan Subscription),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.Register:
			if _, ok := h.Rooms[sub.Room]; !ok {
				h.Rooms[sub.Room] = make(map[chan []byte]bool)
			}
			h.Rooms[sub.Room][sub.Client] = true

		case sub := <-h.Unregister:
			if clients, ok := h.Rooms[sub.Room]; ok {
				if _, ok := clients[sub.Client]; ok {
					delete(clients, sub.Client)
					close(sub.Client)
					if len(clients) == 0 {
						delete(h.Rooms, sub.Room)
					}
				}
			}

		case msg := <-h.Broadcast:
			if clients, ok := h.Rooms[msg.Room]; ok {
				for client := range clients {
					select {
					case client <- msg.Message:
					default:
						delete(clients, client)
						close(client)
					}
				}
			}
		}
	}
}
