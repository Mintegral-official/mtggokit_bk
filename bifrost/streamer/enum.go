package streamer

type UpdatMode int64

const (
	Static = iota
	Dynamic
	Increase
	DynInc
)

var updatModeStrMap = map[UpdatMode]string{
	Static:   "Static",
	Dynamic:  "Dynamic",
	Increase: "Increase",
	DynInc:   "DynInc",
}

func (um *UpdatMode) toString() string {
	v, in := updatModeStrMap[*um]
	if !in {
		return ""
	} else {
		return v
	}
}
