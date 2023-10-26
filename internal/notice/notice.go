package notice

import (
	"github.com/hibiken/asynq"
	"github.com/jorahbi/notice/pkg/client"
)

type NoticeInterface interface {
	asynq.Handler
	Send(payload *client.Payload)
}

const (
	NOTICE_WECHAT = "wechat"
)

// var notices = map[string]NoticeInterface{
// 	NOTICE_WECHAT: &wechat{ch: make(chan struct{})},
// }

type app struct {
	app NoticeInterface
	ch  chan struct{}
}

var notices = map[string]app{
	// NOTICE_WECHAT: app{app: &wechat{ch: make(chan struct{})}},
}

func Wechat() NoticeInterface {
	app := notices[NOTICE_WECHAT]
	return app.app
}

func setNotice(key string, notice NoticeInterface) {
	notices[key] = app{app: notice}
}
