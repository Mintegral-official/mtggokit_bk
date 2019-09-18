package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
	"os"
	"os/signal"
	"time"
)

func main() {

	sp1 := streamer.NewStreamerProvider(&streamer.StreamerProviderCfg{
		Name:       "sp1",
		ExpireTime: 10,
		Logger:     logrus.New(),
	})
	spm := streamer.NewStreamerProviderManager()
	if e := spm.RegiterProvider("sp1", sp1); e != nil {
		fmt.Println("Register error:" + e.Error())
	}

	biServer := streamer.NewBifrostServer(spm)
	s := server.NewServer()
	if e := s.RegisterName("BifrostServer", biServer, ""); e != nil {
		fmt.Println("RegisterName BifrostServer error")
		os.Exit(-1)
	}

	go func() {
		if e := s.Serve("tcp", "localhost:7878"); e != nil {
			fmt.Println("server closed, err=" + e.Error())
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			t := time.After(time.Second)
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

	go func() {
		for {
			t := time.After(time.Millisecond * 300)
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
	_ = s.Close()
}
