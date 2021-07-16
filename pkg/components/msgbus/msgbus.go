package msgbus

import (
	"sync"

	"github.com/asaskevich/EventBus"
)

const (
	WildCardTopic string = "*"
)

type Topics struct {
	sync.RWMutex
	topicTracker map[string]bool
}
type MsgBus struct {
	bus    EventBus.Bus
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

func NewMsgBus() *MsgBus {
	return &MsgBus{
		bus: EventBus.New(),
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

func (m *MsgBus) Subscribe(eventName string, fn interface{}) {
	m.topics.Lock()
	m.topics.topicTracker[eventName] = true
	m.topics.Unlock()
	m.bus.Subscribe(eventName, fn)
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

func (m *MsgBus) Publish(event *Event) {
	topic := event.Name
	if topic == WildCardTopic {
		for _, t := range m.Topics() {
			m.bus.Publish(t, event)
		}
	} else {
		m.bus.Publish(topic, event)
	}
}
