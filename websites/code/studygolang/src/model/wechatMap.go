package model

type WechatMap struct {
	Id      int
	Name    string
	Mosters []int
}

func NewWechatMap(id int, name string, mosters []int) *WechatMap {
	return &WechatMap{
		Id:      id,
		Name:    name,
		Mosters: mosters,
	}
}
