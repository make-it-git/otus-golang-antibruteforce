package contract

import "errors"

var ErrDeclined = errors.New("declined")

type LeakyBucket interface {
	Try(login, password, ip string) error
}
