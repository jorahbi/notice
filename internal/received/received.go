package received

import (
	"context"

	"github.com/jorahbi/notice/pkg/client"
	"github.com/sashabaranov/go-openai"
)

const (
	RECE_KEY_GPT = "@gpt"
)

type ReveResp interface {
	openai.ChatCompletionResponse
}

type EventInterface interface {
	Event(ctx context.Context, payload client.Payload) (string, error)
}
