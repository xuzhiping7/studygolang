// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
	"model"
	"strings"
)

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
	//初始化所有对话模板
	textTemplate = make(map[string]string)

	date, err := ioutil.ReadFile(config.ROOT + "/conf/wechatTextTemplate.json")

	if err != nil {
		logger.Errorln("Read wechatTextTemplate.json fail error:", err)
	}

	if err2 := json.Unmarshal(date, &textTemplate); err2 != nil {
		logger.Errorln("Unmarshal wechatTextTemplate.json fail error:", err2)
	}

	commandPrefix = make(map[int]string)
	commandPrefix[0] = "我"
	commandPrefix[1] = "当前"
	commandPrefix[2] = "前往"
	commandPrefix[3] = "修炼"
	commandPrefix[4] = "状态"
	commandPrefix[5] = "搜寻"
	commandPrefix[6] = "帮助"

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
		if content == textTemplate["2"] {
			b_Reg := CreateWechatPlayer(openid)

			if b_Reg {
				s_ReturnContent = textTemplate["0"]
			}
		} else {
			s_ReturnContent = textTemplate["3"]
		}

	} else {
		player := GetWechatPlayer(openid)
		logger.Debugln(player)

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
				s_ReturnContent = fmt.Sprintf(textTemplate["6"], player.NickName, map_MapName[player.Location], player.Level, "吟游诗人", "三寸黄金", "无", player.Mobility, player.Reputation, "453/656", "56/100", "25", player.Attack, player.Defense, player.Stamina, player.Agility, player.NoDistribution)
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

			case strings.HasPrefix(content, commandPrefix[3]):

			case strings.HasPrefix(content, commandPrefix[6]):
				s_ReturnContent = textTemplate["11"]
			default:
				s_ReturnContent = textTemplate["1"]
			}

		}

	}
	return s_ReturnContent

}

//func UpdateNickName(openid string, name string) bool {

//}
