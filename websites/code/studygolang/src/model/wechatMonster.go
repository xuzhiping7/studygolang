package model

//import (
//	"logger"
//)

type WechatMonster struct {
	//野怪ID
	Id   int
	Name string

	//野怪基本属性
	HP         int
	Attack     int
	Defense    int
	Agility    int
	Resistance int

	//打败野怪可以获得的经验
	Exp int
	//打败野怪可以获得的道具
	Props []int
	//获得道具的概率
	PropsGetRate []int
}

func NewWechatMonster(id int, name string, hp int, attack int, defense int, agility int, resistance int, exp int, props []int, propsRate []int) *WechatMonster {
	return &WechatMonster{
		Id:           id,
		Name:         name,
		HP:           hp,
		Attack:       attack,
		Defense:      defense,
		Agility:      agility,
		Resistance:   resistance,
		Exp:          exp,
		Props:        props,
		PropsGetRate: propsRate,
	}
}
