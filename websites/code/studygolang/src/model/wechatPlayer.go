package model

import (
	"logger"
	"strconv"
	"time"
	"util"
)

// 帖子信息
type WechatPlayer struct {
	/*
		数据表中储存的数据
	*/
	Id             int    `json:"id"`
	OpenId         string `json:"openid"`
	NickName       string `json:"nickname"`
	UserName       string `json:"username"`
	Sex            int    `json:"sex"`
	Level          int    `json:"level"`
	Exp            int    `json:"exp"`
	Coin           int    `json:"coin"`
	Mobility       int    `json:"mobility"`
	Reputation     int    `json:"reputation"`
	Attack         int    `json:"attack"`
	Defense        int    `json:"defense"`
	Stamina        int    `json:"stamina"`
	Agility        int    `json:"agility"`
	Wisdom         int    `json:"wisdom"`
	NoDistribution int    `json:"no_distribution"`
	Location       int    `json:"location"`
	Flag           int    `json:"flag"`

	/*
		数据库访问对象
	*/
	*Dao

	/*
		扩展的动态数据
	*/
	//HP值
	Max_HP int
	Cur_HP int
	//当前行动力
	Cur_Mobility int
	//当前抗性
	Cur_Resistance int
	//负重
	Max_Burden int
	Cur_Burden int
	//玩家道具
	Map_PlayerProp map[int]*WechatPlayerProp

	/*
		辅助数据
	*/
	//当前命令前缀
	CommentPrefixStr string
	//当前前置事件提醒，例如被挑战
	ReplyPreifixStr string
	//玩家历史事件记录
	RecordEvent *WechatPlayerRecord
	//玩家当前状态:0-无，1-休息
	Status int
	//玩家单独计时器
	Timer *time.Ticker
}

func NewWechatPlayer() *WechatPlayer {
	return &WechatPlayer{
		CommentPrefixStr: "",
		Dao:              &Dao{tablename: "wechat_base"},
	}
}

func (this *WechatPlayer) Insert() (int, error) {
	this.prepareInsertData()
	//logger.Debugln(this)
	result, err := this.Dao.Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

//更新用户当前事件节点
func (this *WechatPlayer) UpdateFlag() error {
	err := this.Set("flag=" + strconv.Itoa(this.Flag)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新用户昵称
func (this *WechatPlayer) UpdateNickName() error {
	err := this.Set("nickname=" + this.NickName).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新地名
func (this *WechatPlayer) UpdateLocation() error {
	err := this.Set("location=" + strconv.Itoa(this.Location)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新经验值
func (this *WechatPlayer) UpdateExp() error {
	err := this.Set("exp=" + strconv.Itoa(this.Exp)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新等级
func (this *WechatPlayer) UpdateLevel() error {
	err := this.Set("level=" + strconv.Itoa(this.Level)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新金币
func (this *WechatPlayer) UpdateCoin() error {
	err := this.Set("coin=" + strconv.Itoa(this.Coin)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新行动力
func (this *WechatPlayer) UpdateMobility() error {
	err := this.Set("mobility=" + strconv.Itoa(this.Mobility)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新分配点(NoDistribution)
func (this *WechatPlayer) UpdateNoDistribution() error {
	err := this.Set("no_distribution=" + strconv.Itoa(this.NoDistribution)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
}

//更新属性点
func (this *WechatPlayer) UpdateAttributes() error {
	//logger.Debugln("attack=" + strconv.Itoa(this.Attack) + ",defense=" + strconv.Itoa(this.Defense) + ",stamina=" + strconv.Itoa(this.Stamina) + ",agility=" + strconv.Itoa(this.Agility) + ",wisdom=" + strconv.Itoa(this.Wisdom))
	err := this.Set("attack=" + strconv.Itoa(this.Attack) + ",defense=" + strconv.Itoa(this.Defense) + ",stamina=" + strconv.Itoa(this.Stamina) + ",agility=" + strconv.Itoa(this.Agility) + ",wisdom=" + strconv.Itoa(this.Wisdom)).Where("openid=" + this.OpenId).Update()
	this.SetEmpty()
	return err
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
	this.columns = []string{"openid", "username", "nickname", "exp", "mobility", "attack", "defense", "stamina", "agility", "wisdom"}
	this.colValues = []interface{}{this.OpenId, this.UserName, this.NickName, this.Exp, this.Mobility, this.Attack, this.Defense, this.Stamina, this.Agility, this.Wisdom}
}

func (this *WechatPlayer) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":              &this.Id,
		"openid":          &this.OpenId,
		"username":        &this.UserName,
		"nickname":        &this.NickName,
		"exp":             &this.Exp,
		"coin":            &this.Coin,
		"sex":             &this.Sex,
		"mobility":        &this.Mobility,
		"location":        &this.Location,
		"flag":            &this.Flag,
		"level":           &this.Level,
		"attack":          &this.Attack,
		"defense":         &this.Defense,
		"stamina":         &this.Stamina,
		"agility":         &this.Agility,
		"wisdom":          &this.Wisdom,
		"no_distribution": &this.NoDistribution,
	}
}
