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

func WechatEntrance(rw http.ResponseWriter, req *http.Request) {

	/*
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
	*/

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

	if !service.OpenidExists(v.FromUserName) {
		b_Reg := service.CreateWechatPlayer(v.FromUserName)

		if b_Reg {
			responXML.Content = "用户注册成功"
		}

	} else {
		responXML.Content = "已经注册"
	}

	result, _ := xml.Marshal(responXML)

	fmt.Fprint(rw, string(result))
}

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
