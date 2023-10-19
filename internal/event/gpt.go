package event

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/jsonx"

	"github.com/jorahbi/notice/internal/svc"
	"github.com/jorahbi/notice/pkg/client"
	openai "github.com/sashabaranov/go-openai"
)

var lock sync.Mutex
var lockFlag bool

type gpt struct{}

type GptReqMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func newGpt() *gpt {
	return &gpt{}
}

func (e *gpt) Event(ctx context.Context, svcCtx *svc.ServiceContext, payload client.Payload) (string, error) {
	if lockFlag {
		return "", errors.New("忙着呢，稍候在问!!")
	}
	lock.Lock()
	lockFlag = true

	defer func() {
		lockFlag = false
		lock.Unlock()
	}()
	qust := strings.Trim(payload.String(), " ")
	idx := e.isCall(ctx, svcCtx, qust)
	if idx < 0 {
		return "", nil
	}
	qust = qust[idx+len(svcCtx.Config.GptKeywords):]
	if len(qust) == 0 {
		return "", errors.New("请说出你想问的问题")
	}
	//return e.proxy(qust)
	return e.qustion(ctx, svcCtx, qust)
}

func (e *gpt) qustion(ctx context.Context, svcCtx *svc.ServiceContext, qust string) (string, error) {
	fmt.Printf("开始提问:%v", qust)
	cmd := exec.Command("curl", "--max-time", "300", "--request", "POST", `https://api.openai.com/v1/chat/completions`,
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("Authorization: Bearer %v", svcCtx.Config.GptKey),
		"-d", fmt.Sprintf(`{"model":"%v","messages":[{"role":"%v","content":"%v"}]}`,
			openai.GPT3Dot5Turbo, openai.ChatMessageRoleUser, qust))
	out, err := cmd.Output()

	resp := &openai.ChatCompletionResponse{}
	err = jsonx.Unmarshal(out, resp)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	chLen := len(resp.Choices)
	if chLen == 0 {
		fmt.Println(string(out))
		return "", fmt.Errorf("ChatCompletion len = 0 or error: %v", err)
	}

	return resp.Choices[chLen-1].Message.Content, nil
}

func (e *gpt) isCall(ctx context.Context, svcCtx *svc.ServiceContext, msg string) int {
	idx := strings.Index(msg, svcCtx.Config.GptKeywords)
	if idx >= 0 {
		return idx
	}
	idx = strings.Index(msg, "justyolo ")
	if idx >= 0 {
		return idx
	}
	return -1
}

func (e *gpt) proxy(ctx context.Context, svcCtx *svc.ServiceContext, qust string) (string, error) {
	os.Setenv("HTTP_PROXY", svcCtx.Config.Proxy)
	os.Setenv("HTTPS_PROXY", svcCtx.Config.Proxy)
	defer func() {
		os.Unsetenv("HTTP_PROXY")
		os.Unsetenv("HTTPS_PROXY")
		os.Unsetenv("NO_PROXY")
	}()
	config := openai.DefaultConfig(svcCtx.Config.GptKey)
	config.HTTPClient.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}
	fmt.Printf("开始提问%v", qust)
	client := openai.NewClientWithConfig(config)
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
	fmt.Println("返回", err)
	chLen := len(resp.Choices)
	if err != nil || chLen == 0 {
		return "", fmt.Errorf("ChatCompletion len = 0 or error: %v", err)
	}

	return resp.Choices[chLen-1].Message.Content, nil
}