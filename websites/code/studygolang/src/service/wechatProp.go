package service

import (
	"logger"
	"model"
	"strconv"
)

var Map_PropsData map[int]*model.WechatProp

func init() {
	Map_PropsData = make(map[int]*model.WechatProp)
	Map_PropsData[1] = model.NewWechatProp(1, "止血草", "使用能够回复生命值20点。", 5, 10, model.PropType_恢复生命值, 20)
	Map_PropsData[2] = model.NewWechatProp(2, "风信子", "一级素材，可以卖给商人。", 10, 15, model.PropType_没有任何作用, 0)
	Map_PropsData[3] = model.NewWechatProp(3, "结实黄泥", "一级素材，可以卖给商人。", 10, 15, model.PropType_没有任何作用, 0)
	Map_PropsData[4] = model.NewWechatProp(4, "清魂酒", "使用能够回复行动力20点。", 100, 150, model.PropType_恢复行动力, 20)
	Map_PropsData[5] = model.NewWechatProp(5, "隐世符", "重新角色命名。", 888, 1000, model.PropType_角色昵称更改, 0)
	//logger.Debugln(map_MapData["2"])
}

// 查看该表是否存在
func CheckPlayerPropExists(player_id int, prop_id int) int {
	prop := model.NewWechatPlayerProp()
	if err := prop.Where("player_id=" + strconv.Itoa(player_id) + " and prop_id=" + strconv.Itoa(prop_id)).Find("id"); err != nil {
		logger.Errorln("service CheckPlayerPropExists error:", err)
		return -1
	}
	if prop.Id != 0 {
		return prop.Id
	}
	return -1
}
