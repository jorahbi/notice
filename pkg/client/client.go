package client

import (
	"github.com/jorahbi/notice/internal/aqueue/jobtype"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/jsonx"
)

type RdsConf struct {
	Addr     string
	Password string
	PoolSize int
}

type Client struct {
	aqueue *asynq.Client
}

func NewAsynqClient(c RdsConf) *Client {
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
	return c.aqueue.Enqueue(asynq.NewTask(jobtype.JOB_KEY_WECHAT_NOTICE, data))
}

func (c Client) Close() error {
	return c.aqueue.Close()
}

func (c Client) Native() *asynq.Client {
	return c.aqueue
}
