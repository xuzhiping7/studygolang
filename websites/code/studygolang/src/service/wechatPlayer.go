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
)

//储存所有对话模板
var textTemplate map[int]string

func init() {
	//初始化所有对话模板
	textTemplate = make(map[int]string)

	textTemplate[0] = "注册成功!\n\n感谢注册微信奇幻网游《传说》，希望您能够享受和喜欢这个世界。\n\n（请输出‘传说’两个字，开启您的游戏旅程！）"
	textTemplate[1] = "您已经注册过，请输出‘我’查看您的最新状态。"
	textTemplate[2] = "注册"
	textTemplate[3] = "欢迎来到微信奇幻网游《传说》，请输出'注册'，确认注册游戏。"
	textTemplate[4] = "创建角色中,请输入您的角色名。(例如‘一叶之秋’，8个汉字内。)"
	textTemplate[5] = "角色【%s】成功创建！请输入'传说'两字开始游戏。"
	textTemplate[100] = "网络错误，请重新输入。"
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
		if player.Flag == 0 {
			s_ReturnContent = textTemplate[4]
			player.Flag = 1

			if err := player.UpdateFlag(); err != nil {
				logger.Errorln("wechat UpdateFlag Error:", err)
				s_ReturnContent = textTemplate[100]
			}
		} else if player.Flag == 1 {
			player.Flag = 2
			runes := []rune(content)
			if len(runes) > 8 {
				s_ReturnContent = textTemplate[4]
			} else {
				player.NickName = content
				s_ReturnContent = fmt.Sprintf(textTemplate[5], content)

				if err := player.UpdateNickName(); err != nil {
					logger.Errorln("wechat UpdateFlag Error:", err)
					s_ReturnContent = textTemplate[100]
				}

				if err := player.UpdateFlag(); err != nil {
					logger.Errorln("wechat UpdateFlag Error:", err)
					s_ReturnContent = textTemplate[100]
				}
			}

		} else {
			s_ReturnContent = textTemplate[1]
		}
	}
	return s_ReturnContent

}

//func UpdateNickName(openid string, name string) bool {

//}
