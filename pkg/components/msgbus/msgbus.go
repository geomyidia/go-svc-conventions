package msgbus

import (
	"context"
	"encoding/gob"
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	log "github.com/sirupsen/logrus"
)

const (
	WildCardTopic string = "*"
)

type Topics struct {
	sync.RWMutex
	// XXX Change this bool to struct{} and update the logic below
	topicTracker map[string]bool
}
type MsgBus struct {
	bus    *gochannel.GoChannel
	router *message.Router
	topics *Topics
}

type Handler struct {
	handlerID string
	topic     string
	callback  func()
}

func NewMsgBus() *MsgBus {
	gob.Register(Event{})
	logger := watermill.NewStdLogger(true, false)
	bus := gochannel.NewGoChannel(gochannel.Config{}, logger)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(middleware.Recoverer)
	if err != nil {
		log.Panic(err)
	}
	return &MsgBus{
		bus:    bus,
		router: router,
		topics: &Topics{
			topicTracker: make(map[string]bool),
		},
	}
}

func (m *MsgBus) Subscribe(ctx context.Context, eventName string) <-chan *message.Message {
	m.topics.Lock()
	m.topics.topicTracker[eventName] = true
	m.topics.Unlock()
	messages, err := m.bus.Subscribe(ctx, eventName)
	if err != nil {
		log.Errorf("Couldn't subscribe to %s: %+v", eventName, err)
	}
	return messages
}

func (m *MsgBus) Topics() []string {
	m.topics.RLock()
	var topics []string
	for k := range m.topics.topicTracker {
		topics = append(topics, k)
	}
	m.topics.RUnlock()
	return topics
}

func (m *MsgBus) HasTopic(topic string) bool {
	return m.topics.topicTracker[topic]
}

func (m *MsgBus) Publish(event *Event) {
	topic := event.Name
	event.ID = watermill.NewUUID()
	msg := message.NewMessage(event.ID, Encode(event))
	if m.HasTopic(topic) {
		log.Debugf("Publishing to topic '%s' ...", topic)
		m.bus.Publish(topic, msg)
	}
	if m.HasTopic(WildCardTopic) {
		log.Debugf("Publishing to topic '%s' ...", WildCardTopic)
		m.bus.Publish(WildCardTopic, msg)
	}
}

func (m *MsgBus) Process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Debugf("Received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
