package hub

import (
	"encoding/json"
	"time"

	log "github.com/buglloc/simplelog"
	"github.com/gorilla/websocket"
)

type Message struct {
	Time  time.Time `json:"time"`
	QType string    `json:"type"`
	Name  string    `json:"name"`
	RR    string    `json:"rr"`
	Ok    bool      `json:"ok"`
}

type channelMessage struct {
	msg     Message
	channel string
}

type subscription struct {
	conn *connection
	id   string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	channels map[string]map[*connection]struct{}

	// Inbound messages from the connections.
	broadcast chan channelMessage

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan channelMessage),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	channels:   make(map[string]map[*connection]struct{}),
}

// TODO(buglloc): bullshit
func init() {
	go h.run()
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			connections := h.channels[s.id]
			if connections == nil {
				connections = make(map[*connection]struct{})
				h.channels[s.id] = connections
			}
			h.channels[s.id][s.conn] = struct{}{}
		case s := <-h.unregister:
			connections := h.channels[s.id]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.channels, s.id)
					}
				}
			}
		case m := <-h.broadcast:
			data, err := json.Marshal(m.msg)
			if err != nil {
				log.Error("can't marshal hub message", "name", m.msg.Name, "err", err)
				continue
			}

			connections := h.channels[m.channel]
			for c := range connections {
				select {
				case c.send <- data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.channels, m.channel)
					}
				}
			}
		}
	}
}

func Register(conn *websocket.Conn, id string) {
	s := subscription{
		conn: &connection{
			send: make(chan []byte, 256),
			ws:   conn,
		},
		id: id,
	}

	h.register <- s
	go s.writePump()
	s.readPump()
}

func Send(channel string, msg Message) {
	h.broadcast <- channelMessage{
		msg:     msg,
		channel: channel,
	}
}
