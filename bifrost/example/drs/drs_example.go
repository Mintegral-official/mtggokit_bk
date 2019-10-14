package main

import (
	"context"
	"fmt"
	"github.com/Mintegral-official/mtggokit/bifrost"
	"github.com/Mintegral-official/mtggokit/bifrost/container"
	"github.com/Mintegral-official/mtggokit/bifrost/streamer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"time"
)

type UserData struct {
	Bifrost             *bifrost.Bifrost
	PackageName         []string
	CampaignUptime      int64
	CreativeUptime      int64
	AuditCreativeUptime int64
}

type CampaignIdsParser struct {
}

func (*CampaignIdsParser) Parse(data []byte, userData interface{}) []streamer.ParserResult {
	s := string(data)
	items := strings.SplitN(s, "\t", 3)
	if len(items) != 3 {
		return nil
	}

	campaignIdsStr := strings.Split(items[1], ",")
	campaignIds := make([]int64, 0, len(campaignIdsStr))
	for _, v := range campaignIdsStr {
		if id, e := strconv.ParseInt(v, 10, 64); e == nil {
			campaignIds = append(campaignIds, id)
		}
	}
	return []streamer.ParserResult{{container.DataModeAdd, container.StrKey(items[0]), campaignIds, nil}}
}

type CampaignInfo struct {
	CampaignId   int64  `bson:"campaignId,omitempty"`
	AdvertiserId *int32 `bson:"advertiserId,omitempty"`
	PackageName  string `bson:"packageName,omitempty"`
	Uptime       int64  `bson:"updated,omitempty"`
}

type CreativeInfo struct {
	CampaignId int64 `bson:"campaignId,omitempty"`
	CreativeId int64
	Uptime     int64 `bson:"updated,omitempty"`
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
	if ud.CampaignUptime < campaign.Uptime {
		ud.CampaignUptime = campaign.Uptime
	}
	ud.PackageName = append(ud.PackageName, campaign.PackageName)
	return []streamer.ParserResult{{container.DataModeAdd, container.I64Key(campaign.CampaignId), &campaign, nil}}
}

func GetCampaigns(data *UserData) []int64 {
	d, err := data.Bifrost.Get("campaignsIds", container.StrKey("CampaignList"))
	if err != nil {
		return nil
	}
	campaignIds, ok := d.([]int64)
	if !ok {
		return nil
	}
	return campaignIds
}

type CreativeParser struct {
}

func (*CreativeParser) Parse(data []byte, userData interface{}) []streamer.ParserResult {
	ud, ok := userData.(*UserData)
	if !ok {
		return nil
	}
	creative := &CreativeInfo{}

	if err := bson.Unmarshal(data, &creative); err != nil {
		fmt.Println("bson.Unmarsnal error:" + err.Error())
	}
	if ud.CreativeUptime < creative.Uptime {
		ud.CreativeUptime = creative.Uptime
	}
	return []streamer.ParserResult{{container.DataModeAdd, container.I64Key(creative.CampaignId), &creative, nil}}
}

type AdxAuditCreativeInfo struct {
	CampaignId  int64  `bson:"campaignId,omitempty" json:"campaignId"`
	CountryCode string `bson:"countryCode,omitempty" json:"countryCode"`
	PackageName string `bson:"packageName,omitempty" json:"packageName"`
}
type AdxAuditCreativeParser struct {
}

func (*AdxAuditCreativeParser) Parse(data []byte, userData interface{}) []streamer.ParserResult {
	ud, ok := userData.(*UserData)
	if !ok {
		return []streamer.ParserResult{{container.DataModeAdd, nil, nil, errors.New("user data parse error")}}
	}
	creative := &AdxAuditCreativeInfo{}

	if err := bson.Unmarshal(data, &creative); err != nil {
		fmt.Println("bson.Unmarsnal error:" + err.Error())
	}
	if ud.AuditCreativeUptime < ud.CreativeUptime {
		ud.AuditCreativeUptime = ud.CreativeUptime
	}
	return []streamer.ParserResult{{container.DataModeAdd, container.I64Key(creative.CampaignId), creative, nil}}
}

func getCampaigIdsStreamer() streamer.Streamer {
	lfs := streamer.NewFileStreamer(&streamer.LocalFileStreamerCfg{
		Name:       "campaignsIds",
		Path:       "bifrost/data/caimpaigns.txt",
		UpdatMode:  streamer.Dynamic,
		Interval:   5,
		IsSync:     true,
		DataParser: &CampaignIdsParser{},
		Logger:     logrus.New(),
	})
	if lfs == nil {
		fmt.Println("Init local file streamer error!")
		return nil
	}
	lfs.SetContainer(&container.BufferedMapContainer{
		Tolerate: 0.5,
	})

	if err := lfs.UpdateData(context.Background()); err != nil {
		fmt.Println("streamer campaignsIds udpateData error: ", err.Error())
	}
	return lfs
}

func getCampaignInfoStreamer(bf *bifrost.Bifrost, ud *UserData) streamer.Streamer {
	// 创建 campaignInfo Streamer
	ms, err := streamer.NewMongoStreamer(&streamer.MongoStreamerCfg{
		Name:           "campaignsInfo",
		UpdatMode:      streamer.Dynamic,
		IncInterval:    5,
		IsSync:         true,
		URI:            "mongodb://13.250.108.190:27017",
		DB:             "new_adn",
		Collection:     "campaign",
		ConnectTimeout: 1000000,
		ReadTimeout:    2000000,
		BaseParser:     &CampaignParser{},
		IncParser:      &CampaignParser{},
		UserData:       ud,
		Logger:         logrus.New(),
		OnBeforeBase: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			campaignIds := GetCampaigns(ud)
			if campaignIds == nil {
				return nil
			}
			baseQuery := bson.M{"campaignId": bson.M{"$in": campaignIds}, "publisherId": 0, "status": 1, "system": 5}
			return baseQuery
		},
		OnBeforeInc: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			campaignIds := GetCampaigns(ud)
			if campaignIds == nil {
				return nil
			}
			incQuery := bson.M{
				"campaignId": bson.M{"$in": campaignIds}, "publisherId": 0, "status": 1, "system": 5,
				"updated": bson.M{"$gte": ud.CampaignUptime - 5, "$lte": int(time.Now().Unix())},
			}
			return incQuery
		},
	})
	if err != nil {
		fmt.Println("streamer init err, error:", err.Error())
	}
	if ms == nil {
		fmt.Println("streamer init err")
		return nil
	}
	ms.SetContainer(container.CreateBlockingMapContainer(100, 0))
	if err := ms.UpdateData(context.Background()); err != nil {
		fmt.Println("CampaignInfoStreamer updateData error: ", err.Error())
	}
	return ms
}

func getCreativeStreamer(bf *bifrost.Bifrost, ud *UserData) streamer.Streamer {
	// 创建 creative Streamer
	ms, err := streamer.NewMongoStreamer(&streamer.MongoStreamerCfg{
		Name:           "creativeInfo",
		UpdatMode:      streamer.Dynamic,
		IncInterval:    5,
		IsSync:         true,
		URI:            "mongodb://13.250.108.190:27017",
		DB:             "new_adn",
		Collection:     "creative",
		ConnectTimeout: 100000,
		ReadTimeout:    200000,
		BaseParser:     &CreativeParser{},
		IncParser:      &CreativeParser{},
		UserData:       ud,
		Logger:         logrus.New(),
		OnBeforeBase: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			campaignIds := GetCampaigns(ud)
			if campaignIds == nil {
				return nil
			}
			incQuery := bson.M{"$or": []bson.M{
				{"campaignId": bson.M{"$in": campaignIds}},
				{"packageName": bson.M{"$in": ud.PackageName}},
			}, "status": 1}
			return incQuery
		},
		OnBeforeInc: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			campaignIds := GetCampaigns(ud)
			if campaignIds == nil {
				return nil
			}
			incQuery := bson.M{"$or": []bson.M{
				{"campaignId": bson.M{"$in": campaignIds}},
				{"packageName": bson.M{"$in": ud.PackageName}},
			}, "updated": bson.M{"$gte": ud.CampaignUptime - 5, "$lte": int(time.Now().Unix())},
			}
			return incQuery
		},
	})
	if err != nil {
		fmt.Println("streamer init err, error:", err.Error())
	}
	if ms == nil {
		fmt.Println("streamer init err")
		return nil
	}
	ms.SetContainer(container.CreateBlockingMapContainer(100, 0))
	if err := ms.UpdateData(context.Background()); err != nil {
		fmt.Println("Creative Streamer updateData error: ", err.Error())
	}
	return ms
}

func getAdxAuditCreativeStreamer(ud *UserData) streamer.Streamer {
	// 创建 creative Streamer
	ms, err := streamer.NewMongoStreamer(&streamer.MongoStreamerCfg{
		Name:           "adxAuditCreativeInfo",
		UpdatMode:      streamer.Dynamic,
		IncInterval:    5,
		IsSync:         true,
		URI:            "mongodb://47.252.0.66:27050",
		DB:             "dsp_audit",
		Collection:     "group_creative_audit",
		ConnectTimeout: 100000000,
		ReadTimeout:    2000000,
		BaseParser:     &AdxAuditCreativeParser{},
		IncParser:      &AdxAuditCreativeParser{},
		UserData:       ud,
		Logger:         logrus.New(),
		OnBeforeBase: func(userData interface{}) interface{} {
			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			campaignIds := GetCampaigns(ud)
			if campaignIds == nil {
				return nil
			}
			incQuery := bson.M{"campaignId": bson.M{"$in": campaignIds}, "status": 1}
			return incQuery
		},
		OnBeforeInc: func(userData interface{}) interface{} {

			ud, ok := userData.(*UserData)
			if !ok {
				return nil
			}
			campaignIds := GetCampaigns(ud)
			if campaignIds == nil {
				return nil
			}
			incQuery := bson.M{"campaignId": bson.M{"$in": campaignIds}, "updated": bson.M{"$gte": ud.AuditCreativeUptime - 5, "$lte": int(time.Now().Unix())}}
			return incQuery
		},
	})
	if err != nil {
		fmt.Println("streamer init err, error:", err.Error())
	}
	if ms == nil {
		fmt.Println("streamer init err")
		return nil
	}
	ms.SetContainer(container.CreateBlockingMapContainer(100, 0))
	if err := ms.UpdateData(context.Background()); err != nil {
		fmt.Println("Creative Streamer updateData error: ", err.Error())
	}
	return ms
}

func run() {
	// 初始化 Bifrost
	bf := bifrost.NewBifrost()

	ud := &UserData{
		Bifrost: bf,
	}

	// 创建 campaignsIds streamer
	idStreamer := getCampaigIdsStreamer()
	if err := bf.Register("campaignsIds", idStreamer); err != nil {
		fmt.Println("bf.Register campaignsIds error ")
	}

	// 创建 campaignInfo Streamer
	infoStreamer := getCampaignInfoStreamer(bf, ud)
	if err := bf.Register("campaignsInfo", infoStreamer); err != nil {
		fmt.Println("bf.Register campaignsInfo error ")
	}

	// 创建 creative Streamer
	creativeStreamer := getCreativeStreamer(bf, ud)
	if err := bf.Register("creativeInfo", creativeStreamer); err != nil {
		fmt.Println("bf.Register creativeInfo error ")
	}

	// 创建 creative Streamer
	auditCreativeStream := getAdxAuditCreativeStreamer(ud)
	if err := bf.Register("AuditCreativeInfo", auditCreativeStream); err != nil {
		fmt.Println("bf.Register creativeInfo error ")
	}

	// test
	data, err := bf.Get("campaignsIds", container.StrKey("CampaignList"))
	if err != nil {
		fmt.Println("Get [campaignsIds.CampaignList] error: " + err.Error())
	}
	campaignIds, ok := data.([]int64)
	if !ok {
		fmt.Println("transfer data to []int64 error")
	}

	fmt.Println("len: ", len(campaignIds))
}

func main() {
	run()
}
