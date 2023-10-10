package client

import (
	"github.com/jorahbi/notice/internal/aqueue/jobtype"
	"github.com/jorahbi/notice/internal/notice"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/jsonx"
)

type Client struct {
	*asynq.Client
}

func NewAsynqClient(c notice.RdsConf) *Client {
	return &Client{
		asynq.NewClient(asynq.RedisClientOpt{
			Addr:     c.Addr,
			Password: c.Password,
			PoolSize: c.PoolSize,
		}),
	}
}

func (c Client) Send(payload any) (*asynq.TaskInfo, error) {
	var data []byte
	var err error

	if data, err = jsonx.Marshal(payload); err != nil {
		return nil, err
	}
	return c.Enqueue(asynq.NewTask(jobtype.JOB_KEY_WECHAT_NOTICE, data))
}
