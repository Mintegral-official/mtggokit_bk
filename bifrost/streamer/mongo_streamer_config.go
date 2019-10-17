package streamer

import (
	"github.com/Mintegral-official/mtggokit/bifrost/log"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	FindOpt        *options.FindOptions
	OnBeforeBase   func(interface{}) interface{}
	OnBeforeInc    func(interface{}) interface{}
	Logger         log.BiLogger
}
