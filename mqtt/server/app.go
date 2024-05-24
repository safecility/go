package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/mqtt/lib"
	"github.com/safecility/go/mqtt/messages"
	"time"
)

type GooglePubSub struct {
	Joins            *pubsub.Topic
	Uplinks          *pubsub.Topic
	UplinkErrors     *pubsub.Topic
	Downlinks        *pubsub.Subscription
	DownlinkReceipts *pubsub.Topic
	DownlinkErrors   *pubsub.Topic
}

// MqttProxyConfig - for ttn the Username has the form: username = fmt.Sprintf("%s@ttn", p.AppID)
type MqttProxyConfig struct {
	AppID           string
	AppKey          string
	MqttAddress     string
	Username        string
	CanDownlink     bool
	GooglePubSub    GooglePubSub
	Transformer     lib.PahoTransformer
	PayloadAdjuster lib.PayloadAdjuster
}

type PahoPubSub struct {
	client  *lib.PahoClient
	AppID   string
	DevID   string
	SubPubs map[string]string
	Errors  chan error
}

func (p *PahoPubSub) Subscribe(topic string) (<-chan paho.Message, error) {
	uplink := make(chan paho.Message, mqttBufferSize)

	token := p.client.Subscribe(topic, func(c paho.Client, msg paho.Message) {
		uplink <- msg
	})
	token.Wait()
	err := token.Error()

	return uplink, err
}

func (p *PahoPubSub) SendDownlink(topic string, payload []byte) error {
	token := p.client.Publish(topic, payload)
	token.WaitTimeout(time.Second * 2)
	return token.Error()
}

func NewPahoProxy(p MqttProxyConfig) (*PahoProxy, error) {

	mqttClient := NewClient(p.AppID, p.Username, p.AppKey, p.MqttAddress)

	ps, err := mqttClient.GetPahoPubSub()
	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("could not create PahoPubsub %s", p.AppID))
		return nil, err
	}

	pp := &PahoProxy{
		appID:            p.AppID,
		pahoPubsub:       ps,
		joins:            p.GooglePubSub.Joins,
		uplinks:          p.GooglePubSub.Uplinks,
		uplinkErrors:     p.GooglePubSub.Uplinks,
		downlinkReceipts: p.GooglePubSub.DownlinkReceipts,
		downlinks:        p.GooglePubSub.Downlinks,
		transformer:      p.Transformer,
		adjuster:         p.PayloadAdjuster,
	}

	return pp, nil
}

type PahoProxy struct {
	appID            string
	pahoPubsub       *PahoPubSub
	pahoJoin         <-chan paho.Message
	pahoUplinks      <-chan paho.Message
	pahoUplinkErrors <-chan paho.Message
	pahoDownQueued   <-chan paho.Message
	pahoDownSent     <-chan paho.Message
	pahoDownAck      <-chan paho.Message
	pahoDownNack     <-chan paho.Message
	pahoDownFail     <-chan paho.Message
	pahoDownErrors   <-chan paho.Message
	transformer      lib.PahoTransformer
	adjuster         lib.PayloadAdjuster
	uplinks          *pubsub.Topic
	uplinkErrors     *pubsub.Topic
	joins            *pubsub.Topic
	// not all services need downlinks so nil check before subscribing
	downlinkReceipts *pubsub.Topic
	downlinkErrors   *pubsub.Topic
	downlinks        *pubsub.Subscription
}

// Run TODO make it so we can configure which channels to listen to and forward
func (pp *PahoProxy) Run() error {
	err := pp.initPubsub()
	if err != nil {
		return err
	}

	go pp.listenAndPublishUplinks()
	go pp.listenAndPublishUplinkErrors()
	go pp.listenAndPublishJoins()
	// not all services need downlinks
	if pp.downlinks != nil {
		log.Debug().Str("proxy", pp.appID).Msg("starting downlink listeners")
		go pp.listenToGoogle()
		go pp.listenAndPublishDownlinks(messages.Queued)
		go pp.listenAndPublishDownlinks(messages.Sent)
		go pp.listenAndPublishDownlinks(messages.Ack)
		go pp.listenAndPublishDownlinks(messages.Nack)
		go pp.listenAndPublishDownlinks(messages.Failed)
	}
	return nil
}

func (pp *PahoProxy) initPubsub() error {
	joinTopic := pp.transformer.MqttTopic(messages.AllDevices, messages.Join)
	j, err := pp.pahoPubsub.Subscribe(joinTopic)
	if err != nil {
		log.Err(err).Msg("could not subscribe to join")
		return err
	}
	pp.pahoJoin = j

	uplinkAllTopic := pp.transformer.UplinkTopic(messages.AllDevices)
	u, err := pp.pahoPubsub.Subscribe(uplinkAllTopic)
	if err != nil {
		log.Err(err).Msg("could not subscribe to uplinks")
		return err
	}
	pp.pahoUplinks = u

	uplinkErrorsAllTopic := pp.transformer.UplinkErrorsTopic(messages.AllDevices)
	ue, err := pp.pahoPubsub.Subscribe(uplinkErrorsAllTopic)
	if err != nil {
		log.Err(err).Msg("could not subscribe to uplinks")
		return err
	}
	pp.pahoUplinkErrors = ue

	downlinkQueuedTopic := pp.transformer.MqttTopic(messages.AllDevices, messages.Queued)
	d, err := pp.pahoPubsub.Subscribe(downlinkQueuedTopic)
	if err != nil {
		log.Err(err).Msg("could not subscribe to downlinks")
		return err
	}
	pp.pahoDownQueued = d

	downlinkSentTopic := pp.transformer.MqttTopic(messages.AllDevices, messages.Sent)
	s, err := pp.pahoPubsub.Subscribe(downlinkSentTopic)
	if err != nil {
		log.Err(err).Msg("could not subscribe to downlinks")
		return err
	}
	pp.pahoDownSent = s

	downlinkFailedTopic := pp.transformer.MqttTopic(messages.AllDevices, messages.Failed)
	f, err := pp.pahoPubsub.Subscribe(downlinkFailedTopic)
	if err != nil {
		log.Err(err).Msg("could not subscribe to downlinks")
		return err
	}
	pp.pahoDownFail = f
	return nil
}

func (pp *PahoProxy) listenAndPublishUplinks() {
	for event := range pp.pahoUplinks {
		sm, err := pp.transformer.TransformPahoUplinkMessage(event)
		if err != nil {
			log.Err(err).Msg("transformer err")
		}
		if err := pp.adjuster.AdjustPayload(sm); err != nil {
			log.Err(err).Msg("adjuster err")
		}
		event.Ack()

		err = pp.publishUplink(sm)
		if err != nil {
			log.Err(err).Msg("gPubSub err")
		}
		log.Debug().Str("id", sm.DeviceUID).Str("eui", fmt.Sprintf("%v", sm.DeviceEUI)).
			Str("topic", pp.uplinks.String()).
			Msg(fmt.Sprintf("published sm %+v", sm))
	}
}

func (pp *PahoProxy) listenAndPublishJoins() {
	for event := range pp.pahoJoin {
		sm, err := pp.transformer.TransformPahoJoinMessage(event)
		if err != nil {
			log.Err(err).Msg("transformer join err")
		}
		err = pp.publishJoin(sm)
		if err != nil {
			log.Err(err).Msg("gPubSub err")
		}
		log.Debug().Str("topic", pp.uplinks.String()).Msg("published uplink")
		event.Ack()
	}
}

func (pp *PahoProxy) listenAndPublishUplinkErrors() {
	for event := range pp.pahoUplinkErrors {
		log.Warn().Str("error", fmt.Sprintf("%s", event.Payload())).Msg("uplink error")
		event.Ack()
	}
}

func (pp *PahoProxy) listenAndPublishDownlinks(channel messages.MqttChannel) {
	var downlinks <-chan paho.Message
	switch channel {
	case messages.Queued:
		downlinks = pp.pahoDownQueued
		break
	case messages.Sent:
		downlinks = pp.pahoDownSent
		break
	case messages.Ack:
		downlinks = pp.pahoDownAck
		break
	case messages.Nack:
		downlinks = pp.pahoDownNack
		break
	case messages.Failed:
		downlinks = pp.pahoDownFail
		break
	default:
		downlinks = pp.pahoDownQueued
	}

	for event := range downlinks {
		log.Debug().Str("topic", event.Topic()).Str("channel", string(channel)).
			Msg("downlink receipt")

		sm, err := pp.transformer.TransformPahoDownlinkMessage(event, channel)
		if err != nil {
			log.Err(err).Str("downlink", "err").Msg("transformer err")
		}

		if channel == messages.Ack {
			log.Info().Str("downlink", "ack").Msg("event ack")
		}

		err = pp.publishDownlinkReceipt(sm)
		if err != nil {
			log.Err(err).Msg("could not publish downlink receipt")
		}
		log.Debug().Str("channel", string(channel)).Str("pubsub", "google").
			Msg("published downlink receipt")
		event.Ack()
	}
}

func (pp *PahoProxy) listenToGoogle() {
	log.Debug().Str("subscription", pp.appID).Msg("starting google sub listener")

	err := pp.downlinks.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {

		var byteMessage stream.SimpleMessage

		if err := json.Unmarshal(msg.Data, &byteMessage); err != nil {
			log.Err(fmt.Errorf("could not decode message data: %s", msg.Data)).Msg("decode err")
			msg.Ack()
			return
		}
		log.Debug().Msg(fmt.Sprintf("got downlink byte message %+v", byteMessage))

		if byteMessage.Source != pp.appID {
			log.Warn().Str("appID", fmt.Sprintf("this: %s, message: %s", pp.appID, byteMessage.Source)).
				Msg("appID of message did not match proxy")
		}

		//we don't want downlink queue for a device to have more than one downlink
		deviceTopic := pp.transformer.GetDownlinkTopicReplace(byteMessage.DeviceUID)
		correlationID := fmt.Sprintf("%s-%d", byteMessage.DeviceUID, time.Now().Unix())

		rawPayload, err := pp.transformer.CreateDownlink(byteMessage, []string{correlationID})
		if err != nil {
			log.Err(err).Msg("create downlink err")
			//handle these in the compliance code rather than by sending a nack
			msg.Ack()
			return
		}

		if err := pp.pahoPubsub.SendDownlink(deviceTopic, rawPayload); err != nil {
			//handle these in the compliance code rather than by sending a nack
			log.Err(err).Msg("could not send downlink")
			msg.Ack()
			return
		}
		log.Info().Str("device", byteMessage.DeviceUID).
			Str("payload", fmt.Sprintf("%02x", byteMessage.Payload)).Msg("sent downlink")

		msg.Ack()
	})
	if err != nil {
		log.Err(err).Msg("google pubsub err")
	}
}

func (pp *PahoProxy) publishUplink(message *messages.LoraMessage) error {
	_, err := stream.PublishToTopic(message, pp.uplinks)
	return err
}

func (pp *PahoProxy) publishJoin(message *messages.LoraMessage) error {
	if pp.joins == nil {
		log.Warn().Msg("no topic for joins")
	}
	_, err := stream.PublishToTopic(message, pp.joins)
	return err
}

func (pp *PahoProxy) publishDownlinkReceipt(message *messages.LoraMessage) error {
	if pp.downlinkReceipts == nil {
		log.Warn().Msg("no topic for downlinkReceipts")
	}
	_, err := stream.PublishToTopic(message, pp.downlinkReceipts)
	return err
}
