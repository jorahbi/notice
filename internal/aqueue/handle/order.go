package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jorahbi/notice/internal/aqueue/jobtype"
	"github.com/jorahbi/notice/internal/notice"
	"text/template"

	"github.com/go-vgo/robotgo"
	"github.com/hibiken/asynq"
)

//https://learnku.com/articles/75268

type OrderNoticeHandler struct {
	svcCtx *notice.ServiceContext
}

func NewOrderNoticeHandler(svcCtx *notice.ServiceContext) *OrderNoticeHandler {
	return &OrderNoticeHandler{svcCtx: svcCtx}
}

func (l *OrderNoticeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p jobtype.PayloadOrderNotice
	var err error
	if err = json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	fmt.Println(p, string(t.Payload()))
	fpid, err := robotgo.FindIds("wechat")
	if len(fpid) == 0 || err != nil {
		fmt.Println("没有获取到wechat窗口")
	}
	// fmt.Println(robotgo.ActivePID(5991))
	isExist, err := robotgo.PidExists(fpid[0])
	if err != nil && !isExist {
		fmt.Println("微信窗口已关闭，请重新打开", err.Error())
		return nil
	}
	robotgo.ActivePID(fpid[0])
	notice := ""
	if p.Order.State == 1 {
		notice = `
		订单号: {{ .Order.Oid }}
		下单时间: {{ .Order.CreateTime }}
		类型: 堂食
		{{- range $index, $value := .Goods -}}
		菜名: {{ $value.Name }}
		价格*数量: {{ $value.Price }} * {{ $value.Num }}
		{{ end }}
		`
	} else {
		notice = `
	   订单号: {{ .Order.Oid }}
	   下单时间: {{ .Order.CreateTime }}
		类型: 预定
	    用餐时间: {{ .Order.PreDatetime }}
		电话: {{ .Order.Phone }}
		人数: {{ .Order.Member }}
		桌号: {{ .Order.DesktopName }}
		{{- range $index, $value := .Goods }}
		菜名: {{ $value.Name }}
		价格*数量: {{ $value.Price }} * {{ $value.Num }}
		{{ end }}
		`
	}
	builder := bytes.NewBuffer([]byte{})
	tpl := template.Must(template.New("query").Parse(notice))
	p.Order.DesktopName = desktop(int(p.Order.Desktop))
	tpl.Execute(builder, p)

	robotgo.TypeStr(builder.String())
	robotgo.KeyTap("enter")
	fmt.Println(robotgo.GetActive(), robotgo.GetTitle())
	return nil
}

func desktop(idx int) string {
	mapping := map[int]string{
		0: "大包房", 1: "小包房", 2: "大圆桌", 3: "小圆桌", 4: "方桌", 5: "条桌1", 6: "条桌2",
	}
	if result, isOk := mapping[idx]; isOk {
		return result
	}
	return ""
}
