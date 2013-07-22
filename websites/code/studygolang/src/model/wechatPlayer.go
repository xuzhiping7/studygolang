package model

import (
	"logger"
	"util"
)

// 帖子信息
type WechatPlayer struct {
	Id       int    `json:"id"`
	OpenId   string `json:"openid"`
	NickName string `json:"nickname"`
	UserName string `json:"username"`
	Exp      int    `json:"exp"`
	Mobility int    `json:"mobility"`

	// 数据库访问对象
	*Dao
}

func NewWechatPlayer() *WechatPlayer {
	return &WechatPlayer{
		Dao: &Dao{tablename: "wechat_base"},
	}
}

func (this *WechatPlayer) Insert() (int, error) {
	this.prepareInsertData()
	logger.Debugln(this)
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

func (this *WechatPlayer) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *WechatPlayer) FindAll(selectCol ...string) ([]*WechatPlayer, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	playerList := make([]*WechatPlayer, 0, 10)
	logger.Debugln("selectCol", selectCol)
	colNum := len(selectCol)
	for rows.Next() {
		player := NewWechatPlayer()
		err = this.Scan(rows, colNum, player.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("WechatPlayer FindAll Scan Error:", err)
			continue
		}
		playerList = append(playerList, player)
	}
	return playerList, nil
}

// 为了支持连写
func (this *WechatPlayer) Set(clause string) *WechatPlayer {
	this.Dao.Set(clause)
	return this
}

// 为了支持连写
func (this *WechatPlayer) Where(condition string) *WechatPlayer {
	this.Dao.Where(condition)
	return this
}

// 为了支持连写
func (this *WechatPlayer) Limit(limit string) *WechatPlayer {
	this.Dao.Limit(limit)
	return this
}

// 为了支持连写
func (this *WechatPlayer) Order(order string) *WechatPlayer {
	this.Dao.Order(order)
	return this
}

func (this *WechatPlayer) prepareInsertData() {

	this.columns = []string{"openid", "username", "nickname", "exp", "mobility"}
	this.colValues = []interface{}{this.OpenId, this.UserName, this.NickName, this.Exp, this.Mobility}

}

func (this *WechatPlayer) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":       &this.Id,
		"openid":   &this.OpenId,
		"username": &this.UserName,
		"nickname": &this.NickName,
		"exp":      &this.Exp,
		"mobility": &this.Mobility,
	}
}
