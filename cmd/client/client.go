package main

import (
	"context"
	"flag"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/riverchu/pkg/log"
)

var serverAddr = flag.String("addr", "localhost:7749", "http service address")
var interrupt = make(chan os.Signal, 1)

func main() {
	ctx := context.Background()

	flag.Parse()

	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *serverAddr, Path: "/ws"}
	log.Info("connecting to %s", u.String())

	c, _, err := ConnectWebsocket(ctx, u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	Write(c, []byte("ping"))

	Communicate(c)
}

// ConnectWebsocket connect websocket
func ConnectWebsocket(ctx context.Context, url string, header http.Header) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.DialContext(ctx, url, header)
}

// Communicate ...
func Communicate(c *websocket.Conn) {
	msgCh := Read(c)

	go func() {
		for t := range time.Tick(time.Second) {
			Write(c, []byte(t.String()))
		}
	}()

	for {
		select {
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			log.Info("recv: %s", msg)
		case <-interrupt:
			log.Info("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Info("write close:", err)
				return
			}
			select {
			case <-msgCh:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func Read(c *websocket.Conn) <-chan []byte {
	msg := make(chan []byte, 64)
	go func() {
		defer close(msg)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Info("read:", err)
				return
			}
			msg <- message
		}
	}()
	return msg
}

func Write(c *websocket.Conn, msg []byte) error {
	return c.WriteMessage(websocket.TextMessage, msg)
}
