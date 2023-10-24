package crontab

import (
	"context"

	"github.com/jorahbi/notice/internal/svc"
)

type nomarl struct{}

func (c nomarl) Start(ctx context.Context, svcCtx *svc.ServiceContext) func() {
	return func() {}
}
func (c nomarl) Stop(ctx context.Context) {

}
