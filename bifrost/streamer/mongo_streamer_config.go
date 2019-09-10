package streamer

import "github.com/Mintegral-official/mtggokit/bifrost/log"

type MongoStreamerCfg struct {
	Name           string
	UpdatMode      UpdatMode
	IncInterval    int
	IsSync         bool
	URI            string
	DB             string
	Collection     string
	ConnectTimeout int
	ReadTimeout    int
	BaseParser     DataParser
	IncParser      DataParser
	BaseQuery      interface{}
	IncQuery       interface{}
	UserData       interface{}
	OnIncFinish    func(interface{}) interface{}
	Logger         log.BiLogger
}
