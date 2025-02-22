package messenger

import (
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioService struct {
	cfg    *Config
	client *twilio.RestClient
}

var _ Service = (*TwilioService)(nil)

func (twil *TwilioService) Load(cfg *Config) {
	twil.cfg = cfg
	twil.client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.AccountSid,
		Password: cfg.AuthToken,
	})
}

func (twil *TwilioService) Send(param Param) (bool, error) {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(param.To)
	params.SetFrom(param.From)
	params.SetBody(string(param.Message))

	_, err := twil.client.Api.CreateMessage(params)

	return err == nil, err
}
