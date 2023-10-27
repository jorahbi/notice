package notice

import (
	"math/rand"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jorahbi/notice/pkg/client"
)

type Notice interface {
	asynq.Handler
	Send(payload *client.Payload)
}

const (
	NOTICE_WECHAT = "wechat"
)

var (
	wx      *wechat
	notices map[string]Notice
)

func init() {
	wx = &wechat{timer: time.NewTimer(time.Duration(rand.Intn(240)+240) * time.Second)}
	notices = map[string]Notice{
		NOTICE_WECHAT: wx,
	}
}

func Wechat() Notice {
	return notices[NOTICE_WECHAT]
}
