package streamer

type FileStreamerCfg struct {
	Name       string
	Path       string
	UpdatMode  UpdatMode
	Interval   int
	IsSync     bool
	DataParser DataParser
	UserData   interface{}
}
