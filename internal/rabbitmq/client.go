package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ilivestrong/otp-service/internal/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	sendotp_exchange_name = "verification"
	exchange_type_topic   = "topic"
	sendotp_queue_name    = "otp_request"

	otpcreated_queue_binding_key = "OtpCreated.*"
	otpcreated_routing_key       = "OtpCreated.newaccount"
	otpcreated_queue_name        = "otps_created"
)

type (
	MQClient interface {
		Consume(ctx context.Context, options *messaging.TwilioOptions)
		Publish(ctx context.Context, info otpSentInfo)
	}

	otpSentInfo struct {
		phone string
		otp   string
	}

	otpMQClient struct {
		ch        *amqp.Channel
		smsClient messaging.Messenger
	}
)

func (otpRPub *otpMQClient) Publish(ctx context.Context, info otpSentInfo) {
	eventPayload := fmt.Sprintf(`{"phone_number":"%s","otp":"%s"}`, info.phone, info.otp)

	err := otpRPub.ch.PublishWithContext(ctx,
		sendotp_exchange_name,
		otpcreated_routing_key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(eventPayload),
		})
	failOnError(err, "Failed to publish a message")
}

func (otpEC *otpMQClient) Consume(ctx context.Context, options *messaging.TwilioOptions) {
	msgs, err := otpEC.ch.Consume(sendotp_queue_name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to consume messages from queue")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received SendOTP request for phone number: %s", d.Body)
			otp := generateOTP()
			otpEC.smsClient.Send(fmt.Sprintf("Here is your 6 digit OTP code: %s", otp), string(d.Body))
			otpEC.Publish(ctx, otpSentInfo{
				phone: string(d.Body),
				otp:   otp,
			})
			log.Println("triggered OtpCreated event.")
		}
	}()
	<-forever
}

func declareExchange(ch *amqp.Channel, name string) {
	err := ch.ExchangeDeclare(name, exchange_type_topic, true, false, false, false, nil)
	failOnError(err, fmt.Sprintf("failed to declare exchange: %s\n", name))
}

func declarePublishQueue(ch *amqp.Channel) amqp.Queue {
	q, err := ch.QueueDeclare(otpcreated_queue_name, false, false, false, false, nil)
	failOnError(err, "failed to declare a queue")
	return q
}

func bindQueueToExchange(q amqp.Queue, exchange string, ch *amqp.Channel) {
	err := ch.QueueBind(q.Name, otpcreated_queue_binding_key, exchange, false, nil)
	failOnError(err, fmt.Sprintf("failed to bind queue to exchange: %s\n", exchange))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NewOtpMQClient(amqpConn *amqp.Connection, smsClient messaging.Messenger) MQClient {
	ch, err := amqpConn.Channel()
	failOnError(err, "failed to create message channel")

	declareExchange(ch, sendotp_exchange_name)
	bindQueueToExchange(declarePublishQueue(ch), sendotp_exchange_name, ch)
	return &otpMQClient{ch, smsClient}
}

func generateOTP() string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}
