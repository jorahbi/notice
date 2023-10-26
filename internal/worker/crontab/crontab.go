package crontab

import (
	"context"
	"fmt"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

const JOB_MORNING = "morning"

type CrontabIface interface {
	Start(ctx context.Context, svcCtx *svc.ServiceContext, conf conf.Job) func()
}

var (
	ctx    context.Context
	cancel context.CancelFunc
)

type jobfn func(ctx context.Context, svcCtx *svc.ServiceContext, fo string, conf conf.Job) func()

type Crontab struct {
	jobs map[string]jobfn
}

func NewCrontab() *Crontab {
	return &Crontab{
		jobs: map[string]jobfn{
			JOB_MORNING: morning,
		},
	}
}

func (crontab *Crontab) Start(ctx context.Context, svc *svc.ServiceContext) {
	// ctx, cancel = context.WithCancel(ctx)
	c := cron.New()
	fmt.Println("crontab start ...")
	for _, jobConf := range svc.Config.Jobs {
		job, ok := crontab.jobs[jobConf.Name]
		if !ok {
			logx.Infof("job[%v] not found", jobConf.Name)
			continue
		}
		for _, item := range jobConf.To {
			_, err := c.AddFunc(jobConf.Spec, job(ctx, svc, item, jobConf))
			if err != nil {
				logx.Errorf("job exec error name[%v] error[%v]", jobConf.Spec, err)
			}
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
