package internal

import (
	"github.com/ilivestrong/otp-service/internal/messaging"
	mq "github.com/ilivestrong/otp-service/internal/rabbitmq"
)

type (
	otpService struct {
		consumer mq.MQClient
		sms      messaging.Messenger
	}
)

func (otpSvc *otpService) StartListening() {

}

func NewOtpService(consumer mq.MQClient, sms messaging.Messenger) *otpService {
	return &otpService{consumer, sms}
}
