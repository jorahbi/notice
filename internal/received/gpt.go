package received

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/pkg/client"
	openai "github.com/sashabaranov/go-openai"
)

var lock sync.Mutex
var lockFlag bool

type Gpt struct {
	Model       string       `json:"model"`
	Temperature float32      `json:"temperature"`
	Messages    []*GptReqMsg `json:"messages"`
}

type GptReqMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewGpt() *Gpt {
	return &Gpt{}
}

func (gpt *Gpt) Event(ctx context.Context, svcConf conf.Config, payload client.Payload) (string, error) {
	if lockFlag {
		return "", errors.New("忙着呢，稍候在问!!")
	}
	lock.Lock()
	lockFlag = true

	defer func() {
		lockFlag = false
		lock.Unlock()
	}()
	if len(payload.String()) <= len(svcConf.GptKeywords) {
		return "", errors.New("请说出你想问的问题")
	}
	qust := payload.String()
	idx := strings.Index(qust, svcConf.GptKeywords)
	qust = strings.Trim(qust[idx:], "")
	if len(qust) == 0 {
		return "", errors.New("请输入正确的问题")
	}
	fmt.Printf("开始提问%v", qust)
	client := openai.NewClient(svcConf.GptKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{{
				Role:    openai.ChatMessageRoleUser,
				Content: qust},
			},
		},
	)
	chLen := len(resp.Choices)
	if err != nil || chLen == 0 {
		return "", fmt.Errorf("ChatCompletion len = 0 or error: %v", err)
	}

	return resp.Choices[chLen-1].Message.Content, nil
}
