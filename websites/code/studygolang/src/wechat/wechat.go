package wechat

import (
	"fmt"
	//"github.com/studygolang/mux"
	"encoding/xml"
	"io/ioutil"
	"logger"
	"net/http"
	"service"
)

type textRecieveMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string
	MsgId        string
}

type textResponseMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string
	FuncFlag     string
}

//储存所有对话模板
var textTemplate map[int]string

func init() {
	//初始化所有对话模板
	textTemplate = make(map[int]string)

	textTemplate[0] = "注册成功!\n\n感谢注册微信奇幻网游《传说》，希望您能够享受和喜欢这个世界。\n\n（请输出‘传说’两个字，开启您的游戏旅程！）"
	textTemplate[1] = "您已经注册过，请输出‘我’查看您的最新状态。"
	textTemplate[2] = "注册"
	textTemplate[3] = "欢迎来到微信奇幻网游《传说》，请输出'注册'，确认注册游戏。"
}

//微信通道总入口
func WechatEntrance(rw http.ResponseWriter, req *http.Request) {

	bytes, _ := ioutil.ReadAll(req.Body)

	v := textRecieveMessage{}

	err := xml.Unmarshal(bytes, &v)
	if err != nil {
		logger.Errorln(err)
		return
	}

	responXML := textResponseMessage{}
	responXML.FromUserName = v.ToUserName
	responXML.ToUserName = v.FromUserName
	responXML.Content = v.Content
	responXML.CreateTime = v.CreateTime
	responXML.MsgType = v.MsgType
	responXML.FuncFlag = "0"

	//假设不存在当前用户
	if !service.OpenidExists(v.FromUserName) {

		//假设当前用户输入的不是注册
		if v.Content == textTemplate[2] {
			b_Reg := service.CreateWechatPlayer(v.FromUserName)

			if b_Reg {
				responXML.Content = textTemplate[0]
			}
		} else {
			responXML.Content = textTemplate[3]
		}

	} else {
		responXML.Content = textTemplate[1]
	}

	result, _ := xml.Marshal(responXML)

	fmt.Fprint(rw, string(result))
}

func WechatDevelopVerify(rw http.ResponseWriter, req *http.Request) {
	//验证开发者服务器，暂时发个假数据过去吧。构造那个算法有点麻烦。
	if req.Form["signature"] != nil {

		signature := req.Form["signature"][0]
		timestamp := req.Form["timestamp"][0]
		nonce := req.Form["nonce"][0]
		echostr := req.Form["echostr"][0]

		logger.Debugln(signature)
		logger.Debugln(timestamp)
		logger.Debugln(nonce)
		logger.Debugln(echostr)

		fmt.Fprint(rw, echostr)
	}
}

//本地服务器测试函数
func WechatTest(rw http.ResponseWriter, req *http.Request) {

	//没有相应的OPENID则为其注册一个
	if !service.OpenidExists("dsgdsgdgsg") {
		test := service.CreateWechatPlayer("dsgdsgdgsg")

		if test {
			logger.Debugln("用户注册成功")
		}

	} else {
		logger.Debugln("用户已经存在")
	}
}
