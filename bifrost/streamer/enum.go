package streamer

type UpdatMode int64

const (
	Static = iota
	Dynamic
	Increment
	DynInc
)

var updatModeStrMap = map[UpdatMode]string{
	Static:    "Static",
	Dynamic:   "Dynamic",
	Increment: "Increment",
	DynInc:    "DynInc",
}

func (um *UpdatMode) toString() string {
	v, in := updatModeStrMap[*um]
	if !in {
		return ""
	} else {
		return v
	}
}
