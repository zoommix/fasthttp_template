package utils

import (
	"encoding/json"
	"fmt"

	"github.com/fasthttp/websocket"
	"github.com/gofrs/uuid"
)

// Client ...
type Client struct {
	ID         string
	Connection *websocket.Conn
	UserID     int
}

// PubSub ...
type PubSub struct {
	Clients       []Client
	Subscriptions []Subscription
}

// Subscription ...
type Subscription struct {
	Topic  string
	Client *Client
}

// Message ...
type Message struct {
	Action   string          `json:"action"`
	Topic    string          `json:"topic"`
	Message  json.RawMessage `json:"message"`
	Filters  json.RawMessage `json:"filters"`
	Event    json.RawMessage `json:"event"`
	Slot     json.RawMessage `json:"slot"`
	TimeZone *string         `json:"timezone"`
}

// NewClient ...
func NewClient(conn *websocket.Conn) Client {
	return Client{ID: autoID(), Connection: conn}
}

// NewMessage ...
func NewMessage() Message {
	return Message{}
}

func autoID() string {
	v4, _ := uuid.NewV4()
	return uuid.Must(v4, nil).String()
}

// SendJSON ...
func (client *Client) SendJSON(obj interface{}) error {
	return client.Connection.WriteJSON(obj)
}

// SendJSONError ...
func (client *Client) SendJSONError(action, str string) {
	client.SendJSON(map[string]string{"Action": action, "Error": str})
}

// SendJSONStatus ...
func (client *Client) SendJSONStatus(action, str string) {
	client.SendJSON(map[string]string{"Action": action, "Status": str})
}

// Send ...
func (client *Client) Send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}

// AddClient ...
func (ps *PubSub) AddClient(client Client) *PubSub {
	ps.Clients = append(ps.Clients, client)

	LogInfo("Adding new client: UserID=", client.UserID, " to the list. Total clients: ", len(ps.Clients))

	return ps
}

// RemoveClient ...
func (ps *PubSub) RemoveClient(client Client) *PubSub {
	// first remove all subscriptions by this client
	for index, sub := range ps.Subscriptions {
		if client.ID == sub.Client.ID {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	// remove client from the list
	for index, c := range ps.Clients {
		if c.ID == client.ID {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}

	return ps
}

// GetSubscriptions ...
func (ps *PubSub) GetSubscriptions(topic string, client *Client) []Subscription {
	var subscriptionList []Subscription

	for _, subscription := range ps.Subscriptions {

		if client != nil {
			// used for subscribing, checks if user already subscribed
			if subscription.Client.ID == client.ID && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		} else {
			// used for publishing, returns all subs
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}

	return subscriptionList
}

// Subscribe ...
func (ps *PubSub) Subscribe(client *Client, topic string) *PubSub {
	clientSubs := ps.GetSubscriptions(topic, client)

	if len(clientSubs) > 0 {
		// client is subscribed this topic before
		return ps
	}

	newSubscription := Subscription{
		Topic:  topic,
		Client: client,
	}

	ps.Subscriptions = append(ps.Subscriptions, newSubscription)

	return ps
}

// Publish ...
func (ps *PubSub) Publish(topic string, message []byte, excludeClient *Client) {
	subscriptions := ps.GetSubscriptions(topic, nil)

	for _, sub := range subscriptions {
		LogInfo(fmt.Sprintf("Sending to client id %s, message: %s \n", sub.Client.ID, message))
		sub.Client.SendJSON(map[string]string{"Message": string(message)})
	}
}

// Unsubscribe ...
func (ps *PubSub) Unsubscribe(client *Client, topic string) *PubSub {
	//clientSubscriptions := ps.GetSubscriptions(topic, client)
	for index, sub := range ps.Subscriptions {

		if sub.Client.ID == client.ID && sub.Topic == topic {
			// found this subscription from client and we do need remove it
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	return ps
}
