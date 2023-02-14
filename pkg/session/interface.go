package session

import "errors"

var UnauthorizedTokenErr = errors.New("unauthorized token")

type Interface interface {
	Store(userID uint64) (string, error)
	Authorize(token string) (uint64, error)
	Delete(token string) error
}
