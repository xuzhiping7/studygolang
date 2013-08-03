package service

import (
	//"logger"
	"model"
)

var Map_MapData map[int]*model.WechatMap

func init() {
	Map_MapData = make(map[int]*model.WechatMap)
	Map_MapData[0] = model.NewWechatMap(0, "林风角酒馆", []int{})
	Map_MapData[1] = model.NewWechatMap(1, "林风村", []int{})
	Map_MapData[2] = model.NewWechatMap(2, "林风南海岸", []int{0, 1})

	//logger.Debugln(map_MapData["2"])
}
