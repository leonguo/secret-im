package util

import (
	"time"
	"github.com/gorilla/websocket"
	"log"
	"github.com/golang/protobuf/proto"
	"../app/models"
	"encoding/json"
	"encoding/base64"
	"../app/lib"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 1024
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
	Id   string
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			responseMessage := new(models.WebSocketMessage)
			proto.Unmarshal(message, responseMessage)
			log.Printf(">>> id %v >>>> responseMessage >>>> %v", c.Id, responseMessage)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("::>>> %v", err)
			}
			break
		}
		//c.Hub.Broadcast <- message
		messageContent := new(models.OutMessage)
		err = proto.Unmarshal(message, messageContent)
		if err != nil {
			break
		}
		log.Printf(">>>id  %v >>>>>> request messageContent >>>> %v", c.Id, messageContent)
		// 处理加密消息
		if messageContent.GetType() == models.OutMessage_Type(models.OutMessage_CIPHERTEXT) {
			source := []byte(messageContent.GetSource())
			requestMessage := new(models.WebSocketRequestMessage)
			err = proto.Unmarshal(source, requestMessage)
			if err != nil {
				break
			}
			log.Printf(">>>id  %v >>>>>> request messageContent >>>> %v", c.Id, requestMessage)
			// send keepalive
			if requestMessage.GetVerb() == "GET" && requestMessage.GetPath() == "/v1/keepalive" {
				webSocketMessage := new(models.WebSocketMessage)
				webSocketMessage.Type = models.WebSocketMessage_Type(models.WebSocketMessage_RESPONSE).Enum()
				responseMessage := new(models.WebSocketResponseMessage)
				responseMessage.Id = requestMessage.Id
				responseMessage.Status = proto.Uint32(200)
				responseMessage.Message = proto.String("OK")
				webSocketMessage.Response = responseMessage
				response, _ := proto.Marshal(webSocketMessage)
				c.Send <- response
			} else if requestMessage.GetVerb() == "PUT" {
				// send message
				inMessages := new(models.IncomingMessageList)
				err := json.Unmarshal(requestMessage.Body, inMessages)
				if err != nil {
					break
				}
				source, err := lib.GetAccountInfoByNumber(c.Id)

				if inMessages.Timestamp == 0 {
					inMessages.Timestamp = time.Now().UnixNano() / 1e6
				}
				for _, message := range inMessages.Messages {
					content, _ := base64.StdEncoding.DecodeString(message.Content)
					outMessage := &models.OutMessage{}
					outMessage.Timestamp = proto.Int64(inMessages.Timestamp)
					outMessage.SourceDevice = proto.Int64(source.Devices[0].Id)
					outMessage.Source = proto.String(source.Number)
					outMessage.Content = content
					outMessage.Type = models.OutMessage_Type(int32(message.Type)).Enum()
					out, err := proto.Marshal(outMessage)
					if err != nil {
						log.Println("Failed to encode message:", err)
					}
					singleMsg := new(SingleMessage)
					singleMsg.Id = inMessages.Destination
					destinationAccount, _ := lib.GetAccountInfoByNumber(inMessages.Destination)
					destinationDevice := destinationAccount.Devices[0]
					singleMsgMessage, _ := AESCBCEncrypt(out, destinationDevice.SignalingKey)

					webSocketMessage := new(models.WebSocketMessage)
					webSocketMessage.Type = models.WebSocketMessage_Type(models.WebSocketMessage_REQUEST).Enum()
					requestMessage := new(models.WebSocketRequestMessage)
					requestMessage.Verb = proto.String("PUT")
					requestMessage.Path = proto.String("/api/v1/message")
					requestMessage.Body = singleMsgMessage
					webSocketMessage.Request = requestMessage
					request, _ := proto.Marshal(webSocketMessage)
					singleMsg.Message = request
					c.Hub.SendSingle <- singleMsg
				}
				//response
				wsResponseMessage := new(models.WebSocketMessage)
				wsResponseMessage.Type = models.WebSocketMessage_Type(models.WebSocketMessage_RESPONSE).Enum()
				wsResponse := new(models.WebSocketResponseMessage)
				wsResponse.Id = requestMessage.Id
				wsResponse.Status = proto.Uint32(200)
				wsResponseMessage.Response = wsResponse
				wsResponseContent, _ := proto.Marshal(wsResponseMessage)
				c.Send <- wsResponseContent
			}
		}
	}
}

type SingleMessage struct {
	Id      string
	Message []byte
}

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	// Unregister requests from clients.
	SendSingle chan *SingleMessage
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		SendSingle: make(chan *SingleMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		case singleMsg := <-h.SendSingle:
			// todo 优化策略
			for client := range h.Clients {
				if client.Id == singleMsg.Id {
					log.Printf(">> send message %v %v %v ", client.Id, len(singleMsg.Message), singleMsg.Message)
					client.Send <- singleMsg.Message
					break
				}
			}
		}
	}
}
