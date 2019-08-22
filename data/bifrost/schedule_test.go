package bifrost

import (
	"context"
	"github.com/Mintegral-official/mtggokit/data/bifrost/streamer"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestSched_Schedule(t *testing.T) {
	convey.Convey("Test schedule", t, func() {
		loader := NewLoader()
		convey.So(loader, convey.ShouldNotBeNil)
		convey.So(loader.Register("test1", &FakeStreamer{
			Name:      "fake1",
			SchedInfo: &streamer.SchedInfo{TimeInterval: 1},
		}), convey.ShouldBeNil)
		convey.So(loader.Register("test3", &FakeStreamer{
			Name:      "fake3",
			SchedInfo: &streamer.SchedInfo{TimeInterval: 3},
		}), convey.ShouldBeNil)
		convey.So(loader.Register("test7", &FakeStreamer{
			Name:      "fake7",
			SchedInfo: &streamer.SchedInfo{TimeInterval: 7},
		}), convey.ShouldBeNil)
		convey.So(loader.Register("test5", &FakeStreamer{
			Name:      "fake5",
			SchedInfo: &streamer.SchedInfo{TimeInterval: 5},
		}), convey.ShouldBeNil)

		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()
		loader.sched.Schedule(ctx)
	})
}
