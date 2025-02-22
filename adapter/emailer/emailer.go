package emailer

import "errors"

var ErrNotImplemented = errors.New("this method is not yet implemented by the invoking object")

type Emailer interface {
	Send(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error

	SenderEmail() string
}

// UnImplemented struct (Mock Implementation)
type UnImplemented struct{}

var _ Emailer = (*UnImplemented)(nil)

// Implement all required methods
func (unImpl *UnImplemented) Send(
	_ string,
	_ string,
	_ []string,
	_ []string,
	_ []string,
	_ []string,
) error {

	return ErrNotImplemented
}

func (unImpl *UnImplemented) SenderEmail() string {
	return ""
}
