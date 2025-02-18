package sms

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

// Message defines a generic SMS message struct.
type SMSMessage struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

// twillio client
type TwilioClient struct {
	AccountSid string
	AuthToken  string
	From       string
}

// Send implements [sms.SMS] interface.

type SMS interface {
	Send(m *SMSMessage) error
}

// Send sends an SMS message using the Twilio API.

func (c *TwilioClient) Send(m *SMSMessage) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: c.AccountSid,
		Password: c.AuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(m.To)
	params.SetFrom(c.From)
	params.SetBody(m.Message)

	_, err := client.Api.CreateMessage(params)

	if err != nil {
		return err
	}

	return nil

}
