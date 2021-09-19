package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		Subprotocols: []string{wsSecureProtocolType},
	}
)

func getUserToken(c *websocket.Conn, protocols []string) (string, error) {
	var userToken string

	// log.Print("ws client protocols:", strings.Join(protocols, " "))
	// log.Print("selected protocol:", c.Subprotocol())
	if len(protocols) > 1 {
		userToken = protocols[1]
	}

	// reading first message, which is token
	_, firstMsgToken, err := c.ReadMessage()
	if err != nil {
		return userToken, fmt.Errorf("first_message_token: %w", err)
	}
	// this check is not neccessary, actually
	// we check token from first message and from protocol list
	if string(firstMsgToken) != "" && string(firstMsgToken) != userToken {
		userToken = string(firstMsgToken)
		log.Print("tokens differ")
	}

	return userToken, nil
}

func (chat *Chat) readingMessages(c *websocket.Conn, chatID, userToken string) func() error {
	// read messages from all users in this chat
	// getMsgs returns <channels of new msgs(pubsub)>, <array of old messages(streams)>, <func to close first channel>
	ch, oldMsgs, chClose := chat.store.GetMsgs(context.TODO(), chatID)

	go func() {
		// at first, we fill the window by old messages
		for _, msg := range oldMsgs {
			//time: msg.ID
			msgUser := msg.Value("user")
			msgText := msg.Value("text")

			log.Printf("read from redis oldmsg: %s: %s", msgUser, msgText)

			err := c.WriteMessage(websocket.TextMessage, []byte(
				fmt.Sprintf("FROM %s: %s", msgUser, msgText),
			))
			if err != nil {
				log.Println("write oldmsg:", err)
				break
			}
			log.Printf("send to ws oldmsg: %s: %s", msgUser, msgText)
		}
		log.Printf("done reading oldmsgs: %s: %s", chatID, userToken)

		// next, we reading channel of online messages
		for msg := range ch {
			m := struct {
				User string `json:"user"`
				Text string `json:"text"`
			}{}
			_ = json.Unmarshal([]byte(msg.Message()), &m)

			log.Printf("read from redis online: %s: %s", m.User, m.Text)

			err := c.WriteMessage(websocket.TextMessage, []byte(
				fmt.Sprintf("FROM %s: %s", m.User, m.Text),
			))
			if err != nil {
				log.Println("write online:", err)
				break
			}
			log.Printf("send to ws online: %s: %s", m.User, m.Text)
		}
		log.Printf("stop watching redis: %s: %s", chatID, userToken)
	}()

	return chClose
}

func (chat *Chat) writingMessages(c *websocket.Conn, chatID, userToken string) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv from ws: %s: %s: %s", chatID, userToken, message)

		err = chat.store.WriteMessage(
			context.TODO(),
			chatID,
			userToken,
			string(message),
		)
		if err != nil {
			log.Println("store:", err)
			break
		}

		log.Printf("write to redis: %s: %s: %s", chatID, userToken, message)
	}
}

func (chat *Chat) handleChatWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "empty chat_id", http.StatusBadRequest)
		return
	}

	userToken, err := getUserToken(c, websocket.Subprotocols(r))
	if err != nil {
		log.Println("getUserToken:", err)
	}

	// debug info
	log.Printf("channel_info: %s: %s", chatID, userToken)

	// reading msgs part. readingMessages is non-blocking func
	chClose := chat.readingMessages(c, chatID, userToken)
	defer func() {
		_ = chClose()
	}()

	// write msgs part. writingMessages is blocking func
	chat.writingMessages(c, chatID, userToken)

	log.Printf("stop watching ws: %s: %s", chatID, userToken)
}
