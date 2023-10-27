package event

import (
	"context"
	"fmt"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/orcaman/concurrent-map/v2"
)

// 回复
type received struct {
	reveHub cmap.ConcurrentMap[string, chan string]
}

type WaitConf struct {
	Fo       string
	Time     time.Duration
	Num      int
	Callback func()
}

const (
	REVE_CONTENT = "default"
)

var reve = &received{
	reveHub: cmap.New[chan string](),
}

func WaitReve() *received {
	return reve
}

func (e *received) Event(ctx context.Context, svcCtx *svc.ServiceContext, msg *openwechat.Message) (string, error) {
	user, err := msg.Sender()
	if err != nil {
		return "", nil
	}
	reve, ok := e.reveHub.Get(user.RemarkName)
	if !ok {
		return "", nil
	}
	reve <- msg.Content

	return "", nil
}

func (e *received) MustGet(remarkName string) chan string {
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

func (e *received) Wait(ctx context.Context, wait WaitConf) {
	ch := e.MustGet(wait.Fo)
	timer := time.NewTimer(0 * time.Second)
	defer timer.Stop()
	for i := 0; i < wait.Num; i++ {
		select {
		case <-timer.C:
			wait.Callback()
			timer.Reset(wait.Time)
		case content := <-ch:
			fmt.Printf("corn reve down %v", content)
			reve.Remove(wait.Fo)
			return
		case <-ctx.Done():
			fmt.Printf("context exit ")
			return
		}
	}
}
