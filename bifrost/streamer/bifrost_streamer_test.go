package streamer

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewBiFrostStreamer(t *testing.T) {
	convey.Convey("Test NewBiFrostStreamer", t, func() {
		bs := NewBiFrostStreamer(&BiFrostStreamerCfg{
			Name:         "BiFrostStreamer",
			Version:      0,
			Ip:           "",
			Port:         1111,
			BaseFilePath: "",
			Interval:     60,
			IsSync:       true,
			IsOnline:     false,
			WriteFile:    false,
			CacheSize:    10000,
		})
		convey.So(bs, should.NotBeNil)
	})

}
