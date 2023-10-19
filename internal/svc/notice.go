package svc

import (
	"fmt"

	"github.com/jorahbi/notice/internal/conf"
)

type ServiceContext struct {
	Config conf.Config
}

func NewServiceContext(c conf.Config) *ServiceContext {
	c.GPT.Keywords = fmt.Sprintf("@gpt%v", string(rune(8197)))
	return &ServiceContext{
		Config: c,
	}
}
