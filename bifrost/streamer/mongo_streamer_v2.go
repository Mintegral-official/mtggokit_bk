package streamer

//
//import (
//	"context"
//	"fmt"
//	"github.com/Mintegral-official/mtggokit/bifrost/container"
//	"github.com/pkg/errors"
//	mgo "gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	"strings"
//	"time"
//)
//
//type MongoStreamerV2 struct {
//	container     container.Container
//	cfg           *MongoStreamerCfg
//	hasInit       bool
//	totalNum     int64
//	errorNum      int64
//	curParser     DataParser
//	mgoSession    *mgo.Session
//	mgoCollection *mgo.Collection
//	query         *mgo.Query
//	iter          *mgo.Iter
//}
//
//func NewMongoStreamerV2(mongoConfig *MongoStreamerCfg) (*MongoStreamerV2, error) {
//	streamer := &MongoStreamerV2{
//		cfg: mongoConfig,
//	}
//	dialInfo := &mgo.DialInfo{
//		Addrs:     strings.Split(mongoConfig.URI, ","),
//		Direct:    true,
//		Timeout:   time.Duration(mongoConfig.ReadTimeout) * time.Second,
//		PoolLimit: 10,
//	}
//	session, err := mgo.DialWithInfo(dialInfo)
//	if session == nil || err != nil {
//		mongoConfig.Logger.Warn("connect mongo [", mongoConfig.URI, "] failed. ", err.Error())
//		return nil, errors.New("mgo.DialWithInfo failed." + mongoConfig.URI)
//	}
//	session.SetSocketTimeout(time.Duration(mongoConfig.ReadTimeout) * time.Second)
//	session.SetSyncTimeout(time.Duration(mongoConfig.ReadTimeout) * time.Second)
//	session.SetMode(mgo.Monotonic, true)
//	if streamer.mgoSession != nil {
//		streamer.mgoSession.Close()
//	}
//	streamer.mgoSession = session
//	db := session.DB(mongoConfig.DB)
//	streamer.mgoCollection = db.C(mongoConfig.Collection)
//	if streamer.mgoSession == nil || streamer.mgoCollection == nil {
//		return nil, errors.New("mongo connection is nil, " + mongoConfig.Collection)
//	}
//	mongoConfig.Logger.Infof("mongo conn succ, addr=[%s]", mongoConfig.URI)
//	return streamer, nil
//}
//
//func (ms *MongoStreamerV2) SetContainer(container container.Container) {
//	ms.container = container
//}
//
//func (ms *MongoStreamerV2) GetContainer() container.Container {
//	return ms.container
//}
//
//func (ms *MongoStreamerV2) GetSchedInfo() *SchedInfo {
//	return &SchedInfo{
//		TimeInterval: ms.cfg.IncInterval,
//	}
//}
//
//func (ms *MongoStreamerV2) HasNext() bool {
//	ms.iter.
//	return ms.iter.Next()
//}
//
//func (ms *MongoStreamerV2) Next() (container.DataMode, container.MapKey, interface{}, error) {
//	if ms.cursor == nil {
//		ms.errorNum++
//		return container.DataModeAdd, nil, nil, errors.New("cursor is nil")
//	}
//	if ms.cursor.Err() != nil {
//		ms.errorNum++
//		return container.DataModeAdd, nil, nil, errors.New(fmt.Sprintf("cursor is error[%s]", ms.cursor.Err().Error()))
//	}
//	m, k, v, e := ms.curParser.Parse(ms.cursor.Current, nil)
//	if e != nil {
//		ms.errorNum++
//	}
//	ms.totalNum++
//	return m, k, v, e
//}
//
//func (ms *MongoStreamerV2) UpdateData(ctx context.Context) error {
//	if ms.hasInit && ms.cfg.IsSync {
//		if err := ms.loadBase(ctx); err != nil {
//			ms.cfg.Logger.Warnf("streamer[%s] LoadBase error, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totalNum, ms.errorNum, ms.cfg.UserData)
//			return err
//		} else {
//			ms.cfg.Logger.Infof("streamer[%s] LoadBase succ, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totalNum, ms.errorNum, ms.cfg.UserData)
//			return err
//		}
//	}
//	go func() {
//		if ms.hasInit {
//			if err := ms.loadBase(ctx); err != nil {
//				ms.cfg.Logger.Warnf("streamer[%s] LoadBase error, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totalNum, ms.errorNum, ms.cfg.UserData)
//			} else {
//				ms.cfg.Logger.Infof("streamer[%s] LoadBase succ, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totalNum, ms.errorNum, ms.cfg.UserData)
//			}
//		}
//		for {
//			inc := time.After(time.Duration(ms.cfg.IncInterval) * time.Second)
//			select {
//			case <-ctx.Done():
//				ms.cfg.Logger.Infof("streamer[%s] LoadInc succ, totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, ms.totalNum, ms.errorNum, ms.cfg.UserData)
//				return
//			case <-inc:
//				if err := ms.loadInc(ctx); err != nil {
//					ms.cfg.Logger.Warnf("streamer[%s] LoadInc error[%s], totalNum[%d], errorNum[%d], userData[%s]", ms.cfg.Name, err.Error(), ms.totalNum, ms.errorNum, ms.cfg.UserData)
//				}
//			}
//		}
//	}()
//	return nil
//}
//
//func (ms *MongoStreamerV2) loadBase(ctx context.Context) error {
//	ms.totalNum = 0
//	ms.errorNum = 0
//	query := ms.mgoCollection.Find(ms.cfg.BaseQuery)
//	if query == nil {
//		ms.cfg.Logger.Warnf("streamer[%s]: read mongo error, base query[%s]", ms.cfg.Name, ms.cfg.BaseQuery)
//		return errors.New(fmt.Sprintf("streamer[%s]: read mongo error, base query[%s]", ms.cfg.Name, ms.cfg.BaseQuery))
//	}
//
//	ms.query = query
//	err := ms.container.LoadBase(ms)
//	if err != nil {
//		ms.cfg.Logger.Warnf("Loadbase error, totalNum[%d], errorNum[%d], userData[%s]", ms.totalNum, ms.errorNum, ms.cfg.UserData)
//		return err
//	}
//	ms.cfg.Logger.Infof("Loadbase succ, totalNum[%d], errorNum[%d], userData[%s]", ms.totalNum, ms.errorNum, ms.cfg.UserData)
//	return nil
//}
//
//func (ms *MongoStreamerV2) loadInc(ctx context.Context) error {
//	query := ms.mgoCollection.Find(ms.cfg.BaseQuery)
//	if query != nil {
//		ms.cfg.Logger.Warnf("read mongo error, inc query[%s]", ms.cfg.IncQuery)
//		return errors.New("streamer[" + ms.cfg.Name + "]: read mongo error")
//	}
//
//	ms.iter = query.Iter()
//
//	ms.query = query
//	return ms.container.LoadInc(ms)
//}
