package streamer

type MongoStreamer struct {
}

func NewMongoStreamer(mongoConfig *MongoStreamerCfg) *MongoStreamer {
	return &MongoStreamer{}
}
