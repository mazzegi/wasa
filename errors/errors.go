package errors

import "fmt"

type Error struct {
	msg string
}

func (e Error) Error() string {
	return e.msg
}

func Errorf(s string, args ...interface{}) error {
	return Error{
		msg: fmt.Sprintf(s, args...),
	}
}

func Wrap(err error, msg string) error {
	return Error{
		msg: fmt.Sprintf("%s: %s", msg, err.Error()),
	}
}
