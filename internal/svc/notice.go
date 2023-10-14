package svc

import (
	"github.com/jorahbi/notice/pkg/client"
)

type Config struct {
	RdsConf client.RdsConf
}

type ServiceContext struct {
	Config Config
	Client *client.Client
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Client: client.NewAsynqClient(c.RdsConf),
	}
}
