package model

import (
	"logger"
	"strconv"
	"util"
)

//道具表
type WechatProp struct {
	Id            int
	Name          string
	Descript      string
	Worth         int
	OfficialWorth int
	PropType      int
	PropValue     int
}

const (
	PropType_没有任何作用 = 0
	PropType_恢复生命值  = 1
	PropType_恢复行动力  = 2
	PropType_角色昵称更改 = 3
)

/*
PropType:
1 --- 恢复生命值
2 ---

*/

func NewWechatProp(id int, name string, descript string, worth int, officialWorth int, propType int, propValue int) *WechatProp {
	return &WechatProp{
		Id:            id,
		Name:          name,
		Descript:      descript,
		Worth:         worth,
		OfficialWorth: officialWorth,
		PropType:      propType,
		PropValue:     propValue,
	}
}

//玩家道具表
type WechatPlayerProp struct {
	Id       int `json:"id"`
	PlayerId int `json:"player_id"`
	PropId   int `json:"prop_id"`
	PropNum  int `json:"prop_num"`

	/*
		数据库访问对象
	*/
	*Dao
}

func NewWechatPlayerProp() *WechatPlayerProp {
	return &WechatPlayerProp{
		Dao: &Dao{tablename: "wechat_player_prop"},
	}
}

//根据ID来更新数目
func (this *WechatPlayerProp) UpdatePlayerPropNum() error {
	err := this.Set("prop_num=" + strconv.Itoa(this.PropNum)).Where("id=" + strconv.Itoa(this.Id)).Update()
	this.SetEmpty()
	return err
}

func (this *WechatPlayerProp) Insert() (int, error) {
	this.prepareInsertData()
	//logger.Debugln(this)
	result, err := this.Dao.Insert()
	this.SetEmpty()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return int(id), err
}

// 为了支持连写
func (this *WechatPlayerProp) Set(clause string) *WechatPlayerProp {
	this.Dao.Set(clause)
	return this
}

// 为了支持连写
func (this *WechatPlayerProp) Where(condition string) *WechatPlayerProp {
	this.Dao.Where(condition)
	return this
}

func (this *WechatPlayerProp) Find(selectCol ...string) error {
	return this.Dao.Find(this.colFieldMap(), selectCol...)
}

func (this *WechatPlayerProp) FindAll(selectCol ...string) ([]*WechatPlayerProp, error) {
	if len(selectCol) == 0 {
		selectCol = util.MapKeys(this.colFieldMap())
	}
	rows, err := this.Dao.FindAll(selectCol...)
	if err != nil {
		return nil, err
	}
	// TODO:
	nodeList := make([]*WechatPlayerProp, 0, 10)
	//logger.Debugln("selectCol", selectCol)

	colNum := len(selectCol)
	for rows.Next() {
		node := NewWechatPlayerProp()
		err = this.Scan(rows, colNum, node.colFieldMap(), selectCol...)
		if err != nil {
			logger.Errorln("TopicNode FindAll Scan Error:", err)
			continue
		}
		nodeList = append(nodeList, node)
	}
	return nodeList, err
}

func (this *WechatPlayerProp) prepareInsertData() {
	this.columns = []string{"player_id", "prop_id", "prop_num"}
	this.colValues = []interface{}{this.PlayerId, this.PropId, this.PropNum}
}

func (this *WechatPlayerProp) colFieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        &this.Id,
		"player_id": &this.PlayerId,
		"prop_id":   &this.PropId,
		"prop_num":  &this.PropNum,
	}
}
