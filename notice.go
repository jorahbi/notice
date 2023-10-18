package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/hibiken/asynq"
	"github.com/jorahbi/notice/internal/aqueue"
	"github.com/jorahbi/notice/internal/conf"
	"github.com/jorahbi/notice/internal/svc"

	zconf "github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/notice.yaml", "the config file")

func main() {
	// qust := "@gpt  历史学家"
	// keywords := fmt.Sprintf("@gpt%v", string(rune(8197)))
	// // strings.ReplaceAll(qust, string(rune(8197)), " ")
	// idx := strings.Index(qust, keywords)
	// fmt.Println("keywords idx", idx)
	// // fmt.Println(qust, len(qust))
	// fmt.Println(rune(qust[4]), rune(' '), string(rune(8197)), "===")
	// //\x226\x128\x168

	// return
	flag.Parse()
	var c conf.Config
	zconf.MustLoad(*configFile, &c)
	fmt.Println(c.GptKeywords)
	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()
	// 这里可以看源码，类似go-zero的rest，也可以看做http
	job := aqueue.NewQueue(ctx, svcCtx)
	// 注册路由
	mux := job.Register()
	// 启动asynq服务连接redis
	server := newAsynqServer(c.RdsConf)
	if err := server.Run(mux); err != nil {
		logx.WithContext(ctx).Errorf("!!!CronJobErr!!! run err:%+v", err)
		os.Exit(1)
	}
}

func newAsynqServer(c conf.RdsConf) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     c.Addr,
			Password: c.Password,
			PoolSize: c.PoolSize,
		},
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("asynq server exec task IsFailure ======== >>>>>>>>>>> err : %+v  \n", err)
				return true
			},
			Concurrency: 20, //max concurrent process job task nu
		},
	)
}
