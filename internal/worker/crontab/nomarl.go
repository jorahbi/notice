package crontab

import (
	"context"
	"fmt"
	"time"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/event"
	"github.com/jorahbi/notice/internal/notice"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
)

func morning(ctx context.Context, svcCtx *svc.ServiceContext, fo string, conf conf.Job) func() {
	return func() {
		reve := event.WaitReve()
		reve.Wait(ctx, event.WaitConf{
			Fo:   fo,
			Num:  10,
			Time: 30 * time.Second,
			Callback: func() {
				fmt.Printf("corn time reve %v", conf.Tpl)
				wechat := notice.Wechat()
				wechat.Send(&client.Payload{Fo: fo, Data: conf.Tpl})
			},
		})
	}
}

func noon(ctx context.Context, svcCtx *svc.ServiceContext, fo string, conf conf.Job) func() {
	return func() {
		reve := event.WaitReve()
		reve.Wait(ctx, event.WaitConf{
			Fo:   fo,
			Num:  10,
			Time: 30 * time.Second,
			Callback: func() {
				fmt.Printf("corn time reve %v", conf.Tpl)
				wechat := notice.Wechat()
				wechat.Send(&client.Payload{Fo: fo, Data: conf.Tpl})
			},
		})
	}
}

func night(ctx context.Context, svcCtx *svc.ServiceContext, fo string, conf conf.Job) func() {
	return func() {
		reve := event.WaitReve()
		reve.Wait(ctx, event.WaitConf{
			Fo:   fo,
			Num:  10,
			Time: 30 * time.Second,
			Callback: func() {
				fmt.Printf("corn time reve %v", conf.Tpl)
				wechat := notice.Wechat()
				wechat.Send(&client.Payload{Fo: fo, Data: conf.Tpl})
			},
		})
	}
}

func run(ctx context.Context, svcCtx *svc.ServiceContext, fo string, conf conf.Job) func() {
	return func() {
		reve := event.WaitReve()
		reve.Wait(ctx, event.WaitConf{
			Fo:   fo,
			Num:  10,
			Time: 30 * time.Second,
			Callback: func() {
				fmt.Printf("corn time reve %v", conf.Tpl)
				wechat := notice.Wechat()
				wechat.Send(&client.Payload{Fo: fo, Data: conf.Tpl})
			},
		})
	}
}
