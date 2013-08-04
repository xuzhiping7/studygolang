package model

type WechatMap struct {
	Id          int
	Name        string
	MapDescript string
	Mosters     []int
	MostersRate []int
}

func NewWechatMap(id int, name string, mapDescript string, mosters []int, mostersRate []int) *WechatMap {
	return &WechatMap{
		Id:          id,
		Name:        name,
		MapDescript: mapDescript,
		Mosters:     mosters,
		MostersRate: mostersRate,
	}
}
