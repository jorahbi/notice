package svc

import (
	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/received"
)

type ServiceContext struct {
	Config  conf.Config
	ReveGpt map[string]received.EventInterface
}

func NewServiceContext(c conf.Config) *ServiceContext {
	reve := make(map[string]received.EventInterface)
	reve[received.RECE_KEY_GPT] = received.NewGpt()
	return &ServiceContext{
		Config:  c,
		ReveGpt: reve,
	}
}
