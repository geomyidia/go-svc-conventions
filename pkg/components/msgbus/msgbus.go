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
	callback  func(*message.Message) error
}

func AddHandler(topic string, callback func(*message.Message) error) Handler {
	return Handler{
		handlerID: watermill.NewUUID(),
		topic:     topic,
		callback:  callback,
	}
}

func NewMsgBus() *MsgBus {
	gob.Register(Event{})
	logger := watermill.NewStdLogger(true, false)
	bus := gochannel.NewGoChannel(gochannel.Config{}, logger)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.Recoverer,
		middleware.CorrelationID,
	)
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

func (m *MsgBus) AddHandlers(handlers []Handler) {
	for _, h := range handlers {
		m.AddTopic(h.topic)
		m.router.AddNoPublisherHandler(
			h.handlerID,
			h.topic,
			m.bus,
			h.callback,
		)
	}
}

func (m *MsgBus) Subscribe(ctx context.Context, eventName string) <-chan *message.Message {
	m.AddTopic(eventName)
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

func (m *MsgBus) AddTopic(topic string) {
	m.topics.Lock()
	m.topics.topicTracker[topic] = true
	m.topics.Unlock()
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

func (m *MsgBus) Serve(ctx context.Context) {
	log.Infof("Message bus waiting for events ...")
	if err := m.router.Run(ctx); err != nil {
		log.Panic(err)
	}
}
