package model

type WechatMap struct {
	Id          int
	Name        string
	MapDescript string
	Mosters     []int
	MostersRate []int
	Functions   []int
	SellItems   []int
}

const (
	Func_修炼 = 1
	Func_买  = 2
	Func_卖  = 3
)

func NewWechatMap(id int, name string, mapDescript string, mosters []int, mostersRate []int, functions []int, sellItems []int) *WechatMap {
	return &WechatMap{
		Id:          id,
		Name:        name,
		MapDescript: mapDescript,
		Mosters:     mosters,
		MostersRate: mostersRate,
		Functions:   functions,
		SellItems:   sellItems,
	}
}
