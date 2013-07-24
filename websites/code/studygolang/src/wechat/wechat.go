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
	textTemplate[4] = "创建角色中,请输入您的角色名。(例如‘一叶之秋’，8个汉字内。)"
	textTemplate[5] = "角色【%s】成功创建！请输入'传说'两字开始游戏。"
	textTemplate[100] = "网络错误，请重新输入。"
}

//微信通道总入口
func WechatEntrance(rw http.ResponseWriter, req *http.Request) {

	v := textRecieveMessage{}

	bytes, _ := ioutil.ReadAll(req.Body)
	err := xml.Unmarshal(bytes, &v)
	if err != nil {
		logger.Errorln(err)
		return
	}
	//v.FromUserName = "xuzhipingtest"
	//v.Content = "超帅的烧饼2"

	responXML := textResponseMessage{}
	responXML.FromUserName = v.ToUserName
	responXML.ToUserName = v.FromUserName
	responXML.Content = v.Content
	responXML.CreateTime = v.CreateTime
	responXML.MsgType = v.MsgType
	responXML.FuncFlag = "0"

	/*
		假设不存在当前用户,让玩家进入注册流程
		存在当前用户，则读取当前用户的所有信息。
	*/
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
		player := service.GetWechatPlayer(v.FromUserName)
		logger.Debugln(player)
		if player.Flag == 0 {
			responXML.Content = textTemplate[4]
			player.Flag = 1

			if err := player.UpdateFlag(); err != nil {
				logger.Errorln("wechat UpdateFlag Error:", err)
				responXML.Content = textTemplate[100]
			}
		} else if player.Flag == 1 {
			player.Flag = 2
			runes := []rune(v.Content)
			if len(runes) > 8 {
				responXML.Content = textTemplate[4]
			} else {
				player.NickName = v.Content
				responXML.Content = fmt.Sprintf(textTemplate[5], v.Content)

				if err := player.UpdateNickName(); err != nil {
					logger.Errorln("wechat UpdateFlag Error:", err)
					responXML.Content = textTemplate[100]
				}

				if err := player.UpdateFlag(); err != nil {
					logger.Errorln("wechat UpdateFlag Error:", err)
					responXML.Content = textTemplate[100]
				}
			}

		} else {
			responXML.Content = textTemplate[1]
		}
	}

	result, _ := xml.Marshal(responXML)

	fmt.Fprint(rw, string(result))
}

//开发者认证返回值
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
