package streamer

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewBiFrostStreamer(t *testing.T) {
	convey.Convey("Test NewBiFrostStreamer", t, func() {
		bs := NewBiFrostStreamer(&BiFrostStreamerCfg{
			Name:         "BifrostStreamer",
			Version:      0,
			BaseFilePath: "",
			Interval:     60,
			IsSync:       true,
			IsOnline:     false,
			WriteFile:    false,
		})
		convey.So(bs, should.NotBeNil)
	})

}
