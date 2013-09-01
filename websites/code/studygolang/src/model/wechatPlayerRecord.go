package model

import (
	//"logger"
	//"strconv"
	"time"
)

type WechatPlayerRecord struct {
	Flag   int
	Record [10]string
}

func NewWechatPlayerRecord() *WechatPlayerRecord {
	//logger.Debugln("NewWechatPlayerRecord")
	temp := &WechatPlayerRecord{}
	temp.InitRecord()
	return temp
}

//初始化所有记录
func (this *WechatPlayerRecord) InitRecord() {
	this.Flag = 0
	for i := 0; i < len(this.Record); i++ {
		this.Record[i] = ""
	}
}

//增加一个玩家事件记录
func (this *WechatPlayerRecord) AddRecord(s string) {
	this.Record[this.Flag] = time.Now().Format("[01.02 15:04]") + s
	this.Flag++
	if this.Flag >= len(this.Record) {
		this.Flag = 0
	}
	//logger.Debugln(this.Record)
}

//遍历输出当前所有事件记录
func (this *WechatPlayerRecord) GetAllRecord() (s string) {

	s = ""
	for i := this.Flag - 1; i >= 0; i-- {
		if this.Record[i] != "" {
			s += this.Record[i] + "\n"
		}
	}

	for i := this.Flag; i < len(this.Record); i++ {
		if this.Record[i] != "" {
			s += this.Record[i] + "\n"
		}
	}

	return s
}
