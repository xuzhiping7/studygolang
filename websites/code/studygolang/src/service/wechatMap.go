package service

import (
	"logger"
	"math/rand"
	"model"
)

var Map_MapData map[int]*model.WechatMap

func init() {
	Map_MapData = make(map[int]*model.WechatMap)
	Map_MapData[0] = model.NewWechatMap(0, "林风角酒馆", "吵杂的林风村的小酒馆，你可以向酒侍询问消息，可以跟玩家组队。", 0, []int{}, []int{}, []int{}, []int{})
	Map_MapData[1] = model.NewWechatMap(1, "林风村", "林风角的一个小村落，常年大风，在海角之上临海而拔起，空气潮湿。春秋之季常遭海怪袭击。\n\n'买'：购买物品\n'卖'：出售物品", 0, []int{}, []int{}, []int{2, 3}, []int{1})
	Map_MapData[2] = model.NewWechatMap(2, "林风南海岸", "海怪的聚集地，是初出勇者作为'修炼'的好地方。\n\n'修炼'：在此地修炼\n", 0, []int{0, 1}, []int{50, 50}, []int{1}, []int{})
	Map_MapData[3] = model.NewWechatMap(3, "流放平原", "多年来林城帝国流放重犯之地，妖魂遍野，也有不少勇士在此扎根修炼，但有去无回的传闻实在太多了。\n\n'修炼'：在此地修炼\n", 5, []int{1, 2}, []int{30, 70}, []int{1}, []int{})
	//logger.Debugln(map_MapData["2"])
}

//根据传入的地图索引，获取遇到的怪物索引
func GetMosterByMap(mapIndex int) (mosterIndex int) {
	mosterIndex = -1

	number := rand.Intn(100)
	//logger.Debugln(number)
	for i := 0; i < len(Map_MapData[mapIndex].MostersRate); i++ {
		if number < Map_MapData[mapIndex].MostersRate[i] {
			mosterIndex = Map_MapData[mapIndex].Mosters[i]
			break
		}
		number -= Map_MapData[mapIndex].MostersRate[i]
	}

	return mosterIndex
}

//传入地图索引，查看该地图是否能进行某操作
func CanMapPractice(mapIndex int, funcType int) (b bool) {

	if len(Map_MapData[mapIndex].Functions) > 0 {
		for _, b := range Map_MapData[mapIndex].Functions {
			if b == funcType {
				return true
			}
		}
	}

	//旧的判断方法
	//if len(Map_MapData[mapIndex].MostersRate) > 0 {
	//	return true
	//}

	return false
}

//检测该地图是否有售卖相应的物品
func CheckMapSellProps(mapIndex int, propName string) (propId int, b bool) {
	b = false
	propId = -1
	for _, v := range Map_MapData[mapIndex].SellItems {
		propInfo, ok := Map_PropsData[v]
		if ok {
			if propInfo.Name == propName {
				propId = propInfo.Id
				b = true
				break
			}
		} else {
			logger.Debugln("CheckMapSellProps : Not find prop in Map_MapData")
		}
	}

	return propId, b
}
