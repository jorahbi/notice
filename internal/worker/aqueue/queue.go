package aqueue

import (
	"context"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/logx"

	// "github.com/jorahbi/notice/internal/aqueue"
	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/notice"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/internal/worker/aqueue/jobtype"
)

type Asynq struct{}

func NewAsynq() *Asynq {
	return &Asynq{}
}

func (q *Asynq) Start(ctx context.Context, svc *svc.ServiceContext) {
	server := q.newAsynq(svc.Config.RdsConf)
	mux := q.register(ctx, svc)

	if err := server.Start(mux); err != nil {
		logx.WithContext(ctx).Errorf("!!!CronJobErr!!! run err:%+v", err)
		os.Exit(1)
	}
	for {
		<-ctx.Done()
		server.Shutdown()
		fmt.Println("queue down")
		return
	}
}

func (q *Asynq) newAsynq(c conf.RdsConf) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     c.Addr,
			Password: c.Password,
			PoolSize: c.PoolSize,
		},
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("asynq server exec task IsFailure ======== >>>>>>>>>>> err : %+v  \n", err)
				return true
			},
			Concurrency: 20, //max concurrent process job task nu
		},
	)
}

// register job 这里一看就和go-zero的router类似
func (q *Asynq) register(ctx context.Context, svc *svc.ServiceContext) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.Handle(jobtype.JOB_KEY_WECHAT_NOTICE, notice.Wechat())
	return mux
}
