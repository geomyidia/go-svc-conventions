package msgbus

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	log "github.com/sirupsen/logrus"
)

func SetupAuditor(ctx context.Context, bus *MsgBus) {
	log.Debug("Setting up event bus auditor ...")
	messages := bus.Subscribe(ctx, "*")
	log.Info("Auditor is listening for new events ...")
	HandleMessages(messages)
}

func HandleMessages(messages <-chan *message.Message) {
	for msg := range messages {
		log.Warnf("Auditor got event: %+v", Decode(msg.Payload))
		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}
