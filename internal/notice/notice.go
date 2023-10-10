package notice

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RdsConf struct {
	Addr     string
	Password string
	PoolSize int
}

type Config struct {
	RdsConf redis.RedisConf
}

type ServiceContext struct {
	Config Config
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func NewAsynqServer(c redis.RedisConf) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     c.Host,
			Password: c.Pass,
			PoolSize: 10,
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
