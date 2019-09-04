package streamer

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewMongoStreamer(t *testing.T) {
	convey.Convey("", t, func() {
		ms := NewMongoStreamer(&MongoStreamerCfg{
			Name:        "mongo_test",
			UpdatMode:   Dynamic,
			IncInterval: 60,
			IsSync:      true,
			IP:          "127.0.0.1",
			Port:        21017,
			BaseParser:  &DefaultTextParser{},
			IncParser:   &DefaultTextParser{},
			BaseQuery:   "mongo base query",
			IncQuery:    "mongo inc query",
			UserData:    "user defined data",
			OnIncFinish: func(userData interface{}) interface{} {
				return "nfew inc base query"
			},
		})
		convey.So(ms, convey.ShouldNotBeNil)
	})
}
