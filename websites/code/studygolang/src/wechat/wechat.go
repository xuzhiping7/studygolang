package wechat

import (
	//"filter"
	"fmt"
	//"github.com/studygolang/mux"
	"logger"
	"net/http"
)

func WechatTest() {

	logger.Debugln("wechat go begin !")
}

func WechatTest2(rw http.ResponseWriter, req *http.Request) {
	logger.Debugln(req.RequestURI)

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
	} else {

		fmt.Fprint(rw, `xuzhiping test`)
	}
}
