package service

import (
	//"logger"
	"math/rand"
	"model"
)

var Map_MapData map[int]*model.WechatMap

func init() {
	Map_MapData = make(map[int]*model.WechatMap)
	Map_MapData[0] = model.NewWechatMap(0, "林风角酒馆", "吵杂的林风村的小酒馆，你可以向酒侍询问消息，可以跟玩家组队。", []int{}, []int{})
	Map_MapData[1] = model.NewWechatMap(1, "林风村", "林风角的一个小村落，常年大风，在海角之上临海而拔起，空气潮湿。春秋之季常遭海怪袭击。", []int{}, []int{})
	Map_MapData[2] = model.NewWechatMap(2, "林风南海岸", "海怪的聚集地，是个'修炼'的好地方。", []int{0, 1}, []int{30, 70})

	//logger.Debugln(map_MapData["2"])
}

//根据传入的地图索引，获取遇到的怪物索引
func GetMosterByMap(mapIndex int) (mosterIndex int) {
	mosterIndex = -1

	number := rand.Intn(100)

	for i := 0; i < len(Map_MapData[mapIndex].MostersRate); i++ {
		if number < Map_MapData[mapIndex].MostersRate[i] {
			mosterIndex = Map_MapData[mapIndex].Mosters[i]
		}
		number -= Map_MapData[mapIndex].MostersRate[i]
	}

	return mosterIndex
}

//传入地图索引，查看该地图是否能进行修炼
func CanMapPractice(mapIndex int) (b bool) {
	if len(Map_MapData[mapIndex].MostersRate) > 0 {
		return true
	}
	return false
}
