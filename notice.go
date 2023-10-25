package main

import (
	"flag"

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
	notice.Start(svcCtx, crontab.NewCrontab(), aqueue.NewAsynq())
}
