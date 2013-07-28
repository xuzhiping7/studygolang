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
	"math/rand"
	"model"
	"strings"
)

//储存玩家信息，24小时后删除
var map_PlayerData map[string]*model.WechatPlayer

//储存所有对话模板
var textTemplate map[string]string

//储存所有命令前缀
var commandPrefix map[int]string

//储存地图ID对应名
var map_MapName map[int]string

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

	map_MapName = make(map[int]string)
	map_MapName[0] = "林风角酒馆"
	map_MapName[1] = "林风角"
	map_MapName[2] = "林风南海岸"
}

// 创建一个wechat玩家
func CreateWechatPlayer(openid string) bool {
	player := model.NewWechatPlayer()

	player.OpenId = openid
	player.NickName = "EmptyNow"
	player.UserName = "EmptyNow"
	player.Exp = 0
	player.Mobility = 0

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
	return player
}

// 判断该OpenID是否已经被注册了
func OpenidExists(openid string) bool {
	player := model.NewWechatPlayer()
	if err := player.Where("openid=" + openid).Find("id"); err != nil {
		logger.Errorln("service EmailExists error:", err)
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
		switch {
		case strings.HasPrefix(content, commandPrefix[0]):
			s_ReturnContent = fmt.Sprintf(textTemplate["6"], player.NickName, map_MapName[player.Location], player.Level, player.Exp, 100, "吟游诗人", "三寸黄金", "无", player.Mobility, player.Reputation, "453/656", "56/100", "25", player.Attack, player.Defense, player.Stamina, player.Agility, player.NoDistribution)
			logger.Debugln(s_ReturnContent)
		case strings.HasPrefix(content, commandPrefix[1]):

		case strings.HasPrefix(content, commandPrefix[2]):
			str_AimMap := strings.TrimPrefix(content, commandPrefix[2])

			b_Match := false

			//匹配玩家所在地
			for k, v := range map_MapName {
				if str_AimMap == v {
					b_Match = true
					player.Location = k

					s_ReturnContent = fmt.Sprintf(textTemplate["8"], v)

					if err := player.UpdateLocation(); err != nil {
						logger.Errorln("service wechat UpdateLocation Error:", err)
						s_ReturnContent = textTemplate["100"]
					}

					break
				}
				//fmt.Printf("%s -> %s\n", k, v)
			}

			//如果没有匹配到地点，则输出当前玩家可以前往的地点
			if !b_Match {
				s_ReturnContent += textTemplate["10"]
				for _, v := range map_MapName {
					s_ReturnContent += fmt.Sprintf(textTemplate["9"], v)
				}
			}
			logger.Debugln(s_ReturnContent)

		//修炼
		case strings.HasPrefix(content, commandPrefix[3]):
			s_ReturnContent = PlayerPratctice(player)
		case strings.HasPrefix(content, commandPrefix[6]):
			s_ReturnContent = textTemplate["11"]
		//传说
		case strings.HasPrefix(content, commandPrefix[7]):
			s_ReturnContent = fmt.Sprintf(textTemplate["900000"], len(map_PlayerData))
		default:
			s_ReturnContent = textTemplate["1"]
		}

	}

	return s_ReturnContent
}

func PlayerPratctice(player *model.WechatPlayer) (s string) {
	if rand.Intn(100) > 50 {
		s = fmt.Sprintf(textTemplate["800000"], map_MapName[player.Location], "风铃怪", 10, "风信子", 2)
		player.Exp += 10
	} else {
		s = fmt.Sprintf(textTemplate["800000"], map_MapName[player.Location], "泥巴怪", 5, "粘土", 2)
		player.Exp += 5
	}

	//假设升级就减少经验
	if player.Exp > 100 {
		player.Exp -= 100
		player.Level++
		player.UpdateLevel()

		s = s + "\n" + textTemplate["800001"]

	}
	player.UpdateExp()

	return s
}
