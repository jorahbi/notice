package aqueue

import (
	"context"
	"sync"

	"github.com/jorahbi/notice/internal/aqueue/jobtype"
	"github.com/jorahbi/notice/internal/notice"
	"github.com/jorahbi/notice/internal/svc"

	"github.com/hibiken/asynq"
)

type Queue struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	cancel func()
	wg     *sync.WaitGroup
}

func NewQueue(svcCtx *svc.ServiceContext, wg *sync.WaitGroup, cancal func()) *Queue {
	return &Queue{
		svcCtx: svcCtx,
		cancel: cancal,
		wg:     wg,
	}
}

// register job 这里一看就和go-zero的router类似
func (l *Queue) Register(ctx context.Context) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	l.wg.Add(1)
	wechatCancel := func() {
		l.wg.Done()
		l.cancel()
	}
	mux.Handle(jobtype.JOB_KEY_WECHAT_NOTICE, notice.NewWechatNoticeHandler(ctx, l.svcCtx, wechatCancel))
	return mux
}
