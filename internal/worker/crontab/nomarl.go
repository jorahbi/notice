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
		fmt.Printf("corn morning reve")
		reve := event.ReveInstance()
		ch := reve.Get(fo)
		wechat := notice.Wechat()
		// wechat.Send(&client.Payload{Fo: fo, Data: conf.Tpl})
		timer := time.NewTimer(0 * time.Second)
		defer timer.Stop()
		for i := 0; i < 10; i++ {
			fmt.Println("corn time before")
			select {
			case <-timer.C:
				fmt.Printf("corn time reve %v", conf.Tpl)
				wechat.Send(&client.Payload{Fo: fo, Data: conf.Tpl})
				timer.Reset(30 * time.Second)
			case content := <-ch:
				fmt.Printf("corn reve %v", content)
				reve.Remove(fo)
				return
			}
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
