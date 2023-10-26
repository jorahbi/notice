package crontab

import (
	"context"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/notice"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
)

func morning(ctx context.Context, svcCtx *svc.ServiceContext, conf conf.Job) func() {
	return func() {
		for _, item := range conf.To {
			notice.Wechat().Send(&client.Payload{Fo: item, Data: conf.Tpl})
		}
	}
}

func eat(ctx context.Context, svcCtx *svc.ServiceContext, conf conf.Job) func() {
	return func() {
		for _, item := range conf.To {
			notice.Wechat().Send(&client.Payload{Fo: item, Data: conf.Tpl})
		}
	}
}

func run(ctx context.Context, svcCtx *svc.ServiceContext, conf conf.Job) func() {
	return func() {
		for _, item := range conf.To {
			notice.Wechat().Send(&client.Payload{Fo: item, Data: conf.Tpl})
		}
	}
}
