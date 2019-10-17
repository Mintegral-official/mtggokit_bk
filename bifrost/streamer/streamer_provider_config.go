package streamer

import "github.com/Mintegral-official/mtggokit/bifrost/log"

type StreamerProviderCfg struct {
	Name       string
	ExpireTime int64
	Logger     log.BiLogger
}
