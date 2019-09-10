package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

type UserData struct {
	Uptime int64
}

type CampaignInfo struct {
	CampaignId   int64  `bson:"campaignId,omitempty"`
	AdvertiserId *int32 `bson:"advertiserId,omitempty"`
	Uptime       int64
}

type CampaignParser struct {
}

func (cp *CampaignParser) Parse(data []byte, userData interface{}) (container.DataMode, container.MapKey, interface{}, error) {
	ud, ok := userData.(*UserData)
	if !ok {
		return container.DataModeAdd, nil, nil, errors.New("user data parse error")
	}
	campaign := &CampaignInfo{}

	if err := bson.Unmarshal(data, &campaign); err != nil {
		fmt.Println("bson.Unmarsnal error:" + err.Error())
	}
	ud.Uptime = campaign.Uptime

	return container.DataModeAdd, container.I64Key(campaign.CampaignId), &campaign, nil
}

func main() {
	ms, err := streamer.NewMongoStreamer(&streamer.MongoStreamerCfg{
		Name:           "mongo_test",
		UpdatMode:      streamer.Dynamic,
		IncInterval:    60,
		IsSync:         true,
		URI:            "mongodb://127.0.0.1:21017",
		ConnectTimeout: 100,
		ReadTimeout:    20,
		BaseParser:     &CampaignParser{},
		IncParser:      &CampaignParser{},
		BaseQuery:      bson.M{},
		IncQuery:       bson.M{},
		UserData:       &UserData{},
		OnIncFinish: func(userData interface{}) interface{} {
			return "nfew inc base query"
		},
	})

	if err != nil {
		fmt.Println("Init mongo streamer error! err=" + err.Error())
		os.Exit(1)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = ms.UpdateData(ctx)
}
