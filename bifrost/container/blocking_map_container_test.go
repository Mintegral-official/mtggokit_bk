package container

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"sync/atomic"
	"testing"
	"unsafe"
)

func TestBlockingMap_Get(t *testing.T) {
	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1)
		convey.So(bm.LoadBase(NewTestDataIter([]string{})), convey.ShouldBeNil)
		convey.So(bm.ErrorNum, convey.ShouldEqual, 0)
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1)
		convey.So(bm.LoadBase(NewTestDataIter([]string{
			"1\t2",
			"a\tb",
		})), convey.ShouldBeNil)
		convey.So(bm.ErrorNum, convey.ShouldEqual, 0)
		convey.So(bm.innerData, convey.ShouldNotBeNil)
		v, e := bm.Get(StrKey("1"))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(StrKey("a"))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")
	})

	convey.Convey("Test BufferedMapContainer Get", t, func() {
		bm := CreateBlockingMapContainer(1)
		convey.So(bm.LoadBase(NewTestIntDataIter([]string{
			"1\t2",
			"4\tb",
		})), convey.ShouldBeNil)
		convey.So(bm.ErrorNum, convey.ShouldEqual, 0)
		v, e := bm.Get(I64Key(1))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "2")

		v, e = bm.Get(I64Key(4))
		convey.So(e, convey.ShouldBeNil)
		convey.So(v, convey.ShouldEqual, "b")
	})
}

func TestAomicSwap(t *testing.T) {
	p := &map[int]int{1: 2}
	q := &map[int]int{3: 4}
	fmt.Printf("%v, %p, %v, %p\n", p, &p, q, &q)

	u_p := unsafe.Pointer(p)
	u_q := unsafe.Pointer(q)
	fmt.Printf("%v, %v\n", u_p, u_q)
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	atomic.CompareAndSwapPointer(&u_p, u_p, u_q)
	fmt.Printf("%v, %v\n", u_p, u_q)
	fmt.Printf("%v, %p, %v, %p\n", p, &p, q, &q)
}
