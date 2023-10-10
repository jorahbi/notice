package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"notice/internal/aqueue/jobtype"
	"notice/internal/aqueue/queue"
	"notice/internal/notice"
	"os"
	"text/template"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/notice.yaml", "the config file")

func main() {
	flag.Parse()
	// 这里还是去引用go-zero的配置
	// flag.Parse()

	var c notice.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := notice.NewServiceContext(c)
	ctx := context.Background()
	// 这里可以看源码，类似go-zero的rest，也可以看做http
	job := queue.NewQueue(ctx, svcCtx)
	// 注册路由
	mux := job.Register()
	// 启动asynq服务连接redis
	server := notice.NewAsynqServer(c.RdsConf)
	if err := server.Run(mux); err != nil {
		logx.WithContext(ctx).Errorf("!!!CronJobErr!!! run err:%+v", err)
		os.Exit(1)
	}
}

func test() {
	var p jobtype.PayloadOrderNotice
	payload := `{
		"Order":{
			"oid":447129156399749,
			"uid":420368187617925,
			"member":1,
			"goods_info":"",
			"state":0,
			"money":0,
			"pre_money":0,
			"pre_datetime":"2023-08-31T13:13:03+08:00",
			"phone":"13077367670",
			"extra":"",
			"desktop":1,
			"operation":0,
			"create_time":"2023-08-06T13:13:13.69958+08:00",
			"update_time":"0001-01-01T00:00:00Z"
		},
		"Goods":[
			{
				"group":447128567702149,
				"gid":6,
				"desktop":1,
				"num":2,
				"name":"test",
				"subTitle":"test",
				"img":"/static/images/WechatIMG238.jpg",
				"price":111
			}
		]
	}`
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		fmt.Println(err.Error())
		return
	}

	notice := ""
	if p.Order.State == 1 {
		notice = `
		订单号: {{ .Order.Oid }}
		类型: 堂食
		{{- range $index, $value := .Goods -}}
		菜名: {{ $value.Name }}
		价格*价格: {{ $value.Price }} * {{ $value.Num }}
		{{ end }}
		`
	} else {
		notice = `
	   订单号: {{ .Order.Oid }}
		类型: 预定
	    用餐时间: {{ .Order.PreDatetime }}
		人数: {{ .Order.Member }}
		桌号: {{ .Order.Desktop }}
		{{- range $index, $value := .Goods }}
		菜名: {{ $value.Name }}
		价格*价格: {{ $value.Price }} * {{ $value.Num }}
		{{ end }}
		`
	}
	builder := bytes.NewBuffer([]byte{})
	tpl := template.Must(template.New("query").Parse(notice))
	tpl.Execute(builder, p)
	fmt.Println(string(builder.Bytes()))
}
