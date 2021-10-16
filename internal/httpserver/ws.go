package httpserver

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		Subprotocols: []string{wsSecureProtocolType},
	}
)

func (s *Server) getUserToken(c *websocket.Conn, protocols []string) (string, error) {
	var userToken string

	if len(protocols) > 1 {
		userToken = protocols[1]
	}

	// reading first message, which is token
	_, firstMsgToken, err := c.ReadMessage()
	if err != nil {
		return userToken, fmt.Errorf("first_message_token: %w", err)
	}
	// this check is not neccessary
	// actually, we checking token from first message and from protocol list
	if string(firstMsgToken) != "" && string(firstMsgToken) != userToken {
		userToken = string(firstMsgToken)
		s.logger.Debugw("tokens differ",
			"userToken", userToken,
			"firstMsgToken", string(firstMsgToken),
		)
	}

	return userToken, nil
}
