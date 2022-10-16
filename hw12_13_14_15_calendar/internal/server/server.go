package server

import "context"

type Instance interface {
	Start(context.Context) error
	Stop(context.Context) error
}
