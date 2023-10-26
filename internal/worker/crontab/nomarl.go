package crontab

import (
	"context"
	"fmt"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/notice"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
)

type nomarl0800 struct{}

func (c nomarl0800) Start(ctx context.Context, svcCtx *svc.ServiceContext, conf conf.Notice) func() {
	return func() {
		fmt.Println("crontab ...")
		for _, item := range conf.To {
			notice.Wechat().Send(&client.Payload{Fo: item, Data: conf.Tpl})
		}
	}
}
func (c nomarl0800) Stop(ctx context.Context) {

}
