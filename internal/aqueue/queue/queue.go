package queue

import (
	"context"
	handler "notice/internal/aqueue/handle"
	"notice/internal/aqueue/jobtype"
	"notice/internal/notice"

	"github.com/hibiken/asynq"
)

type Queue struct {
	ctx    context.Context
	svcCtx *notice.ServiceContext
}

func NewQueue(ctx context.Context, svcCtx *notice.ServiceContext) *Queue {
	return &Queue{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// register job 这里一看就和go-zero的router类似
func (l *Queue) Register() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	// mux.Handle(jobtype.JOB_KEY_ORDER_NOTICE, handler.NewOrderNoticeHandler(l.svcCtx))
	mux.Handle(jobtype.JOB_KEY_WECHAT_NOTICE, handler.NewWechatNoticeHandler(l.svcCtx))
	return mux
}
