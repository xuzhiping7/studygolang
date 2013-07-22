// Copyright 2013 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author：polaris	studygolang@gmail.com

package service

import (
	"logger"
	"model"
)

// 创建一个wechat玩家
func CreateWechatPlayer(openid string) bool {
	player := model.NewWechatPlayer()

	player.OpenId = openid
	player.NickName = "葱烧烙饼"
	player.UserName = "xuzhiping"
	player.Exp = 0
	player.Mobility = 0

	if _, err := player.Insert(); err != nil {
		logger.Errorln("player service CreateWechatPlayer error:", err)
		return false
	}
	return true
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
