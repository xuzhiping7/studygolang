// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"fmt"
	"logger"
	"model"
	"strings"
)

//储存所有对话模板
var textTemplate map[int]string

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
	//初始化所有对话模板
	textTemplate = make(map[int]string)

	textTemplate[0] = "注册成功!\n\n感谢注册微信奇幻网游《传说》，目前游戏处于删档测试阶段，有什么问题和建议请联系邮箱xuzhiping7@qq.com，希望您能够享受和喜欢这个世界。\n\n（请输出‘传说’两个字，开启您的游戏旅程！）"
	textTemplate[1] = "您已经注册过，请输出‘我’查看您的最新状态。"
	textTemplate[2] = "注册"
	textTemplate[3] = "欢迎来到微信奇幻网游《传说》，请输出'注册'，确认注册游戏。"
	textTemplate[4] = "创建角色中,请输入您的角色名。(例如‘一叶之秋’，8个汉字内。)"
	textTemplate[5] = "角色【%s】成功创建！请输入'传说'两字开始游戏。"
	textTemplate[6] = "%s\n当前地点:%s\n\n等级:%d\n职业:%s\n称号:%s\n状态：%s\n\n行动力:%d\n声望:%d\n\n生命:%s\n负重:%s\n抗性:%s\n\n攻击:%d\n防御:%d\n体力:%d\n敏捷:%d\n剩余分配点数:%d\n\n(您可以输入'当前'查看你所能做的事情)"
	textTemplate[7] = "抱歉，你不能到达这个地方。"
	textTemplate[8] = "你到达了%s"
	textTemplate[100] = "网络错误，请重新输入。"

	commandPrefix = make(map[int]string)
	commandPrefix[0] = "我"
	commandPrefix[1] = "当前"
	commandPrefix[2] = "前往"
	commandPrefix[3] = "修炼"
	commandPrefix[4] = "状态"
	commandPrefix[5] = "搜寻"

	map_MapName = make(map[int]string)
	map_MapName[0] = "林风阁酒馆"
	map_MapName[1] = "林风村"
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

	if !OpenidExists(openid) {

		//假设当前用户输入的不是注册
		if content == textTemplate[2] {
			b_Reg := CreateWechatPlayer(openid)

			if b_Reg {
				s_ReturnContent = textTemplate[0]
			}
		} else {
			s_ReturnContent = textTemplate[3]
		}

	} else {
		player := GetWechatPlayer(openid)
		logger.Debugln(player)

		switch player.Flag {
		case flag_注册完成:
			s_ReturnContent = textTemplate[4]
			player.Flag = flag_用户传入角色名申请更名操作
			if err := player.UpdateFlag(); err != nil {
				logger.Errorln("wechat UpdateFlag Error:", err)
				s_ReturnContent = textTemplate[100]
			}
		case flag_用户传入角色名申请更名操作:
			player.Flag = flag_用户传入角色名申请更名操作
			runes := []rune(content)
			if len(runes) > 8 {
				s_ReturnContent = textTemplate[4]
			} else {
				player.NickName = content
				s_ReturnContent = fmt.Sprintf(textTemplate[5], content)
				player.Flag = flag_暂无
				if err := player.UpdateNickName(); err != nil {
					logger.Errorln("wechat UpdateFlag Error:", err)
					s_ReturnContent = textTemplate[100]
				}

				if err := player.UpdateFlag(); err != nil {
					logger.Errorln("wechat UpdateFlag Error:", err)
					s_ReturnContent = textTemplate[100]
				}
			}
		default:
			switch {
			case strings.HasPrefix(content, commandPrefix[0]):
				s_ReturnContent = fmt.Sprintf(textTemplate[6], player.NickName, map_MapName[player.Location], player.Level, "吟游诗人", "三寸黄金", "无", player.Mobility, player.Reputation, "453/656", "56/100", "25", player.Attack, player.Defense, player.Stamina, player.Agility, player.NoDistribution)
				logger.Debugln(s_ReturnContent)
			case strings.HasPrefix(content, commandPrefix[1]):

			case strings.HasPrefix(content, commandPrefix[2]):
				str_AimMap := strings.TrimPrefix(content, commandPrefix[2])

				for k, v := range map_MapName {
					if str_AimMap == v {

						player.Location = k

						s_ReturnContent = fmt.Sprintf(textTemplate[8], v)

						if err := player.UpdateLocation(); err != nil {
							logger.Errorln("service wechat UpdateLocation Error:", err)
							s_ReturnContent = textTemplate[100]
						}

						break
					} else {
						s_ReturnContent = textTemplate[7]
					}
					//fmt.Printf("%s -> %s\n", k, v)
				}

			case strings.HasPrefix(content, commandPrefix[3]):

			default:
				s_ReturnContent = textTemplate[1]
			}

		}

	}
	return s_ReturnContent

}

//func UpdateNickName(openid string, name string) bool {

//}
