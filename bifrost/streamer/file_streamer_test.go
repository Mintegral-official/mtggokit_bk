package streamer

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalFileStreamer_UpdateData(t *testing.T) {
	curPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fmt.Println(curPath)
	filename := filepath.Join(curPath, "aaa")
	convey.Convey("TestLocalFileStreamer_UpdateData", t, func() {
		lfs := NewFileStreamer(&FileStreamerCfg{
			Name:       "test1",
			Path:       filename,
			UpdatMode:  Dynamic,
			Interval:   2,
			IsSync:     true,
			DataParser: &DefaultTextParser{},
		})
		lfs.SetContainer(&container.BufferedMapContainer{
			Tolerate: 0.1,
		})

		convey.Convey("TestNoFile", func() {

			err := lfs.UpdateData(context.Background())
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(strings.Contains(err.Error(), "no such file or directory"), convey.ShouldBeTrue)
		})

		convey.Convey("TestOneFile", func() {
			s1 := "a\taa\nb\tbb"
			f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
			convey.So(err, convey.ShouldBeNil)
			n, err := f.WriteString(s1)
			convey.So(err, convey.ShouldBeNil)
			convey.So(n, convey.ShouldEqual, len(s1))
			convey.So(f.Close(), convey.ShouldBeNil)
			defer os.Remove(filename)

			err = lfs.UpdateData(context.Background())
			convey.So(err, convey.ShouldBeNil)
			v, err := lfs.GetContainer().Get(container.StrKey("a"))
			convey.So(err, convey.ShouldBeNil)
			convey.So(v, convey.ShouldEqual, "aa")
			v, err = lfs.GetContainer().Get(container.StrKey("b"))
			convey.So(err, convey.ShouldBeNil)
			convey.So(v, convey.ShouldEqual, "bb")

			convey.So(lfs, convey.ShouldNotBeNil)
		})
	})
}
