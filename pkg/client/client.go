package client

import (
	"github.com/hibiken/asynq"
	"github.com/jorahbi/notice/internal/aqueue/jobtype"
	"github.com/zeromicro/go-zero/core/jsonx"
)

const (
	WECHAT = "wechat"
)

type RdsConf struct {
	Addr     string
	Password string
	PoolSize int
}

type Payload struct {
	Fo   string
	Data any
	Type string
}

func (p Payload) String() string {
	return p.Data.(string)
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

func (c Client) Send(payload *Payload) (*asynq.TaskInfo, error) {
	return c.ReveSend(jobtype.JOB_KEY_WECHAT_NOTICE, payload)
}

func (c Client) ReveSend(jobName string, payload *Payload) (*asynq.TaskInfo, error) {
	var data []byte
	var err error

	if data, err = jsonx.Marshal(payload); err != nil {
		return nil, err
	}
	return c.aqueue.Enqueue(asynq.NewTask(jobName, data))
}

func (c Client) Close() error {
	return c.aqueue.Close()
}

func (c Client) Native() *asynq.Client {
	return c.aqueue
}
