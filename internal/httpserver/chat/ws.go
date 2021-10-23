package chat

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		Subprotocols: []string{wsSecureProtocolType},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (chat *Chat) checkToken(d []byte) (string, error) {
	var userToken struct {
		Username  string `json:"username"`
		RequestID string `json:"request_id"`
	}

	data, err := base64.URLEncoding.DecodeString(string(d))
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &userToken)
	if err != nil {
		return "", err
	}

	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}
	defer httpClient.CloseIdleConnections()

	req, err := http.NewRequest(http.MethodPost, chat.conf.SessionValidateUrl, bytes.NewBuffer(d))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status is not ok")
	}

	decoder := json.NewDecoder(resp.Body)
	var tokenValid bool
	err = decoder.Decode(&tokenValid)
	if err != nil {
		return "", err
	}
	if !tokenValid {
		return "", fmt.Errorf("token is not valid")
	}

	return userToken.Username, nil
}

func (chat *Chat) getUserID(c *websocket.Conn, protocols []string) (string, error) {
	var userTokenData []byte

	// log.Print("ws client protocols:", strings.Join(protocols, " "))
	// log.Print("selected protocol:", c.Subprotocol())
	if len(protocols) > 1 {
		userTokenData = []byte(protocols[1])
	}

	// reading first message, which is token
	_, firstMsgToken, err := c.ReadMessage()
	if err != nil {
		u, errToken := chat.checkToken(userTokenData)
		if errToken != nil {
			return "", errToken
		}
		return u, fmt.Errorf("first_message_token: %w", err)
	}
	// this check is not neccessary, actually
	// we check token from first message and from protocol list
	if string(firstMsgToken) != "" && string(firstMsgToken) != string(userTokenData) {
		userTokenData = firstMsgToken
		chat.logger.Debugw("tokens differ")
	}

	u, err := chat.checkToken(userTokenData)
	if err != nil {
		return "", err
	}
	return u, nil
}

func (chat *Chat) readingMessages(c *websocket.Conn, chatID, userToken string) func() error {
	// read messages from all users in this chat
	// getMsgs returns <channels of new msgs(pubsub)>, <array of old messages(streams)>, <func to close first channel>
	ch, oldMsgs, chClose := chat.store.GetMsgs(context.TODO(), chatID)

	go func() {
		// at first, we fill the window by old messages
		for _, msg := range oldMsgs {
			//time: msg.ID

			m := struct {
				User string `json:"user"`
				Text string `json:"text"`
			}{
				User: msg.Value("user").(string),
				Text: msg.Value("text").(string),
			}
			b, _ := json.Marshal(m)

			chat.logger.Debugw("read from redis oldmsg",
				"user", m.User,
				"text,", m.Text,
			)

			err := c.WriteMessage(websocket.TextMessage, []byte(b))
			if err != nil {
				chat.logger.Errorw("write oldmsg", "error", err)
				break
			}
			chat.logger.Debugw("send to ws oldmsg",
				"user", m.User,
				"text,", m.Text,
			)
		}
		chat.logger.Debugw("done reading oldmsgs",
			"chatID", chatID,
			"userToken,", userToken,
		)

		// next, we reading channel of online messages
		for msg := range ch {
			m := struct {
				User string `json:"user"`
				Text string `json:"text"`
			}{}
			b := []byte(msg.Message())
			_ = json.Unmarshal(b, &m)

			chat.logger.Debugw("read from redis online",
				"user", m.User,
				"text,", m.Text,
			)

			err := c.WriteMessage(websocket.TextMessage, []byte(b))
			if err != nil {
				chat.logger.Errorw("write online", "error", err)
				break
			}
			chat.logger.Debugw("send to ws online",
				"user", m.User,
				"text,", m.Text,
			)
		}
		chat.logger.Debugw("stop watching redis",
			"chatID", chatID,
			"userToken,", userToken,
		)
	}()

	return chClose
}

func (chat *Chat) writingMessages(c *websocket.Conn, chatID, userToken string) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			chat.logger.Errorw("read", "error", err)
			break
		}
		chat.logger.Debugw("recv from ws",
			"chat_id", chatID,
			"user_token", userToken,
			"message", message,
		)

		err = chat.store.WriteMessage(
			context.TODO(),
			chatID,
			userToken,
			string(message),
		)
		if err != nil {
			chat.logger.Errorw("store", "error", err)
			break
		}

		chat.logger.Debugw("write to redis",
			"chat_id", chatID,
			"user_token", userToken,
			"message", message,
		)
	}
}

func (chat *Chat) handleChatWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		chat.logger.Errorw("upgrade", "error", err)
		return
	}
	defer c.Close()

	chatID := r.URL.Query().Get("chat_id")
	if chatID == "" {
		http.Error(w, "empty chat_id", http.StatusBadRequest)
		return
	}

	userToken, err := chat.getUserID(c, websocket.Subprotocols(r))
	if err != nil {
		chat.logger.Errorw("getUserToken", "error", err)
	}

	// debug info
	chat.logger.Debugw("channel_info",
		"chat_id", chatID,
		"user_token", userToken,
	)

	// reading msgs part. readingMessages is non-blocking func
	chClose := chat.readingMessages(c, chatID, userToken)
	defer func() {
		_ = chClose()
	}()

	// write msgs part. writingMessages is blocking func
	chat.writingMessages(c, chatID, userToken)

	chat.logger.Debugw("stop watching ws",
		"chat_id", chatID,
		"user_token", userToken,
	)
}
