package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ilivestrong/otp-service/internal/messaging"
	"github.com/ilivestrong/otp-service/internal/rabbitmq"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type (
	Options struct {
		AMQPAddress   string
		TwilioOptions messaging.TwilioOptions
	}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	options := Options{
		AMQPAddress: mustGetEnv("AMQP_ADDRESS"),
		TwilioOptions: messaging.TwilioOptions{
			AccountSID: mustGetEnv("TWILIO_ACCOUNT_SID"),
			AuthToken:  mustGetEnv("TWILIO_AUTH_TOKEN"),
			From:       mustGetEnv("TWILIO_PHONE_NUMBER"),
		},
	}

	smsClient := messaging.NewSMSMessenger(&options.TwilioOptions)
	mqclient := rabbitmq.NewOtpMQClient(bootMQ(&options), smsClient)

	// otpService := internal.NewOtpService(eventConsumer, twilioClient)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	startConsumingEvents(ctx, mqclient, "SendOTP", &options.TwilioOptions)
}

func bootMQ(options *Options) *amqp.Connection {
	conn, err := amqp.Dial(options.AMQPAddress)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ, %v", err)
	}
	return conn
}

func mustGetEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("failed to get env for: %s", key)
	}
	return v
}

func startConsumingEvents(ctx context.Context, mqclient rabbitmq.MQClient, event string, twOptions *messaging.TwilioOptions) {
	log.Printf("starting listening for %s events. To exit press CTRL+C\n", event)
	mqclient.Consume(ctx, twOptions)
}
