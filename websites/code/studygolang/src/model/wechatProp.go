package model

type WechatProp struct {
	Id       int
	Name     string
	Descript string
}

func NewWechatProp(id int, name string, descript string) *WechatProp {
	return &WechatProp{
		Id:       id,
		Name:     name,
		Descript: descript,
	}
}
