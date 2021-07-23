package msgbus

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func SetupAuditor(ctx context.Context, bus *MsgBus) {
	messages := bus.Subscribe(ctx, "*")
	HandleMessages(messages)
}

func HandleMessages(messages <-chan *message.Message) {
	for msg := range messages {
		log.Warnf("Auditor got message: %s, event: %s", msg.UUID, Decode(msg.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
