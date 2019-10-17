package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

func main() {

	sp1 := streamer.NewStreamerProvider(&streamer.StreamerProviderCfg{
		Name:       "bifrost_streamer_example",
		ExpireTime: 10,
		Logger:     logrus.New(),
	})

	bs := bifrost.NewBifrostServer(&bifrost.BifrostServerCfg{
		Addr:   "localhost:7878",
		Logger: logrus.New(),
	})

	if err := bs.RegisterProvider("bifrost_streamer_example", sp1); err != nil {
		fmt.Println("Rigister provider error," + err.Error())
		os.Exit(-1)
	}

	// 设置基准
	sp1.SetBase(&streamer.BaseInfo{
		Name:     "bifrost_streamer_example",
		Progress: 5,
		Data: map[container.MapKey]interface{}{
			container.StrKey("1"): 1,
			container.StrKey("2"): 4,
		},
	})

	//启动服务
	go func() {
		if e := bs.Serve(); e != nil {
			fmt.Println("server closed, err=" + e.Error())
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	// 写增量
	go func() {
		for {
			t := time.After(10 * time.Second)
			select {
			case <-t:
				keyPrefix := time.Now().Format("2019-01-01_11:12:13")
				progress := time.Now().Unix()
				sp1.AddInc([]*streamer.IncRecord{
					{
						MapKey:   container.StrKey(keyPrefix + "_1"),
						Progress: progress,
						Value:    keyPrefix + "_value1",
					},
					{
						MapKey:   container.StrKey(keyPrefix + "_2"),
						Progress: progress,
						Value:    keyPrefix + "_value2",
					},
					{
						MapKey:   container.StrKey(keyPrefix + "_3"),
						Progress: progress,
						Value:    keyPrefix + "_value3",
					},
				})
			case <-ctx.Done():
				fmt.Println("Inc Finish")
				return
			}
		}
	}()

	// 读增量
	go func() {
		for {
			t := time.After(time.Millisecond * 3000)
			select {
			case <-t:
				records, err := sp1.GetInc(time.Now().Add(-time.Second).Unix(), 2)
				if err != nil {
					fmt.Println("GetIncError: " + err.Error())
					continue
				}
				for _, r := range records {
					fmt.Println(r)
				}
			}
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c)
	fmt.Println("退出信号", <-c)
	cancel()
	_ = bs.Close()
}
