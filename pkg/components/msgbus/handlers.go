package msgbus

import (
	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func HandleWildCard(msg *message.Message) error {
	log.Warnf("Auditor got event: %+v", Decode(msg.Payload))
	// we need to Acknowledge that we received and processed the message,
	// otherwise, it will be resent over and over again.
	msg.Ack()
	return nil
}

func HandlePing(msg *message.Message) error {
	log.Warnf("Ping listeners got event: %+v", Decode(msg.Payload))
	// we need to Acknowledge that we received and processed the message,
	// otherwise, it will be resent over and over again.
	msg.Ack()
	return nil
}

func HandleHealth(msg *message.Message) error {
	log.Warnf("Health listeners got event: %+v", Decode(msg.Payload))
	// we need to Acknowledge that we received and processed the message,
	// otherwise, it will be resent over and over again.
	msg.Ack()
	return nil
}
