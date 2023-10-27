package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/notice"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/internal/worker/aqueue"
	"github.com/jorahbi/notice/internal/worker/crontab"
	zconf "github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/notice.yaml", "the config file")

func main() {
	flag.Parse()
	var c conf.Config
	zconf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	chat := notice.NewWechat(svcCtx)
	chat.Start(crontab.NewCrontab(), aqueue.NewAsynq())
	// chat.Start()
}

func test() {
	timer := time.NewTimer(0 * time.Second)
	defer timer.Stop()
	for i := 0; i < 10; i++ {
		select {
		case <-timer.C:
			fmt.Println("corn time reve ")
			timer.Reset(30 * time.Second)
		}
	}
}
