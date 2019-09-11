package streamer

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func getPath() string {
	fmt.Println(runtime.GOOS)
	if runtime.GOOS == "linux" {
		return filepath.Join("/tmp/bifrost", strconv.FormatInt(time.Now().Unix(), 10), "aaa")
	} else if runtime.GOOS == "darwin" {
		curPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Println(curPath)
		return filepath.Join(curPath, "aaa")
	} else {
		curPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		fmt.Println(curPath)
		return filepath.Join(curPath, "aaa")
	}
}

func TestLocalFileStreamer_UpdateData(t *testing.T) {
	filename := getPath()
	fmt.Println("FilePath:" + filename)
	convey.Convey("TestLocalFileStreamer_UpdateData", t, func() {
		lfs := NewFileStreamer(&LocalFileStreamerCfg{
			Name:       "test1",
			Path:       filename,
			UpdatMode:  Dynamic,
			Interval:   1,
			IsSync:     true,
			DataParser: &DefaultTextParser{},
			Logger:     logrus.New(),
		})
		convey.So(lfs, convey.ShouldNotBeNil)
		lfs.SetContainer(&container.BufferedMapContainer{
			Tolerate: 0.5,
		})

		convey.Convey("TestNoFile", func() {
			ctx, cancel := context.WithTimeout(context.TODO(), time.Microsecond*10)
			defer cancel()
			err := lfs.UpdateData(ctx)
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(strings.Contains(err.Error(), "no such file or directory"), convey.ShouldBeTrue)
		})

		//convey.Convey("TestOneFile", func() {
		//	s1 := "a\taa\nb\tbb"
		//	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		//	convey.So(err, convey.ShouldBeNil)
		//	n, err := f.WriteString(s1)
		//	convey.So(err, convey.ShouldBeNil)
		//	convey.So(n, convey.ShouldEqual, len(s1))
		//	convey.So(f.Close(), convey.ShouldBeNil)
		//	defer os.Remove(filename)
		//
		//	ctx, cancel := context.WithTimeout(context.TODO(), time.Microsecond*10)
		//	defer cancel()
		//	err = lfs.UpdateData(ctx)
		//	convey.So(err, convey.ShouldBeNil)
		//	v, err := lfs.GetContainer().Get(container.StrKey("a"))
		//	convey.So(err, convey.ShouldBeNil)
		//	convey.So(v, convey.ShouldEqual, "aa")
		//	v, err = lfs.GetContainer().Get(container.StrKey("b"))
		//	convey.So(err, convey.ShouldBeNil)
		//	convey.So(v, convey.ShouldEqual, "bb")
		//
		//	convey.Convey("UpatedataFile", func() {
		//		s1 := "ab\tabab\nd\tdd\nee\tEE"
		//		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		//		convey.So(err, convey.ShouldBeNil)
		//		n, err := f.WriteString(s1)
		//		convey.So(err, convey.ShouldBeNil)
		//		convey.So(n, convey.ShouldEqual, len(s1))
		//		convey.So(f.Close(), convey.ShouldBeNil)
		//		defer os.Remove(filename)
		//
		//		ctx, cancel := context.WithTimeout(context.TODO(), time.Microsecond*10)
		//		defer cancel()
		//		err = lfs.UpdateData(ctx)
		//		convey.So(err, convey.ShouldBeNil)
		//
		//		v, err := lfs.GetContainer().Get(container.StrKey("ab"))
		//		convey.So(err, convey.ShouldBeNil)
		//		convey.So(v, convey.ShouldEqual, "abab")
		//
		//		v, err = lfs.GetContainer().Get(container.StrKey("d"))
		//		convey.So(err, convey.ShouldBeNil)
		//		convey.So(v, convey.ShouldEqual, "dd")
		//
		//		v, err = lfs.GetContainer().Get(container.StrKey("ee"))
		//		convey.So(err, convey.ShouldBeNil)
		//		convey.So(v, convey.ShouldEqual, "EE")
		//	})
		//
		//	convey.Convey("TestErrorDataFile", func() {
		//		s1 := "ab\tabab\na\nee"
		//		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		//		convey.So(err, convey.ShouldBeNil)
		//		n, err := f.WriteString(s1)
		//		convey.So(err, convey.ShouldBeNil)
		//		convey.So(n, convey.ShouldEqual, len(s1))
		//		convey.So(f.Close(), convey.ShouldBeNil)
		//		defer os.Remove(filename)
		//
		//		err = lfs.UpdateData(context.Background())
		//		convey.So(err, convey.ShouldNotBeNil)
		//		convey.So(strings.Contains(err.Error(), "LoadBase error, tolerate"), convey.ShouldBeTrue)
		//
		//		v, err := lfs.GetContainer().Get(container.StrKey("a"))
		//		convey.So(err, convey.ShouldBeNil)
		//		convey.So(v, convey.ShouldEqual, "aa")
		//
		//		v, err = lfs.GetContainer().Get(container.StrKey("b"))
		//		convey.So(err, convey.ShouldBeNil)
		//		convey.So(v, convey.ShouldEqual, "bb")
		//	})
		//})
	})
}
