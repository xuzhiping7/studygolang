// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"config"
	"csv_controller"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"logger"
	"model"
	"strconv"
	"strings"
	"time"
)

//储存玩家信息，24小时后删除
var map_PlayerData map[string]*model.WechatPlayer

//储存所有对话模板
var textTemplate map[string]string

//储存所有命令前缀
var commandPrefix map[int]string

//储存地图ID对应名
//var map_MapName map[int]string

//定义枚举事件常量
const (
	flag_注册完成 = iota
	flag_用户传入角色名申请更名操作
	flag_暂无
)

func init() {

	//初始化玩家信息MAP
	map_PlayerData = make(map[string]*model.WechatPlayer)
	//初始化所有对话模板
	textTemplate = make(map[string]string)

	//从json读入数据
	/*
		date, err := ioutil.ReadFile(config.ROOT + "/conf/wechatTextTemplate.json")
		if err != nil {
			logger.Errorln("Read wechatTextTemplate.json fail error:", err)
		}
		if err2 := json.Unmarshal(date, &textTemplate); err2 != nil {
			logger.Errorln("Unmarshal wechatTextTemplate.json fail error:", err2)
		}
	*/
	//将数据写入csv
	/*
		testData := make([][]string, len(textTemplate))
		flag := 0
		for k, v := range textTemplate {
			testData[flag] = []string{k, v}
			flag++
		}
		csv_controller.WriteCSV(config.ROOT+"/conf/wechatTextTemplate.csv", []string{"编号", "内容"}, testData)
	*/

	//从csv读取数据
	testData := csv_controller.ReadCSV(config.ROOT + "/conf/wechatTextTemplate.csv")
	for i := 1; i < len(testData); i++ {
		textTemplate[testData[i][0]] = testData[i][1]
	}
	//logger.Debugln(textTemplate)

	//将数据写入json
	/*
		textTemplate["0"] = "注册成功!\n\n感谢注册微信奇幻网游《传说》，目前游戏处于删档测试阶段，有什么问题和建议请联系邮箱xuzhiping7@qq.com，希望您能够享受和喜欢这个世界。\n\n（请输出‘传说’两个字，开启您的游戏旅程！）"
		textTemplate["1"] = "请输出‘我’查看您的最新状态。输入‘帮助’可以查看相关操作指令。"
		textTemplate["2"] = "注册"
		textTemplate["3"] = "欢迎来到微信奇幻网游《传说》，请输出'注册'，确认注册游戏。"
		textTemplate["4"] = "创建角色中,请输入您的角色名。(例如‘德玛西亚皇子’，8个汉字内。)"
		textTemplate["5"] = "角色【%s】成功创建！请输入'帮助'开始游戏简单的指导。"
		textTemplate["6"] = "%s\n当前地点:%s\n\n等级:%d\n职业:%s\n称号:%s\n状态：%s\n\n行动力:%d\n声望:%d\n\n生命:%s\n负重:%s\n抗性:%s\n\n攻击:%d   防御:%d\n体力:%d   敏捷:%d\n剩余分配点数:%d\n\n(您可以输入'当前'查看你所能做的事情)"
		textTemplate["7"] = "抱歉，你不能到达这个地方。"
		textTemplate["8"] = "你到达了%s"
		textTemplate["9"] = "前往%s\n"
		textTemplate["10"] = "你当前可以前往的地方:\n"
		textTemplate["11"] = "您可以输入相应关键字来实现的操作：\n‘我’----查看角色状态\n‘当前’----查看事件状态\n‘前往’----查看你可以去的地方\n‘前往’+ ‘地名’----命令角色前往目标地方\n‘修炼’----在角色当前的地方修炼\n‘搜寻’----搜寻当前地方的物品或事件"
		textTemplate["100"] = "网络错误，请重新输入。"

		jsonData, _ := json.Marshal(textTemplate)
		ioutil.WriteFile(config.ROOT+"/conf/wechatTextTemplate2.json", jsonData, 0666)
	*/
	commandPrefix = make(map[int]string)
	commandPrefix[0] = "我"
	commandPrefix[1] = "当前"
	commandPrefix[2] = "前往"
	commandPrefix[3] = "修炼"
	commandPrefix[4] = "状态"
	commandPrefix[5] = "搜寻"
	commandPrefix[6] = "帮助"
	commandPrefix[7] = "传说"
	commandPrefix[8] = "装备"
	commandPrefix[9] = "道具"
	commandPrefix[10] = "使用"
	commandPrefix[11] = "提升"
	commandPrefix[12] = "提升攻击"
	commandPrefix[13] = "提升体力"
	commandPrefix[14] = "提升防御"
	commandPrefix[15] = "提升敏捷"
	commandPrefix[16] = "提升智慧"
	commandPrefix[17] = "附近"
	commandPrefix[18] = "休息"
	commandPrefix[19] = "买"
	commandPrefix[20] = "卖"
	commandPrefix[21] = "挑战"
	commandPrefix[22] = "确认洗点"
	//map_MapName = make(map[int]string)
	//map_MapName[0] = "林风角酒馆"
	//map_MapName[1] = "林风角"
	//map_MapName[2] = "林风南海岸"
}

// 创建一个wechat玩家
func CreateWechatPlayer(openid string) bool {
	player := model.NewWechatPlayer()

	player.OpenId = openid
	player.NickName = "EmptyNow"
	player.UserName = "EmptyNow"
	player.Exp = 0
	player.Mobility = 100
	player.Attack = 5
	player.Defense = 5
	player.Stamina = 5
	player.Agility = 5
	player.Wisdom = 5

	if _, err := player.Insert(); err != nil {
		logger.Errorln("player service CreateWechatPlayer error:", err)
		return false
	}
	return true
}

//获取一个wechat玩家信息
func GetWechatPlayer(openid string) (player *model.WechatPlayer) {

	player = model.NewWechatPlayer()
	err := player.Where("openid=" + openid).Find()
	if err != nil {
		logger.Errorln("player service GetWechatPlayer Error:", err)
		return
	}
	//logger.Debugln(player)

	//属性配点正确性验证
	if player.Attack+player.Agility+player.Wisdom+player.Stamina+player.Defense+player.NoDistribution > player.Level*3+25 {
		PlayerRedistributeAttribute(player)
		logger.Errorln("玩家数据异常，请查看玩家昵称：" + player.NickName + "  ID:" + strconv.Itoa(player.Id))
	}

	//初始化动态信息
	InitPlayerProp(player)

	player.Cur_Mobility = player.Mobility
	player.Cur_HP = player.Stamina * 10
	player.Max_HP = player.Stamina * 10
	player.Cur_Resistance = player.Stamina + player.Defense*3 + player.Wisdom*3
	player.Max_Burden = player.Stamina * 5
	player.Cur_Burden = 0

	return player
}

func InitPlayerProp(player *model.WechatPlayer) {

	props, err := model.NewWechatPlayerProp().Where("player_id=" + strconv.Itoa(player.Id)).FindAll()

	if err != nil {
		logger.Errorln("service InitPlayerProp error:", err)
	} else {
		player.Map_PlayerProp = make(map[int]*model.WechatPlayerProp)

		for _, prop := range props {
			player.Map_PlayerProp[prop.PropId] = prop
		}
	}

}

// 判断该OpenID是否已经被注册了
func OpenidExists(openid string) bool {
	player := model.NewWechatPlayer()
	if err := player.Where("openid=" + openid).Find("id"); err != nil {
		logger.Errorln("service OpenidExists error:", err)
		return false
	}
	if player.Id != 0 {
		return true
	}
	return false
}

func WechatResponseHandle(openid string, content string) (s_ReturnContent string) {

	/*
		假设不存在当前用户,让玩家进入注册流程
		存在当前用户，则读取当前用户的所有信息。
	*/
	player, ok := map_PlayerData[openid]

	if !ok {
		if !OpenidExists(openid) {

			//假设当前用户输入的不是注册
			if content == textTemplate["2"] {
				b_Reg := CreateWechatPlayer(openid)

				if b_Reg {
					s_ReturnContent = textTemplate["0"]
				}
			} else {
				s_ReturnContent = textTemplate["3"]
			}
			return
		}

		logger.Debugln("没有储存记录，从数据库读出，并且保存到MAP")
		player = GetWechatPlayer(openid)
		map_PlayerData[openid] = player
	}

	//logger.Debugln(player)

	switch player.Flag {

	case flag_注册完成:
		s_ReturnContent = textTemplate["4"]
		player.Flag = flag_用户传入角色名申请更名操作
		if err := player.UpdateFlag(); err != nil {
			logger.Errorln("wechat UpdateFlag Error:", err)
			s_ReturnContent = textTemplate["100"]
		}
	case flag_用户传入角色名申请更名操作:
		player.Flag = flag_用户传入角色名申请更名操作
		runes := []rune(content)
		if len(runes) > 8 || len(runes) < 2 {
			s_ReturnContent = textTemplate["4"]
		} else {
			player.NickName = content
			s_ReturnContent = fmt.Sprintf(textTemplate["5"], content)
			player.Flag = flag_暂无
			if err := player.UpdateNickName(); err != nil {
				logger.Errorln("wechat UpdateFlag Error:", err)
				s_ReturnContent = textTemplate["100"]
			}

			if err := player.UpdateFlag(); err != nil {
				logger.Errorln("wechat UpdateFlag Error:", err)
				s_ReturnContent = textTemplate["100"]
			}
		}
	default:
		//如果有命令补全，则补全命令
		if player.CommentPrefixStr != "" {
			content = player.CommentPrefixStr + content
		}

		switch {
		//我
		case strings.HasPrefix(content, commandPrefix[0]):

			s_ReturnContent = fmt.Sprintf(textTemplate["6"], player.NickName, Map_MapData[player.Location].Name, player.Level, player.Exp, LevelsExp[player.Level], player.Cur_Mobility, player.Mobility, player.Reputation, player.Coin, player.Cur_HP, player.Max_HP, player.Cur_Burden, player.Max_Burden, player.Cur_Resistance, player.Attack, player.Defense, player.Stamina, player.Agility, player.Wisdom, player.NoDistribution)
			//logger.Debugln(s_ReturnContent)
		case strings.HasPrefix(content, commandPrefix[1]):

		//前往
		case strings.HasPrefix(content, commandPrefix[2]):
			str_AimMap := strings.TrimPrefix(content, commandPrefix[2])

			b_Match := false

			//匹配玩家所在地
			for k, v := range Map_MapData {
				if str_AimMap == v.Name {
					//匹配成功，消除命令前缀
					if player.CommentPrefixStr != "" {
						player.CommentPrefixStr = ""
					}

					b_Match = true

					if v.Level <= player.Level {

						player.Location = k

						s_ReturnContent = fmt.Sprintf(textTemplate["800009"], v.Name, Map_MapData[player.Location].MapDescript)

						if err := player.UpdateLocation(); err != nil {
							logger.Errorln("service wechat UpdateLocation Error:", err)
							s_ReturnContent = textTemplate["100"]
						}
					} else {
						s_ReturnContent = fmt.Sprintf(textTemplate["100028"], v.Level)
					}

					break
				}
				//fmt.Printf("%s -> %s\n", k, v)
			}

			//如果没有匹配到地点，则输出当前玩家可以前往的地点，并且进入命令补全模式
			if !b_Match {

				//如果有命令补全，则消掉
				if player.CommentPrefixStr != "" {
					s_ReturnContent = textTemplate["7"]
					player.CommentPrefixStr = ""
				} else {

					s_ReturnContent += textTemplate["10"]
					for _, v := range Map_MapData {
						if v.Level <= player.Level {
							s_ReturnContent += v.Name + "\n"
						}
					}
					player.CommentPrefixStr = commandPrefix[2]
				}
			}

			//logger.Debugln(s_ReturnContent)

		//修炼
		case strings.HasPrefix(content, commandPrefix[3]):
			//查看玩家当前地点是否适合修炼
			b_CanMapPractice := CanMapPractice(player.Location, model.Func_修炼)
			if !b_CanMapPractice {
				s_ReturnContent = textTemplate["800008"]
				break
			}

			//查看是否有足够的行动力来执行改动作
			temp_Str, b := PlayerCheckMobility(player, -5)
			if b {
				PlayerCheckStatus(player)
				s_ReturnContent = PlayerPratctice(player) + "\n\n" + temp_Str
			} else {
				s_ReturnContent = temp_Str
			}

		case strings.HasPrefix(content, commandPrefix[6]):
			s_ReturnContent = textTemplate["11"]

		//传说
		case strings.HasPrefix(content, commandPrefix[7]):
			s_ReturnContent = fmt.Sprintf(textTemplate["900000"], len(map_PlayerData))

		//道具
		case strings.HasPrefix(content, commandPrefix[9]):
			s_ReturnContent = ShowPlayerProps(player)

		//使用
		case strings.HasPrefix(content, commandPrefix[10]):
			str_PropName := strings.TrimPrefix(content, commandPrefix[10])
			s_ReturnContent = PlayerUseProp(player, str_PropName)

		//提升
		case content == commandPrefix[11]:
			player.CommentPrefixStr = commandPrefix[11]
			if player.NoDistribution > 0 {
				s_ReturnContent = fmt.Sprintf(textTemplate["100003"], player.NoDistribution)
			} else {
				s_ReturnContent = textTemplate["100002"]
			}

		//提升攻击
		case strings.HasPrefix(content, commandPrefix[12]):
			player.CommentPrefixStr = ""
			str := strings.TrimPrefix(content, commandPrefix[12])
			point1, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				logger.Errorln(err)
				s_ReturnContent = textTemplate["900001"]
				break
			}
			point := int(point1)

			if player.NoDistribution-point < 0 {
				s_ReturnContent = fmt.Sprintf(textTemplate["100004"], player.NoDistribution, point, "攻击")
				break
			}

			player.Attack += point
			player.NoDistribution -= point
			player.UpdateNoDistribution()
			player.UpdateAttributes()

			s_ReturnContent = fmt.Sprintf(textTemplate["100001"], "攻击", point, "攻击", player.Attack-point, player.Attack)
		//提升体力
		case strings.HasPrefix(content, commandPrefix[13]):
			player.CommentPrefixStr = ""
			str := strings.TrimPrefix(content, commandPrefix[13])
			point1, err := strconv.ParseInt(str, 10, 32)
			point := int(point1)

			if err != nil {
				logger.Errorln(err)
				s_ReturnContent = textTemplate["900001"]
				break
			}

			if player.NoDistribution-point < 0 {
				s_ReturnContent = fmt.Sprintf(textTemplate["100004"], player.NoDistribution, point, "体力")
				break
			}

			player.Stamina += point
			player.NoDistribution -= point

			player.Max_HP += point * 10
			player.Cur_Resistance += point
			player.Max_Burden += point * 5

			player.UpdateNoDistribution()
			player.UpdateAttributes()

			s_ReturnContent = fmt.Sprintf(textTemplate["100001"], "体力", point, "体力", player.Stamina-point, player.Stamina)
		//提升防御
		case strings.HasPrefix(content, commandPrefix[14]):
			player.CommentPrefixStr = ""
			str := strings.TrimPrefix(content, commandPrefix[14])
			point1, err := strconv.ParseInt(str, 10, 32)
			point := int(point1)
			if err != nil {
				logger.Errorln(err)
				s_ReturnContent = textTemplate["900001"]
				break
			}

			if player.NoDistribution-point < 0 {
				s_ReturnContent = fmt.Sprintf(textTemplate["100004"], player.NoDistribution, point, "防御")
				break
			}

			player.Defense += point
			player.NoDistribution -= point

			player.Cur_Resistance += point * 3

			player.UpdateNoDistribution()
			player.UpdateAttributes()

			s_ReturnContent = fmt.Sprintf(textTemplate["100001"], "防御", point, "防御", player.Defense-point, player.Defense)
		//提升敏捷
		case strings.HasPrefix(content, commandPrefix[15]):
			player.CommentPrefixStr = ""
			str := strings.TrimPrefix(content, commandPrefix[15])
			point1, err := strconv.ParseInt(str, 10, 32)
			point := int(point1)
			if err != nil {
				logger.Errorln(err)
				s_ReturnContent = textTemplate["900001"]
				break
			}
			if player.NoDistribution-point < 0 {
				s_ReturnContent = fmt.Sprintf(textTemplate["100004"], player.NoDistribution, point, "敏捷")
				break
			}
			player.Agility += point
			player.NoDistribution -= point
			s_ReturnContent = fmt.Sprintf(textTemplate["100001"], "敏捷", point, "敏捷", player.Agility-point, player.Agility)

		//提升智慧
		case strings.HasPrefix(content, commandPrefix[16]):
			player.CommentPrefixStr = ""
			str := strings.TrimPrefix(content, commandPrefix[16])
			point1, err := strconv.ParseInt(str, 10, 32)
			point := int(point1)
			if err != nil {
				logger.Errorln(err)
				s_ReturnContent = textTemplate["900001"]
				break
			}
			if player.NoDistribution-point < 0 {
				s_ReturnContent = fmt.Sprintf(textTemplate["100004"], player.NoDistribution, point, "智慧")
				break
			}
			player.Wisdom += point
			player.NoDistribution -= point
			player.Cur_Resistance += point * 3

			player.UpdateNoDistribution()
			player.UpdateAttributes()

			s_ReturnContent = fmt.Sprintf(textTemplate["100001"], "智慧", point, "智慧", player.Wisdom-point, player.Wisdom)
		case strings.HasPrefix(content, commandPrefix[17]):
			s_ReturnContent = ShowCurMapPlayer(player.Location, player.OpenId)

		case strings.HasPrefix(content, commandPrefix[18]):
			s_ReturnContent = PlayerResting(player)

		case strings.HasPrefix(content, commandPrefix[19]):
			//如果是有前缀的，代表是直接买东西的。
			if player.CommentPrefixStr == commandPrefix[19] {
				player.CommentPrefixStr = ""
				str_AimProp := strings.TrimPrefix(content, commandPrefix[19])
				s_ReturnContent = PlayerBuyProps(player, str_AimProp)
			} else {
				//查看该地图是否可以进行买东西操作
				if CanMapPractice(player.Location, model.Func_买) {
					if len(Map_MapData[player.Location].SellItems) > 0 {
						player.CommentPrefixStr = commandPrefix[19]
						s_ReturnContent = textTemplate["100015"]
						for _, propIndex := range Map_MapData[player.Location].SellItems {
							s_ReturnContent += "\n" + fmt.Sprintf(textTemplate["100024"], Map_PropsData[propIndex].Name, Map_PropsData[propIndex].Descript, Map_PropsData[propIndex].OfficialWorth)
						}
					} else {
						s_ReturnContent = textTemplate["100017"]
					}
				} else {
					s_ReturnContent = textTemplate["100013"]
				}
			}
		//卖东西
		case strings.HasPrefix(content, commandPrefix[20]):
			//如果是有前缀的，代表是直接卖东西的。
			if player.CommentPrefixStr == commandPrefix[20] {
				player.CommentPrefixStr = ""
				str_AimProp := strings.TrimPrefix(content, commandPrefix[20])
				s_ReturnContent = PlayerSellProps(player, str_AimProp)

			} else {
				//查看该地图是否可以进行卖东西操作，如果可以则输出玩家可以出售的物品
				if CanMapPractice(player.Location, model.Func_卖) {
					if len(player.Map_PlayerProp) > 0 {
						player.CommentPrefixStr = commandPrefix[20]
						s_ReturnContent = textTemplate["100019"]
						for _, prop := range player.Map_PlayerProp {
							s_ReturnContent += "\n" + fmt.Sprintf(textTemplate["100016"], Map_PropsData[prop.PropId].Name, player.Map_PlayerProp[prop.PropId].PropNum, Map_PropsData[prop.PropId].Descript, Map_PropsData[prop.PropId].Worth)
						}
					} else {
						s_ReturnContent = textTemplate["100018"]
					}
				} else {
					s_ReturnContent = textTemplate["100014"]
				}
			}
		//挑战玩家
		case strings.HasPrefix(content, commandPrefix[21]):

			str_AimPlayerNickName := strings.TrimPrefix(content, commandPrefix[21])

			//判断输入的是否自己的昵称
			if str_AimPlayerNickName == player.NickName {
				s_ReturnContent = textTemplate["800015"]
			} else {
				//获取同地图的某玩家
				targetOpenId, ok := GetPlayerOpenIdByNameAndMapId(str_AimPlayerNickName, player.Location)

				if ok {
					b_Win, _ := Player_VS_Player(player, map_PlayerData[targetOpenId])

					if b_Win {
						s_ReturnContent = fmt.Sprintf(textTemplate["800013"], map_PlayerData[targetOpenId].NickName, map_PlayerData[targetOpenId].NickName)
					} else {
						s_ReturnContent = fmt.Sprintf(textTemplate["800014"], map_PlayerData[targetOpenId].NickName, map_PlayerData[targetOpenId].NickName, Map_MapData[player.Location].Name)
					}
				} else {
					s_ReturnContent = fmt.Sprintf(textTemplate["800012"], Map_MapData[player.Location].Name, str_AimPlayerNickName)
				}
			}
		//我确认洗点
		case strings.HasPrefix(content, commandPrefix[22]):
			//所有属性重置为5点
			s_ReturnContent = PlayerRedistributeAttribute(player)
		default:
			s_ReturnContent = textTemplate["900001"]
		}

	}

	return s_ReturnContent
}

func PlayerPratctice(player *model.WechatPlayer) (s string) {

	//根据玩家所在地图，获取玩家能够匹配到的怪物
	mosterIndex := GetMosterByMap(player.Location)

	//if mosterIndex == -1 {
	//	s = textTemplate["800008"]
	//	return s
	//}

	b_Win, HPLoss := Player_VS_Moster(player, mosterIndex)

	if b_Win {
		s = fmt.Sprintf(textTemplate["800000"], Map_MapData[player.Location].Name, Map_MonsterData[mosterIndex].Name)

		s += "\n\n" + fmt.Sprintf(textTemplate["800005"], -HPLoss, player.Cur_HP, player.Max_HP)

		s += "\n" + fmt.Sprintf(textTemplate["800006"], Map_MonsterData[mosterIndex].Exp)

		//判断是否有物品获得
		propId := GetMosterProp(mosterIndex)

		if propId != -1 {
			//获取物品
			b_Get := PlayerGetProp(player, propId, 1)

			if b_Get {
				s += "\n" + fmt.Sprintf(textTemplate["800007"], Map_PropsData[propId].Name, Map_PropsData[propId].Descript)
			}
		}

		player.Exp += Map_MonsterData[mosterIndex].Exp

		//假设升级就减少经验
		if player.Level < len(LevelsExp) && player.Exp > LevelsExp[player.Level] {
			player.Exp -= LevelsExp[player.Level]

			player.Level++
			player.UpdateLevel()

			player.NoDistribution += 3
			player.UpdateNoDistribution()

			s = s + "\n" + textTemplate["800001"]
		}
		player.UpdateExp()
	} else {
		s = fmt.Sprintf(textTemplate["800004"], Map_MapData[player.Location].Name, Map_MonsterData[mosterIndex].Name, Map_MapData[1].Name)
		s += "\n\n" + fmt.Sprintf(textTemplate["800005"], HPLoss, player.Cur_HP, player.Max_HP)

		//死亡回城
		PlayerGoBackHomeTown(player)
	}

	return s
}

//扣除或添加相应行动力
func PlayerCheckMobility(player *model.WechatPlayer, value int) (s string, b bool) {
	if player.Cur_Mobility+value < 0 {
		b = false
		s = textTemplate["800002"]

	} else {
		player.Cur_Mobility += value
		s = fmt.Sprintf(textTemplate["800003"], value, player.Cur_Mobility, player.Mobility)
		b = true
	}

	return s, b
}

//玩家获取道具
func PlayerGetProp(player *model.WechatPlayer, prop_index int, num int) (b bool) {
	//logger.Debugln(prop_index)
	_, ok := player.Map_PlayerProp[prop_index]

	if ok {
		//logger.Debugln(player.Map_PlayerProp[prop_index])

		player.Map_PlayerProp[prop_index].PropNum += 1
		if err := player.Map_PlayerProp[prop_index].UpdatePlayerPropNum(); err != nil {
			delete(player.Map_PlayerProp, prop_index)
			logger.Errorln("player service PlayerGetProp if error:", err)
			return false
		}
	} else {
		//logger.Debugln(player.Map_PlayerProp[prop_index])

		player.Map_PlayerProp[prop_index] = model.NewWechatPlayerProp()
		player.Map_PlayerProp[prop_index].PropId = prop_index
		player.Map_PlayerProp[prop_index].PlayerId = player.Id
		player.Map_PlayerProp[prop_index].PropNum = 1

		id, err := player.Map_PlayerProp[prop_index].Insert()
		if err != nil {
			delete(player.Map_PlayerProp, prop_index)
			logger.Errorln("player service PlayerGetProp else error:", err)
			return false
		}
		player.Map_PlayerProp[prop_index].Id = id
	}

	return true
}

//显示玩家所拥有的道具
func ShowPlayerProps(player *model.WechatPlayer) (s string) {
	if len(player.Map_PlayerProp) > 0 {
		s += textTemplate["100005"] + "\n\n"
		for _, v := range player.Map_PlayerProp {
			if prop, ok := Map_PropsData[v.PropId]; ok {
				s += fmt.Sprintf(textTemplate["100007"], prop.Name, v.PropNum, prop.Descript) + "\n"
			} else {
				logger.Errorln("player service ShowPlayerProps error: No this prop")
			}
		}

	} else {
		s += textTemplate["100006"]
	}
	return s
}

//根据字符串判断玩家是否拥有此道具,拥有的话则返回一个正确的道具ID
func CheckPlayerHasProp(player *model.WechatPlayer, propName string) (propId int, b bool) {
	b = false
	propId = -1
	for _, v := range player.Map_PlayerProp {
		propInfo, ok := Map_PropsData[v.PropId]
		if ok {
			if propInfo.Name == propName {
				propId = propInfo.Id
				b = true
				break
			}
		} else {
			logger.Debugln("CheckPlayerHasProp : Not find prop in Map_PropData")
		}
	}

	return propId, b
}

//玩家使用道具
func PlayerUseProp(player *model.WechatPlayer, propName string) (s string) {

	//查看玩家是否存在此道具,返回ID

	propId, ok := CheckPlayerHasProp(player, propName)

	if !ok {
		s = fmt.Sprintf(textTemplate["100009"], propName)
		return
	}

	targetProp, ok := Map_PropsData[propId]

	//根据道具的类型来使用
	if ok {
		switch targetProp.PropType {
		case model.PropType_恢复生命值:
			addpoint := targetProp.PropValue
			player.Cur_HP += targetProp.PropValue

			if player.Cur_HP > player.Max_HP {
				addpoint = targetProp.PropValue - (player.Cur_HP - player.Max_HP)
				player.Cur_HP = player.Max_HP
			}
			s = fmt.Sprintf(textTemplate["100008"], targetProp.Name, addpoint, player.Cur_HP, player.Max_HP)
		case model.PropType_恢复行动力:
			addpoint := targetProp.PropValue
			player.Cur_Mobility += targetProp.PropValue

			if player.Cur_Mobility > player.Mobility {
				addpoint = targetProp.PropValue - (player.Cur_Mobility - player.Mobility)
				player.Cur_Mobility = player.Mobility
			}
			s = fmt.Sprintf(textTemplate["100026"], targetProp.Name, addpoint, player.Cur_Mobility, player.Mobility)
		case model.PropType_角色昵称更改:
			player.Flag = flag_用户传入角色名申请更名操作
			s = textTemplate["100029"]
		default:
			s = textTemplate["100010"]
		}

		//减少道具
		DecreasePlayerProp(player, targetProp.Id)
	} else {
		s = textTemplate["900002"]
	}
	return s
}

//玩家出售某个道具
func PlayerSellProps(player *model.WechatPlayer, propName string) (s string) {
	propIndex, b := CheckPlayerHasProp(player, propName)
	if b {
		player.Coin += Map_PropsData[propIndex].Worth
		player.UpdateCoin()
		DecreasePlayerProp(player, propIndex)
		s = fmt.Sprintf(textTemplate["100021"], propName, Map_PropsData[propIndex].Worth, player.Coin)
	} else {
		s = fmt.Sprintf(textTemplate["100020"], propName)
	}
	return s
}

//玩家购买某个道具
func PlayerBuyProps(player *model.WechatPlayer, propName string) (s string) {
	propIndex, b := CheckMapSellProps(player.Location, propName)
	if b {
		//假设够钱就买，不够钱则提示。
		if player.Coin-Map_PropsData[propIndex].OfficialWorth >= 0 {
			player.Coin -= Map_PropsData[propIndex].OfficialWorth
			player.UpdateCoin()

			if PlayerGetProp(player, propIndex, 1) {
				s = fmt.Sprintf(textTemplate["100023"], Map_PropsData[propIndex].OfficialWorth, propName)
			} else {
				s = textTemplate["900002"]
			}
		} else {
			s = fmt.Sprintf(textTemplate["100022"], player.Coin, propName, Map_PropsData[propIndex].OfficialWorth)
		}
	} else {
		s = fmt.Sprintf(textTemplate["100025"], propName)
	}
	return s
}

//减少玩家某个道具
func DecreasePlayerProp(player *model.WechatPlayer, prop int) {
	player.Map_PlayerProp[prop].PropNum -= 1
	if player.Map_PlayerProp[prop].PropNum == 0 {
		model.NewWechatPlayerProp().Where("id=" + strconv.Itoa(player.Map_PlayerProp[prop].Id)).Delete()
		delete(player.Map_PlayerProp, prop)
	} else {
		player.Map_PlayerProp[prop].UpdatePlayerPropNum()
	}
}

//返回除了自己以外的所有玩家。“附近”
func ShowCurMapPlayer(mapId int, playerOpenId string) (s string) {

	s = ""
	for _, v := range map_PlayerData {
		if v.Location == mapId && v.OpenId != playerOpenId {
			s += "\n" + v.NickName
		}
	}

	if s != "" {
		s = textTemplate["800010"] + s
	} else {
		s = textTemplate["800011"]
	}

	return s
}

//玩家进入休息状态
func PlayerResting(player *model.WechatPlayer) (s string) {

	if player.Status == 1 {
		s = textTemplate["100012"]
		return s
	} else {
		player.Status = 1
	}
	player.Timer = time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-player.Timer.C:
				if player.Cur_HP < player.Max_HP {
					player.Cur_HP++
				}

				if player.Cur_Mobility < player.Mobility {
					player.Cur_Mobility++
				}
				logger.Debugln("Player Name:" + player.NickName + " PlayerHp:" + strconv.Itoa(player.Cur_HP) + " PlayerMobility:" + strconv.Itoa(player.Cur_Mobility))

				if player.Cur_HP >= player.Max_HP && player.Cur_Mobility >= player.Mobility {
					player.Timer.Stop()
				}
			}
		}
	}()

	s = textTemplate["100011"]
	return s
}

//查看玩家是否有持续性的状态，若有则打断
func PlayerCheckStatus(player *model.WechatPlayer) {
	if player.Status == 1 {
		player.Status = 0
		player.Timer.Stop()
	}
}

//玩家与玩家对战，一方失败位置，扣血量为玩家。死亡的玩家回城。
func Player_VS_Player(me *model.WechatPlayer, target *model.WechatPlayer) (b_Win bool, HPLoss int) {
	//初始化
	HPLoss = 0
	b_Win = false

	//自己的攻击倍率
	rate := float32(me.Agility) / float32(target.Agility)

	//自己的DPS
	meHurt := (float32(me.Attack) - float32(target.Defense)) * rate
	if meHurt <= 0.0 {
		meHurt = 1.0
	}

	//对方的DPS
	targetHurt := float32(target.Attack) - float32(me.Defense)
	if targetHurt <= 0.0 {
		targetHurt = 1.0
	}

	meDPSTime := float32(target.Cur_HP) / meHurt
	targetDPSTime := float32(me.Cur_HP) / targetHurt

	//logger.Debugln(meHurt)
	//logger.Debugln(targetHurt)
	//logger.Debugln(meDPSTime)
	//logger.Debugln(targetDPSTime)

	if meDPSTime <= targetDPSTime {

		HPLoss = int(meDPSTime * targetHurt)
		me.Cur_HP -= HPLoss

		target.Cur_HP = 0
		PlayerGoBackHomeTown(target)

		b_Win = true
	} else {

		HPLoss = int(targetDPSTime * meHurt)
		target.Cur_HP -= HPLoss

		me.Cur_HP = 0
		PlayerGoBackHomeTown(me)

		b_Win = false
	}

	return b_Win, HPLoss
}

//玩家回城
func PlayerGoBackHomeTown(player *model.WechatPlayer) {
	player.Location = 1
	player.UpdateLocation()
}

//根据昵称，获取在某地图上的玩家ID。
func GetPlayerOpenIdByNameAndMapId(playerName string, curMapId int) (playerOpenId string, b bool) {
	playerOpenId = ""
	b = false
	for _, v := range map_PlayerData {
		if v.Location == curMapId && v.NickName == playerName {
			playerOpenId = v.OpenId
			b = true
		}
	}
	return
}

//玩家重新分配点数，洗点
func PlayerRedistributeAttribute(player *model.WechatPlayer) (s string) {
	player.Attack = 5
	player.Defense = 5
	player.Agility = 5
	player.Stamina = 5
	player.Wisdom = 5

	player.UpdateAttributes()

	player.NoDistribution = player.Level * 3
	player.UpdateNoDistribution()

	player.Max_HP = player.Stamina * 10

	if player.Cur_HP > player.Stamina*10 {
		player.Cur_HP = player.Stamina * 10
	}

	player.Cur_Resistance = player.Stamina + player.Defense*3 + player.Wisdom*3
	player.Max_Burden = player.Stamina * 5
	player.Cur_Burden = 0

	s = textTemplate["100027"]
	return s
}
