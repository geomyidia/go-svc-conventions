package msgbus

import (
	"bytes"
	"context"
	"encoding/gob"
	"sync"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
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
	topics *Topics
}

// Event struct is used in two places:
// 1. Subscribing: an event's name is used as the underlying topic in EventBus
// 2. Publishing: an even't name is needed for the topic to publish to, and the
//                event data is what get's published on the EventBus under that
//                topic
type Event struct {
	Name string
	Data string // XXX interface{}? serialized data?
}

func Encode(e *Event) []byte {
	log.Debugf("Attempting to encode %+v ...", e)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(e)
	if err != nil {
		log.Errorf("Possble issue encoding 'event', err: %#v", err)
	}
	log.Debug("Encoded result: ", buf.String())
	return buf.Bytes()
}

func Decode(data []byte) *Event {
	log.Debugf("Attempting to decode bytes: %+v ...", data)
	buf := bytes.NewBuffer(data)
	decoded := &Event{}
	dec := gob.NewDecoder(buf)
	err := dec.Decode(decoded)
	if err != nil {
		log.Error("Couldn't decode event: ", err)
	}
	log.Debugf("Decoded result: %+v", decoded)
	return decoded
}

func NewMsgBus() *MsgBus {
	gob.Register(Event{})
	bus := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)
	return &MsgBus{
		bus: bus,
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
		log.Error("Couldn't subscribe to %s: %+v", eventName, err)
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
	msg := message.NewMessage(watermill.NewUUID(), Encode(event))
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
