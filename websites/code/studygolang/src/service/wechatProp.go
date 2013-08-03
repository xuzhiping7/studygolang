package service

import (
	//"logger"
	"model"
)

var Map_PropsData map[int]*model.WechatProp

func init() {
	Map_PropsData = make(map[int]*model.WechatProp)
	Map_PropsData[0] = model.NewWechatProp(0, "止血草", "使用能够回复生命值20点。")
	Map_PropsData[1] = model.NewWechatProp(1, "风信子", "一级素材，可以卖给商人。")
	Map_PropsData[2] = model.NewWechatProp(2, "结实黄泥", "一级素材，可以卖给商人。")

	//logger.Debugln(map_MapData["2"])
}
