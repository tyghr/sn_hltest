package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr      = "127.0.0.1:8080"
	userToken = "lady_gaga"
	wsSP      = "SPTI"
	chatPath  = "/wschat"
	chatID    = "room123"
)

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// 	rand.Intn(1000)
// }

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: addr, Path: chatPath, RawQuery: fmt.Sprintf("chat_id=%s", chatID)}
	log.Printf("connecting to %s", u.String())

	dialer := websocket.DefaultDialer
	dialer.Subprotocols = []string{wsSP, userToken}
	c, resp, err := dialer.Dial(u.String(), nil)
	if err != nil {
		if resp != nil {
			b, _ := io.ReadAll(resp.Body)
			log.Println(resp.StatusCode, string(b))
		}
		log.Fatal("dial:", err)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(userToken))
	if err != nil {
		log.Println("write token:", err)
		return
	}

	done := make(chan struct{})
	msgs := make(chan string)

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	go func() {
		defer close(done)
		defer close(msgs)
		for i := 0; i < 10_000; i++ {
			msgs <- fmt.Sprintf("msg %d", time.Now().Nanosecond())
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		case msg := <-msgs:
			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}
