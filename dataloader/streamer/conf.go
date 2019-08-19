// ******************************************************************************* //
//                                                                                 //
//                                                                                 //
// File: conf.go                                                                   //
//                                                                                 //
// By: wangjia <jia.wang@mitegral.com>                                             //
//                                                                                 //
// Created: 2019/08/15 12:22:52 by wangjia                                         //
// Updated: 2019/08/15 12:22:52 by wangjia                                         //
//                                                                                 //
// ******************************************************************************* //

package streamer

type StreamerConfig struct {
	StreamerCfg *StreamerCfg `toml:"dataloader"`
}

type StreamerCfg struct {
	FileStreamer  []FileStreamerCfg  `toml:"file_streamer"`
	MongoStreamer []MongoStreamerCfg `toml:"mongo_streamer"`
}

type FileStreamerCfg struct {
	Name   string `toml:"name"`
	Path   string `toml:"path"`
	Mode   string `toml:"mode"`
	Parser string `toml:"parser"`
}

type MongoStreamerCfg struct {
	Mongo      string `toml:"mongo"`
	Timeout    int    `toml:"timeout"`
	Name       string `toml:"name"`
	Db         string `toml:"db"`
	Collection string `toml:"collection"`
	BaseQuery  string `toml:"base_query"`
	IncQuery   string `toml:"inc_query"`
	Parser     string `toml:"parser"`
}