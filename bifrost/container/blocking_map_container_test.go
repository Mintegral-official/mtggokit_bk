package container

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBlockingMapContainer_LoadBase(t *testing.T) {
	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1, 0)
		convey.So(bm.LoadBase(NewTestDataIter([]string{})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 0)
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1, 0)
		convey.So(bm.LoadBase(NewTestDataIter([]string{
			"1\t2",
			"a\tb",
		})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 0)
		convey.So(bm.innerData, convey.ShouldNotBeNil)
		v, e := bm.Get(StrKey("1"))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(StrKey("a"))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1, 0)
		convey.So(bm.LoadBase(NewTestIntDataIter([]string{
			"1\t2",
			"4\tb",
		})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 0)
		v, e := bm.Get(I64Key(1))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(I64Key(4))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1, 0)
		convey.So(bm.LoadBase(NewTestIntDataIter([]string{
			"1\t2",
			"4\tb",
			"2",
		})), convey.ShouldNotBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 1)
		convey.So(bm.totalNum, convey.ShouldEqual, 3)
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1, 0.5)
		convey.So(bm.LoadBase(NewTestIntDataIter([]string{
			"1\t2",
			"4\tb",
			"2",
		})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 1)
		convey.So(bm.totalNum, convey.ShouldEqual, 3)

		v, e := bm.Get(I64Key(1))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(I64Key(4))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")
	})

}

func TestBlockingMapContainer_LoadInc(t *testing.T) {
	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1, 0.5)
		convey.So(bm.LoadBase(NewTestIntDataIter([]string{
			"1\t2",
			"4\tb",
			"2",
		})), convey.ShouldBeNil)
		convey.So(bm.errorNum, convey.ShouldEqual, 1)
		convey.So(bm.totalNum, convey.ShouldEqual, 3)

		v, e := bm.Get(I64Key(1))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(I64Key(4))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")

		convey.Convey("Test LoadIncSucc", func() {
			convey.So(bm.LoadInc(NewTestIntDataIter([]string{
				"5\t3",
				"2",
			})), convey.ShouldBeNil)
			convey.So(bm.errorNum, convey.ShouldEqual, 2)
			convey.So(bm.totalNum, convey.ShouldEqual, 5)
		})

		convey.Convey("Test LoadInc Fail", func() {
			convey.So(bm.LoadInc(NewTestIntDataIter([]string{
				"5",
				"2",
			})), convey.ShouldNotBeNil)
			convey.So(bm.errorNum, convey.ShouldEqual, 3)
			convey.So(bm.totalNum, convey.ShouldEqual, 5)
		})
	})
}
