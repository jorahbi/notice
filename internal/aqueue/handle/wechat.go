package handler

import (
	"context"
	"encoding/json"
	"notice/internal/aqueue/jobtype"
	"notice/internal/notice"

	"github.com/hibiken/asynq"
)

type WechatNoticeHandler struct {
	svcCtx *notice.ServiceContext
}

func NewWechatNoticeHandler(svcCtx *notice.ServiceContext) *OrderNoticeHandler {
	return &OrderNoticeHandler{svcCtx: svcCtx}
}

func (l *WechatNoticeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p jobtype.PayloadNotice
	var err error
	if err = json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	return nil
}
