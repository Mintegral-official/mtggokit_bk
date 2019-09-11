package bifrost

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
)

func main() {

	// init
	Bifrost := bifrost.NewBifrost() // new a bifronst object
	s := streamer.NewFileStreamer(&streamer.LocalFileStreamerCfg{
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
	if err := Bifrost.Register("example1", s); err != nil {
		fmt.Println("Register error", err.Error())
	}

	// use
	value, err := Bifrost.Get("exmaple1", container.StrKey("key"))
	if err != nil {
		fmt.Println(value)
	}
}
