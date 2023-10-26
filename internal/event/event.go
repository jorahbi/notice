package event

import (
	"context"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/jorahbi/notice/internal/svc"
	"github.com/sashabaranov/go-openai"
)

type EventType string

const (
	EVENT_KEY_GPT  EventType = "gpt"
	EVENT_KEY_REVE EventType = "reve"
)

type ReveResp interface {
	openai.ChatCompletionResponse
}

type EventInterface interface {
	Event(ctx context.Context, svcCtx *svc.ServiceContext, msg *openwechat.Message) (string, error)
}

var Events = map[EventType]EventInterface{
	EVENT_KEY_GPT:  &gpt{},
	EVENT_KEY_REVE: reve,
}

func keywords(msg string, keywords []string) int {
	for _, keys := range keywords {
		idx := strings.Index(msg, keys)
		if idx >= 0 {
			return idx + len(keys)
		}
	}
	return -1
}
