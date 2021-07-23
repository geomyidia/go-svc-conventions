package msgbus

import (
	"bytes"
	"context"
	"encoding/gob"
	"sync"
	"time"

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

// Event struct is used in two places:
// 1. Subscribing: an event's name is used as the underlying topic in EventBus
// 2. Publishing: an even't name is needed for the topic to publish to, and the
//                event data is what get's published on the EventBus under that
//                topic
type Event struct {
	ID   string
	Name string
	Data string // XXX interface{}? serialized data?
}

type ProcessedEvent struct {
	ProcessedID string
	Time        time.Time
}

type Handler struct {
	handlerID string
	topic     string
	callback  func()
}

func Encode(e *Event) []byte {
	log.Tracef("Attempting to encode %+v ...", e)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		log.Errorf("Possble issue encoding 'event', err: %#v", err)
	}
	log.Trace("Encoded result: ", buf.String())
	return buf.Bytes()
}

func Decode(data []byte) *Event {
	log.Tracef("Attempting to decode bytes: %+v ...", data)
	decoded := &Event{}
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(decoded)
	if err != nil {
		log.Error("Couldn't decode event: ", err)
	}
	log.Tracef("Decoded event: %+v", decoded)
	return decoded
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

func NewEvent(name string, data string) *Event {
	return &Event{
		Name: name,
		Data: data,
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
