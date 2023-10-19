package aqueue

import (
	"context"

	"github.com/jorahbi/notice/internal/aqueue/jobtype"
	"github.com/jorahbi/notice/internal/notice"
	"github.com/jorahbi/notice/internal/svc"

	"github.com/hibiken/asynq"
)

type Queue struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	cancel func()
}

func NewQueue(svcCtx *svc.ServiceContext, cancal func()) *Queue {
	return &Queue{
		svcCtx: svcCtx,
		cancel: cancal,
	}
}

// register job 这里一看就和go-zero的router类似
func (l *Queue) Register(ctx context.Context) *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.Handle(jobtype.JOB_KEY_WECHAT_NOTICE, notice.NewWechatNoticeHandler(ctx, l.svcCtx, l.cancel))
	return mux
}
