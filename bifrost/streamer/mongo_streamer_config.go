package streamer

type MongoStreamerCfg struct {
	Name        string
	UpdatMode   UpdatMode
	IncInterval int
	IsSync      bool
	IP          string
	Port        int
	BaseParser  DataParser
	IncParser   DataParser
	BaseQuery   interface{}
	IncQuery    interface{}
	UserData    interface{}
	OnIncFinish func(interface{}) interface{}
}
