package msgbus

import (
	"bytes"
	"encoding/gob"
	"time"

	log "github.com/sirupsen/logrus"
)

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

func NewEvent(name string, data string) *Event {
	return &Event{
		Name: name,
		Data: data,
	}
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
