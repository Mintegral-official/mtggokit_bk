package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
)

func main() {

	// init
	bifrost := bifrost.NewBifrost() // new a bifronst object
	s, _ := streamer.NewFileStreamer(&streamer.FileStreamerCfg{
		Name:       "example1",
		Path:       "a.txt",
		Interval:   60,
		IsSync:     true,
		DataParser: &streamer.DefaultTextParser{},
		UserData:   nil,
	})
	c := &container.BufferedMapContainer{}
	s.SetContainer(c)
	_ = s.UpdateData(context.Background())
	if err := bifrost.Register("example1", s); err != nil {
		fmt.Println("Register error", err.Error())
	}

	// use
	value, err := bifrost.Get("exmaple1", container.StrKey("key"))
	if err != nil {
		fmt.Println(value)
	}
}
