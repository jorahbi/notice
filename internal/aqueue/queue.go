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
}

func NewQueue(ctx context.Context, svcCtx *svc.ServiceContext) *Queue {
	return &Queue{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// register job 这里一看就和go-zero的router类似
func (l *Queue) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	// mux.Handle(jobtype.JOB_KEY_WECHAT_NOTICE, handler.NewOrderNoticeHandler(l.svcCtx))
	mux.Handle(jobtype.JOB_KEY_WECHAT_NOTICE, notice.NewWechatNoticeHandler(l.svcCtx))
	return mux
}
