package event

import (
	"context"
	"strings"

	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
	"github.com/sashabaranov/go-openai"
)

type EventType string

const (
	EVENT_KEY_GPT EventType = "gpt"
)

type ReveResp interface {
	openai.ChatCompletionResponse
}

type EventInterface interface {
	Event(ctx context.Context, svcCtx *svc.ServiceContext, payload *client.Payload) (string, error)
}

var Events = map[EventType]EventInterface{
	EVENT_KEY_GPT: &gpt{},
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
