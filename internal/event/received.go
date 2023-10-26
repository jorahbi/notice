package event

import (
	"context"

	"github.com/eatmoreapple/openwechat"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/orcaman/concurrent-map/v2"
)

// 回复
type received struct {
	reveHub cmap.ConcurrentMap[string, chan string]
}

const (
	REVE_CONTENT = "default"
)

var reve = &received{
	reveHub: cmap.New[chan string](),
}

func ReveInstance() *received {
	return reve
}

func (e *received) Event(ctx context.Context, svcCtx *svc.ServiceContext, msg *openwechat.Message) (string, error) {
	user, err := msg.Sender()
	if err != nil {
		return "", nil
	}
	reve := e.Get(user.RemarkName)
	reve <- msg.Content

	return "", nil
}

func (e *received) Get(remarkName string) chan string {
	reve, ok := e.reveHub.Get(remarkName)
	if !ok {
		reve = make(chan string)
		e.reveHub.Set(remarkName, reve)
	}
	return reve
}

func (e *received) Remove(remarkName string) {
	e.reveHub.Remove(remarkName)
}
