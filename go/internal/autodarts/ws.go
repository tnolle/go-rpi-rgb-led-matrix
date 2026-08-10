package autodarts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/tnolle/go-rpi-rgb-led-matrix/internal/keycloak"
)

type AutodartsWSClient struct {
	URL  string
	conn *websocket.Conn

	mu          sync.Mutex
	subscribers map[string]chan any
}

func NewAutodartsWSClient() *AutodartsWSClient {
	return &AutodartsWSClient{
		URL:         "wss://api.autodarts.io/ms/v0/subscribe",
		subscribers: make(map[string]chan any),
	}
}

func (c *AutodartsWSClient) Connect() error {
	var err error
	h := make(http.Header)
	token, err := keycloak.AccessToken()
	h.Set("Authorization", "Bearer "+token)
	c.conn, _, err = websocket.DefaultDialer.Dial(c.URL, h)
	if err != nil {
		return err
	}
	go c.readPump()
	return nil
}

type SubscribeMessage struct {
	Method  string `json:"type"`
	Channel string `json:"channel"`
	Topic   string `json:"topic"`
}

type Message struct {
	Channel string `json:"channel"`
	Topic   string `json:"topic"`
	Data    any    `json:"data"`
}

func (c *AutodartsWSClient) readPump() {
	defer c.conn.Close()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}
		key := fmt.Sprintf("%s:%s", m.Channel, m.Topic)
		c.mu.Lock()
		ch, ok := c.subscribers[key]
		if ok {
			select {
			case ch <- m.Data:
			default:
			}
		}
		c.mu.Unlock()
	}
}

func (c *AutodartsWSClient) subscribe(channel, topic string) chan any {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.WriteJSON(SubscribeMessage{
		Method:  "subscribe",
		Channel: channel,
		Topic:   topic,
	})
	ch := make(chan any)
	key := fmt.Sprintf("%s:%s", channel, topic)
	c.subscribers[key] = ch
	return ch
}

func (c *AutodartsWSClient) unsubscribe(channel, topic string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.WriteJSON(SubscribeMessage{
		Method:  "unsubscribe",
		Channel: channel,
		Topic:   topic,
	})
	key := fmt.Sprintf("%s:%s", channel, topic)
	if ch, ok := c.subscribers[key]; ok {
		close(ch)
		delete(c.subscribers, key)
	}
}
func handle[T any](ctx context.Context, c *AutodartsWSClient, channel, topic string) <-chan T {
	fmt.Println("Subscribing to channel:", channel, "topic:", topic)
	ch := make(chan T)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Context done, unsubscribing")
				c.unsubscribe(channel, topic)
				close(ch)
				return
			case msg := <-c.subscribe(channel, topic):
				var body T
				b, _ := json.Marshal(msg)
				json.Unmarshal(b, &body)
				ch <- body
			}
		}
	}()
	return ch
}

type OnlineUserMessage struct {
	Online int `json:"online"`
}

func (c *AutodartsWSClient) OnOnlineUsersChange(ctx context.Context) <-chan OnlineUserMessage {
	return handle[OnlineUserMessage](ctx, c, "autodarts.users", "online-users")
}

type MatchesCountMessage struct {
	Count int `json:"count"`
}

func (c *AutodartsWSClient) OnMatchCountChange(ctx context.Context) <-chan MatchesCountMessage {
	return handle[MatchesCountMessage](ctx, c, "autodarts.matches", "matches.count")
}
