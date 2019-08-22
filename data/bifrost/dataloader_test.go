package bifrost

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/data/bifrost/streamer"
	"github.com/Mintegral-official/mtggokit/data/container"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type FakeStreamer struct {
	Name      string
	SchedInfo *streamer.SchedInfo
}

func (*FakeStreamer) SetContainer(container.Container) {

}

func (*FakeStreamer) GetContainer() container.Container {
	return nil
}

func (fs *FakeStreamer) UpdateData(ctx context.Context) error {
	fmt.Printf("Name: %s, interval[%d], timestamp[%d]\n", fs.Name, fs.SchedInfo.TimeInterval, time.Now().Unix())
	return nil
}

func (fs *FakeStreamer) GetSchedInfo() *streamer.SchedInfo {
	return fs.SchedInfo
}

func TestLoader_Register(t *testing.T) {
	convey.Convey("Test register duplicate name", t, func() {
		loader := NewLoader()
		convey.So(loader, convey.ShouldNotBeNil)
	})

	convey.Convey("Test register name", t, func() {
		loader := NewLoader()
		convey.So(loader, convey.ShouldNotBeNil)
		convey.So(loader.Register("abc", &FakeStreamer{Name: "fake1"}), convey.ShouldBeNil)
	})

	convey.Convey("Test register duplicate name", t, func() {
		loader := NewLoader()
		convey.So(loader, convey.ShouldNotBeNil)
		convey.So(loader.Register("abc", &FakeStreamer{Name: "fake2"}), convey.ShouldBeNil)
		e := loader.Register("abc", &FakeStreamer{})
		convey.So(e, convey.ShouldNotBeNil)
		convey.So(e.Error(), convey.ShouldEqual, "streamer[abc] has already exist")
	})
}
