package conf

import "github.com/Mintegral-official/mtggokit/bifrost/streamer"

type StreamerConfig struct {
	StreamerCfg *StreamerCfg `toml:"bifrost"`
}

type StreamerCfg struct {
	FileStreamer  []FileStreamerCfg  `toml:"file_streamer"`
	MongoStreamer []MongoStreamerCfg `toml:"mongo_streamer"`
}

type FileStreamerCfg struct {
	Name       string              `toml:"name"`
	Path       string              `toml:"path"`
	Mode       string              `toml:"mode"`
	Parser     string              `toml:"parser"`
	DataParser streamer.DataParser `toml:"-"`
}

type MongoStreamerCfg struct {
	Mongo      string              `toml:"mongo"`
	Timeout    int                 `toml:"timeout"`
	Name       string              `toml:"name"`
	Db         string              `toml:"db"`
	Collection string              `toml:"collection"`
	BaseQuery  string              `toml:"base_query"`
	IncQuery   string              `toml:"inc_query"`
	Parser     string              `toml:"parser"`
	DataParser streamer.DataParser `toml:"-"`
}
