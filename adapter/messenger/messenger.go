package messenger

import "errors"

var ErrNotImplemented = errors.New("this method is not yet implemented by the invoking object")

type Config struct {
	AccountSid string
	AuthToken  string
}
type Param struct {
	From    string
	To      string
	Message []byte
}

type Service interface {
	Load(*Config)
	Send(Param) (bool, error)
}

/*
Struct implements all the methods that are defined on the parent struct
It is useful to embed this struct in the struct that implements the parent struct
to act as a default method in case you want to implement only some of the methods and not all
*/
type UnImplemented struct{}

var _ Service = (*UnImplemented)(nil)

func (unImpl UnImplemented) Load(_ *Config) {}

func (unImpl UnImplemented) Send(Param) (bool, error) {
	return false, ErrNotImplemented
}

// ************************************************** //
