package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"os/signal"
	"time"
)

type UserData struct {
	Uptime int64
}

type CampaignInfo struct {
	CampaignId   int64  `bson:"campaignId,omitempty"`
	AdvertiserId *int32 `bson:"advertiserId,omitempty"`
	Uptime       int64  `bson:"updated,omitempty"`
}

type CampaignParser struct {
}

func (cp *CampaignParser) Parse(data []byte, userData interface{}) []streamer.ParserResult {
	ud, ok := userData.(*UserData)
	if !ok {
		return nil
	}
	campaign := &CampaignInfo{}

	if err := bson.Unmarshal(data, &campaign); err != nil {
		fmt.Println("bson.Unmarsnal error:" + err.Error())
	}
	if ud.Uptime < campaign.Uptime {
		ud.Uptime = campaign.Uptime
	}
	return []streamer.ParserResult{{container.DataModeAdd, container.I64Key(campaign.CampaignId), &campaign, nil}}
}

func main() {
	ms, err := streamer.NewMongoStreamer(&streamer.MongoStreamerCfg{
		Name:           "mongo_test",
		UpdatMode:      streamer.Dynamic,
		IncInterval:    5,
		IsSync:         true,
		URI:            "mongodb://13.250.108.190:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 10000,
		ReadTimeout:    20000,
		BaseParser:     &CampaignParser{},
		IncParser:      &CampaignParser{},
		BaseQuery:      bson.M{"status": 1, "advertiserId": 903},
		IncQuery:       bson.M{"advertiserId": 903},
		UserData:       &UserData{},
		Logger:         logrus.New(),
		OnBeforeInc: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			incQuery := bson.M{"advertiserId": 903, "updated": bson.M{"$gte": ud.Uptime - 5, "$lte": int(time.Now().Unix())}}
			return incQuery
		},
	})
	if ms == nil {
		fmt.Println("streamer init err")
		return
	}
	ms.SetContainer(container.CreateBlockingMapContainer(100, 0))

	if err != nil {
		fmt.Println("Init mongo streamer error! err=" + err.Error())
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := make(chan os.Signal)
	signal.Notify(c)
	_ = ms.UpdateData(ctx)

	value, err := ms.GetContainer().Get(container.StrKey("abc"))
	if err == nil {
		fmt.Println(value)
	}

	s := <-c
	fmt.Println("退出信号", s)
}
