package crontab

import (
	"context"
	"fmt"
	"sync"

	"github.com/jorahbi/notice/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

const JOB_NOMARL = "nomarl"

type CrontabIface interface {
	Start(ctx context.Context, svcCtx *svc.ServiceContext) func()
	Stop(ctx context.Context)
}

var (
	ctx    context.Context
	cancel context.CancelFunc
	jobs   = map[string]CrontabIface{
		JOB_NOMARL: &nomarl{},
	}
)

func CronRun(ctx context.Context, svc *svc.ServiceContext, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	// ctx, cancel = context.WithCancel(ctx)
	c := cron.New()
	for _, jobConf := range svc.Config.Notices {
		job, ok := jobs[jobConf.Mode]
		if !ok {
			panic("job not found")
		}
		_, err := c.AddFunc(jobConf.Spec, job.Start(ctx, svc))
		if err != nil {
			logx.Errorf("job exec error name[%v] error[%v]", jobConf.Spec, err)
		}
	}
	waitStop(ctx, c)
}

func waitStop(ctx context.Context, c *cron.Cron) {
	select {
	case <-ctx.Done():
		select {
		case <-c.Stop().Done():
			fmt.Println("job down")
			return
		}

	}

}
