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
		if len(runes) > 8 {
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

			s_ReturnContent = fmt.Sprintf(textTemplate["6"], player.NickName, Map_MapData[player.Location].Name, player.Level, player.Exp, 100, player.Cur_Mobility, player.Mobility, player.Reputation, player.Coin, player.Cur_HP, player.Max_HP, player.Cur_Burden, player.Max_Burden, player.Cur_Resistance, player.Attack, player.Defense, player.Stamina, player.Agility, player.Wisdom, player.NoDistribution)
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
					player.Location = k

					s_ReturnContent = fmt.Sprintf(textTemplate["800009"], v.Name, Map_MapData[player.Location].MapDescript)

					if err := player.UpdateLocation(); err != nil {
						logger.Errorln("service wechat UpdateLocation Error:", err)
						s_ReturnContent = textTemplate["100"]
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
						s_ReturnContent += v.Name + "\n"
					}
					player.CommentPrefixStr = commandPrefix[2]
				}
			}

			//logger.Debugln(s_ReturnContent)

		//修炼
		case strings.HasPrefix(content, commandPrefix[3]):
			//查看玩家当前地点是否适合修炼
			b_CanMapPractice := CanMapPractice(player.Location)
			if !b_CanMapPractice {
				s_ReturnContent = textTemplate["800008"]
				break
			}

			//查看是否有足够的行动力来执行改动作
			temp_Str, b := PlayerCheckMobility(player, -5)
			if b {
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
			player.UpdateNoDistribution()
			player.UpdateAttributes()

			s_ReturnContent = fmt.Sprintf(textTemplate["100001"], "智慧", point, "智慧", player.Wisdom-point, player.Wisdom)
		case strings.HasPrefix(content, commandPrefix[17]):
			s_ReturnContent = ShowCurMapPlayer(player.Location, player.OpenId)
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
		if player.Exp > 100 {
			player.Exp -= 100

			player.Level++
			player.UpdateLevel()

			player.NoDistribution += 3
			player.UpdateNoDistribution()

			s = s + "\n" + textTemplate["800001"]
		}
		player.UpdateExp()
	} else {
		s = fmt.Sprintf(textTemplate["800004"], Map_MapData[player.Location].Name, Map_MonsterData[mosterIndex].Name)
		s += "\n\n" + fmt.Sprintf(textTemplate["800005"], HPLoss, player.Cur_HP, player.Max_HP)
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

//减少玩家某个道具
func DecreasePlayerProp(player *model.WechatPlayer, prop int) {
	player.Map_PlayerProp[prop].PropNum -= 1
	if player.Map_PlayerProp[prop].PropNum == 0 {
		delete(player.Map_PlayerProp, prop)
		model.NewWechatPlayerProp().Where("id=" + strconv.Itoa(player.Map_PlayerProp[prop].Id)).Delete()
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
