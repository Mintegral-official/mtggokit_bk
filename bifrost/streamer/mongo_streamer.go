package streamer

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoStreamer struct {
	container  container.Container
	cfg        *MongoStreamerCfg
	hasInit    bool
	totalNum   int64
	errorNum   int64
	curParser  DataParser
	client     *mongo.Client
	collection *mongo.Collection
	cursor     *mongo.Cursor
	result     []ParserResult
	curLen     int
	findOpt    *options.FindOptions
	startTime  int64
	endTime    int64
}

func NewMongoStreamer(mongoConfig *MongoStreamerCfg) (*MongoStreamer, error) {
	streamer := &MongoStreamer{
		cfg: mongoConfig,
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Duration(mongoConfig.ConnectTimeout)*time.Microsecond)
	opt := options.Client().ApplyURI(mongoConfig.URI)
	opt.SetReadPreference(readpref.SecondaryPreferred())
	direct := true
	opt.Direct = &direct
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		if mongoConfig.Logger != nil {
			mongoConfig.Logger.Warnf("mongo connect error, err=[%s]", err.Error())
		}
		return nil, err
	}
	streamer.client = client

	streamer.findOpt = options.MergeFindOptions(mongoConfig.FindOpt)
	d := time.Duration(mongoConfig.ReadTimeout) * time.Microsecond
	streamer.findOpt.MaxTime = &d

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		if mongoConfig.Logger != nil {
			mongoConfig.Logger.Warnf("mongo ping error, err=[%s]", err.Error())
		}
		return nil, err
	}

	streamer.collection = client.Database(mongoConfig.DB).Collection(mongoConfig.Collection)
	if streamer.collection == nil {
		if mongoConfig.Logger != nil {
			mongoConfig.Logger.Warnf("[%s.%s] Not found", mongoConfig.DB, mongoConfig.Collection)
		}
		return nil, errors.New(fmt.Sprintf("[%s.%s] Not found", mongoConfig.DB, mongoConfig.Collection))
	}

	return streamer, nil
}

func (ms *MongoStreamer) SetContainer(container container.Container) {
	ms.container = container
}

func (ms *MongoStreamer) GetContainer() container.Container {
	return ms.container
}

func (ms *MongoStreamer) GetSchedInfo() *SchedInfo {
	return &SchedInfo{
		TimeInterval: ms.cfg.IncInterval,
	}
}

func (ms *MongoStreamer) HasNext() bool {

	return ms.curLen < len(ms.result) || ms.cursor.Next(context.Background())
}

func (ms *MongoStreamer) Next() (container.DataMode, container.MapKey, interface{}, error) {

	ms.totalNum++
	if ms.curLen < len(ms.result) {
		r := ms.result[ms.curLen]
		ms.curLen++
		if r.Err != nil {
			ms.errorNum++
		}
		return r.DataMode, r.Key, r.Value, r.Err
	}

	if ms.cursor == nil {
		ms.errorNum++
		return container.DataModeAdd, nil, nil, errors.New("cursor is nil")
	}
	if ms.cursor.Err() != nil {
		ms.WarnStatus(fmt.Sprintf("cursor is error[%s]", ms.cursor.Err().Error()))
		ms.errorNum++
		return container.DataModeAdd, nil, nil, errors.New(fmt.Sprintf("cursor is error[%s]", ms.cursor.Err().Error()))
	}
	result := ms.curParser.Parse(ms.cursor.Current, ms.cfg.UserData)
	if result == nil {
		ms.errorNum++
		return container.DataModeAdd, nil, nil, errors.New("Parse error")
	}
	ms.curLen = 0
	ms.result = result
	if ms.curLen < len(ms.result) {
		r := ms.result[ms.curLen]
		ms.curLen++
		if r.Err != nil {
			ms.errorNum++
		}
		return r.DataMode, r.Key, r.Value, r.Err
	}
	ms.errorNum++
	return container.DataModeAdd, nil, nil, errors.New(fmt.Sprintf("Index[%d] error, len[%d]", ms.curLen, len(ms.result)))
}

func (ms *MongoStreamer) UpdateData(ctx context.Context) error {
	ms.startTime = time.Now().UnixNano()
	if !ms.hasInit && ms.cfg.IsSync {
		err := ms.loadBase(ctx)
		ms.endTime = time.Now().UnixNano()
		if err != nil {
			ms.WarnStatus("LoadBase error:" + err.Error())
			return err
		}
		ms.InfoStatus("LoadBase Succ")
		ms.hasInit = true

	}
	go func() {
		ms.startTime = time.Now().UnixNano()
		if !ms.hasInit {
			err := ms.loadBase(ctx)
			ms.endTime = time.Now().UnixNano()
			if err != nil {
				ms.WarnStatus("LoadBase error:" + err.Error())
			} else {
				ms.InfoStatus("LoadBase succ")
			}
		}
		for {
			inc := time.After(time.Duration(ms.cfg.IncInterval) * time.Second)
			select {
			case <-ctx.Done():
				ms.startTime = time.Now().UnixNano()
				ms.endTime = time.Now().UnixNano()
				ms.InfoStatus("LoadInc Finish:")
				return
			case <-inc:
				ms.startTime = time.Now().UnixNano()
				err := ms.loadInc(ctx)
				ms.endTime = time.Now().UnixNano()
				if err != nil {
					ms.WarnStatus("LoadInc Error:" + err.Error())
				} else {
					ms.InfoStatus("LoadInc Succ:")
				}
			}
		}
	}()
	return nil
}

func (ms *MongoStreamer) loadBase(context.Context) error {
	ms.totalNum = 0
	ms.errorNum = 0
	if ms.cfg.OnBeforeBase != nil {
		ms.cfg.BaseQuery = ms.cfg.OnBeforeBase(ms.cfg.UserData)
	}
	cur, err := ms.collection.Find(nil, ms.cfg.BaseQuery, ms.findOpt)
	if err != nil {
		return errors.Wrap(err, "FindError")
	}

	if ms.cursor != nil {
		_ = ms.cursor.Close(nil)
	}
	ms.cursor = cur
	ms.curParser = ms.cfg.BaseParser
	err = ms.container.LoadBase(ms)
	return err
}

func (ms *MongoStreamer) loadInc(ctx context.Context) error {
	if ms.cfg.OnBeforeInc != nil {
		ms.cfg.IncQuery = ms.cfg.OnBeforeInc(ms.cfg.UserData)
	}
	c, _ := context.WithTimeout(ctx, time.Duration(ms.cfg.ReadTimeout)*time.Microsecond)
	cur, err := ms.collection.Find(nil, ms.cfg.IncQuery, ms.cfg.FindOpt)
	if err != nil {
		return errors.Wrap(err, "FindError")
	}
	if ms.cursor != nil {
		_ = ms.cursor.Close(c)
	}
	ms.cursor = cur
	ms.curParser = ms.cfg.IncParser
	return ms.container.LoadInc(ms)
}

func (ms *MongoStreamer) InfoStatus(s string) {
	if ms.cfg.Logger != nil {
		ms.cfg.Logger.Infof("streamer[%s] %s, totalNum[%d], errorNum[%d], timeUsed[%d]", ms.cfg.Name, s, ms.totalNum, ms.errorNum, (ms.endTime-ms.startTime)/10e6)
	}
}

func (ms *MongoStreamer) WarnStatus(s string) {
	if ms.cfg.Logger != nil {
		ms.cfg.Logger.Warnf("streamer[%s] %s, totalNum[%d], errorNum[%d], timeUsed[%d]", ms.cfg.Name, s, ms.totalNum, ms.errorNum, (ms.endTime-ms.startTime)/10e6)
	}
}
