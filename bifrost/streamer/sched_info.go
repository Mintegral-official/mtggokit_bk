package streamer

type UpdatMode int64

const (
	Static   UpdatMode = 0
	Dynamic  UpdatMode = 1
	Increase UpdatMode = 2
	DynInc   UpdatMode = 3
)

type SchedInfo struct {
	UpdateMode   UpdatMode
	TimeInterval int
	IsDetached   bool
	IsAsync      bool
}
