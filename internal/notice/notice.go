package notice

import (
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
	wx = &wechat{timer: time.NewTimer(300 * time.Second)}
	notices = map[string]Notice{
		NOTICE_WECHAT: wx,
	}
}

func Wechat() Notice {
	return notices[NOTICE_WECHAT]
}
