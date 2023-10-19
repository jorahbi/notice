package event

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/pkg/client"
	"github.com/zeromicro/go-zero/core/jsonx"
)

type Http[T ReveResp] struct {
	client *http.Client
}

func NewHttp[T ReveResp]() *Http[T] {
	c := &http.Client{Timeout: 180 * time.Second}
	return &Http[T]{client: c}
}

func (h *Http[T]) Send(ctx context.Context, conf conf.ReveConfig, payload client.Payload) (*T, error) {
	body := bytes.Buffer{}
	body.WriteString(payload.String())
	request, err := http.NewRequest(conf.Http.Method, conf.Http.Url, &body)
	if err != nil {
		return nil, err
	}
	for key, val := range conf.Http.Header {
		request.Header.Set(key, strings.Join(val, " "))
	}

	resp, err := h.client.Do(request)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	result := new(T)
	defer resp.Body.Close()
	err = jsonx.UnmarshalFromReader(resp.Body, result)
	fmt.Println("response", result)
	return result, err
}
