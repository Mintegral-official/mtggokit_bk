package streamer

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoStreamer struct {
	container  container.Container
	cfg        *MongoStreamerCfg
	hasInit    bool
	totoalNum  int64
	errorNum   int64
	curParser  DataParser
	client     *mongo.Client
	collection *mongo.Collection
	cursor     *mongo.Cursor
}

func NewMongoStreamer(mongoConfig *MongoStreamerCfg) (*MongoStreamer, error) {
	streamer := &MongoStreamer{
		cfg: mongoConfig,
	}

	ctx, _ := context.WithTimeout(context.TODO(), time.Duration(mongoConfig.ConnectTimeout)*time.Microsecond)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.URI))
	if err != nil {
		if mongoConfig.Logger != nil {
			mongoConfig.Logger.Warnf("mongo connect error, err=[%s]", err.Error())
		}
		return nil, err
	}
	streamer.client = client

	//if err = client.Ping(ctx, readpref.Primary()); err != nil {
	//	if mongoConfig.Logger != nil {
	//		mongoConfig.Logger.Warnf("mongo ping error, err=[%s]", err.Error())
	//	}
	//	return nil, err
	//}

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
	return ms.cursor.Next(context.Background())
}

func (ms *MongoStreamer) Next() (container.DataMode, container.MapKey, interface{}, error) {
	if ms.cursor == nil {
		ms.errorNum++
		return container.DataModeAdd, nil, nil, errors.New("cursor is nil")
	}
	if ms.cursor.Err() != nil {
		ms.errorNum++
		return container.DataModeAdd, nil, nil, errors.New(fmt.Sprintf("cursor is error[%s]", ms.cursor.Err().Error()))
	}
	m, k, v, e := ms.curParser.Parse(ms.cursor.Current, nil)
	if e != nil {
		ms.errorNum++
	}
	ms.totoalNum++
	return m, k, v, e
}

func (ms *MongoStreamer) UpdateData(ctx context.Context) error {
	if !ms.hasInit && ms.cfg.IsSync {
		if err := ms.loadBase(ctx); err != nil {
			ms.cfg.Logger.Warnf("streamer[%s] LoadBase error, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totoalNum, ms.errorNum, ms.cfg.UserData)
			return err
		} else {
			ms.cfg.Logger.Infof("streamer[%s] LoadBase succ, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totoalNum, ms.errorNum, ms.cfg.UserData)
			return err
		}
	}
	go func() {
		if !ms.hasInit {
			if err := ms.loadBase(ctx); err != nil {
				ms.cfg.Logger.Warnf("streamer[%s] LoadBase error, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totoalNum, ms.errorNum, ms.cfg.UserData)
			} else {
				ms.cfg.Logger.Infof("streamer[%s] LoadBase succ, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totoalNum, ms.errorNum, ms.cfg.UserData)
			}
		}
		for {
			inc := time.After(time.Duration(ms.cfg.IncInterval) * time.Second)
			select {
			case <-ctx.Done():
				ms.cfg.Logger.Infof("streamer[%s] LoadInc succ, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totoalNum, ms.errorNum, ms.cfg.UserData)
				return
			case <-inc:
				if err := ms.loadInc(ctx); err != nil {
					ms.cfg.Logger.Warnf("streamer[%s] LoadInc error[%s], totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, err.Error(), ms.totoalNum, ms.errorNum, ms.cfg.UserData)
				}
			}
		}
	}()
	return nil
}

func (ms *MongoStreamer) loadBase(ctx context.Context) error {
	ms.totoalNum = 0
	ms.errorNum = 0
	c, _ := context.WithTimeout(ctx, time.Duration(ms.cfg.ReadTimeout)*time.Microsecond)
	cur, err := ms.collection.Find(c, ms.cfg.BaseQuery)
	if err != nil {
		ms.cfg.Logger.Warnf("streamer[%s]: loadBase error[%s]", ms.cfg.Name, err.Error())
		return err
	}

	if ms.cursor != nil {
		_ = ms.cursor.Close(ctx)
	}
	ms.cursor = cur
	err = ms.container.LoadBase(ms)
	if err != nil {
		ms.cfg.Logger.Warnf("Loadbase error, totalNum[%d], errorNum[%d], userData[%s]", ms.totoalNum, ms.errorNum, ms.cfg.UserData)
		return err
	}
	ms.cfg.Logger.Infof("Loadbase succ, totalNum[%d], errorNum[%d], userData[%s]", ms.totoalNum, ms.errorNum, ms.cfg.UserData)
	return nil
}

func (ms *MongoStreamer) loadInc(ctx context.Context) error {
	c, _ := context.WithTimeout(ctx, time.Duration(ms.cfg.ReadTimeout)*time.Microsecond)
	cur, err := ms.collection.Find(c, ms.cfg.IncQuery)
	if err != nil {
		ms.cfg.Logger.Warnf("streamer[%s]: loadInc error[%s]", ms.cfg.Name, err.Error())
		return err
	}

	if ms.cursor != nil {
		_ = ms.cursor.Close(ctx)
	}
	ms.cursor = cur
	return ms.container.LoadInc(ms)
}
