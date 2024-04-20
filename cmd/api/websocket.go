package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is the connected user.
type Client struct {
	Username string
	Socket   *websocket.Conn
	Join     chan bool
	Leave    chan bool
	Message  chan SocketMessage
}

// SocketMessage is the message send through websockets and not through database.
type SocketMessage struct {
	Sender    string `json:"-"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

// clients is map with online clients
var clients = make(map[string]*Client)

// read decodes new incoming message in a SocketMessage struct and then send it to the client's Message channel.
func (c *Client) read() {
	for {
		msg := SocketMessage{}
		msg.Sender = c.Username
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			break
		}
		if err = c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
			break
		}
		c.Message <- msg
	}
	c.Leave <- true
	c.Socket.Close()
}

// write encodes a SocketMessage in JSON.
func (c *Client) write(msg SocketMessage) {
	err := c.Socket.WriteJSON(msg)
	if err != nil {
		fmt.Println(err)
		c.Socket.Close()
	}
}

// AddClient adds client to clients map.
func addClient(c *Client) {
	clients[c.Username] = c
}

// FindAndSend first if the SocketMessage recipient is in clients map. Then writes to recipient socket the SocketMessage.
func FindAndSend(msg SocketMessage) {
	if val, ok := clients[msg.Recipient]; ok {
		if err := val.Socket.WriteJSON(msg); err != nil {
			val.Socket.Close()
		}

	}
}

// NewClient returns a new Client.
func NewClient(conn *websocket.Conn, username string) *Client {
	return &Client{
		Username: username,
		Socket:   conn,
		Join:     make(chan bool),
		Leave:    make(chan bool),
		Message:  make(chan SocketMessage),
	}
}

func (c *Client) listen() {
	for {
		select {
		case <-c.Join:
			addClient(c)
		case <-c.Leave:
			delete(clients, c.Username)
		case msg := <-c.Message:
			// c.write(msg)
			FindAndSend(msg)
		}
		fmt.Println(clients)
	}
}

func (app *application) handleConnections(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	c := NewClient(conn, username)
	go c.listen()
	c.Join <- true
	go c.read()
}
