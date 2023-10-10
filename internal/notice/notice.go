package notice

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jorahbi/notice/pkg/client"
)

type Config struct {
	RdsConf client.RdsConf
}

type ServiceContext struct {
	Config Config
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func NewAsynqServer(c client.RdsConf) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     c.Addr,
			Password: c.Password,
			PoolSize: c.PoolSize,
		},
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("asynq server exec task IsFailure ======== >>>>>>>>>>> err : %+v  \n", err)
				return true
			},
			Concurrency: 20, //max concurrent process job task nu
		},
	)
}
