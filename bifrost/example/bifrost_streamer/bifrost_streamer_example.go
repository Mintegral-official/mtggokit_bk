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
)

func main() {
	// init
	Bifrost := bifrost.NewBifrost() // new a bifronst object
	s := streamer.NewBiFrostStreamer(&streamer.BiFrostStreamerCfg{
		Name:     "bifrost_streamer_example",
		URI:      "localhost:7878",
		Interval: 1,
		IsSync:   true,
		Logger:   logrus.New(),
	})
	c := &container.BufferedMapContainer{}
	s.SetContainer(c)
	ctx, cancel := context.WithCancel(context.Background())
	_ = s.UpdateData(ctx)
	if err := Bifrost.Register("example1", s); err != nil {
		fmt.Println("Register error", err.Error())
	}

	// use
	value, err := Bifrost.Get("exmaple1", container.StrKey("key"))
	if err != nil {
		fmt.Println("xxxxxxxxxxxxxxxx", value)
	}

	ch := make(chan os.Signal)
	signal.Notify(ch)
	fmt.Println("退出信号", <-ch)
	cancel()
}
