package messaging

import (
	"log"

	"github.com/twilio/twilio-go"
	twilioclient "github.com/twilio/twilio-go/client"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type (
	TwilioOptions struct {
		AccountSID string
		AuthToken  string
		From       string
	}

	Messenger interface {
		Send(text string, to string)
	}

	sms struct {
		config *TwilioOptions
	}
)

func NewSMSMessenger(options *TwilioOptions) Messenger {
	return &sms{options}
}

func (m *sms) Send(text string, to string) {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: m.config.AccountSID,
		Password: m.config.AuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(m.config.From)
	params.SetBody(text)

	if _, err := client.Api.CreateMessage(params); err != nil {
		twilioError := err.(*twilioclient.TwilioRestError)
		println(twilioError.Error())
		return
	}

	log.Printf("OTP sent to: %s", to)
}
