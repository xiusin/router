package interfaces

import (
	"net/http"
	"time"
)

type SessionManagerInf interface {
	Session(*http.Request, http.ResponseWriter) (SessionInf, error)
}

type SessionConfigInf interface {
	GetCookieName() string
	GetExpires() time.Duration
	GetHttpOnly() bool
	GetSecure() bool
}

type SessionStoreInf interface {
	GetConfig() SessionConfigInf
	Read(string, interface{}) error
	Save(string, interface{}) error
	Clear(string) error
}

type SessionInf interface {
	Set(string, interface{}) error
	Get(string) (interface{}, error)
	AddFlush(string, interface{}) error
	Remove(string) error
	Save() error
	Clear() error
}