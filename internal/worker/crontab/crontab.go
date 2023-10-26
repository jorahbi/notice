package crontab

import (
	"context"
	"fmt"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

const JOB_NOMARL = "nomarl"

type CrontabIface interface {
	Start(ctx context.Context, svcCtx *svc.ServiceContext, conf conf.Notice) func()
	Stop(ctx context.Context)
}

var (
	ctx    context.Context
	cancel context.CancelFunc
)

type Crontab struct {
	jobs map[string]CrontabIface
}

func NewCrontab() *Crontab {
	return &Crontab{
		jobs: map[string]CrontabIface{
			JOB_NOMARL: &nomarl0800{},
		},
	}
}

func (crontab *Crontab) Start(ctx context.Context, svc *svc.ServiceContext) {
	// ctx, cancel = context.WithCancel(ctx)
	c := cron.New()
	fmt.Println("crontab start ...")
	for _, jobConf := range svc.Config.Notices {
		job, ok := crontab.jobs[jobConf.Mode]
		if !ok {
			panic("job not found")
		}
		fmt.Println(jobConf.Spec)
		_, err := c.AddFunc(jobConf.Spec, job.Start(ctx, svc, jobConf))
		if err != nil {
			logx.Errorf("job exec error name[%v] error[%v]", jobConf.Spec, err)
		}
	}
	c.Start()
	wait(ctx, c)
}

func wait(ctx context.Context, c *cron.Cron) {
	for {
		select {
		case <-ctx.Done():
			select {
			case <-c.Stop().Done():
				fmt.Println("job down")
				return
			}

		}
	}
}
