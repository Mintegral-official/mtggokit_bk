package bifrost

import (
	"context"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

type FakeStreamer struct {
	Name string
}

func (*FakeStreamer) SetContainer(container.Container) {

}

func (*FakeStreamer) GetContainer() container.Container {
	return nil
}

func (*FakeStreamer) GetSchedInfo() *streamer.SchedInfo {
	return nil
}

func (fs *FakeStreamer) UpdateData(ctx context.Context) error {
	return nil
}

func TestLoader_Register(t *testing.T) {
	convey.Convey("Test register duplicate name", t, func() {
		bifrost := NewBifrost()
		convey.So(bifrost, convey.ShouldNotBeNil)
	})

	convey.Convey("Test register name", t, func() {
		bifrost := NewBifrost()
		convey.So(bifrost, convey.ShouldNotBeNil)
		convey.So(bifrost.Register("abc", &FakeStreamer{Name: "fake1"}), convey.ShouldBeNil)
	})

	convey.Convey("Test register duplicate name", t, func() {
		bifrost := NewBifrost()
		convey.So(bifrost, convey.ShouldNotBeNil)
		convey.So(bifrost.Register("abc", &FakeStreamer{Name: "fake2"}), convey.ShouldBeNil)
		e := bifrost.Register("abc", &FakeStreamer{})
		convey.So(e, convey.ShouldNotBeNil)
		convey.So(e.Error(), convey.ShouldEqual, "streamer[abc] has already exist")
	})
}
