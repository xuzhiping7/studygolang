package wechat

import (
	//"filter"
	"fmt"
	"github.com/studygolang/mux"
	"logger"
	"net/http"
)

func WechatTest() {

	logger.Debugln("wechat go begin !")
}

func WechatTest2(rw http.ResponseWriter, req *http.Request) {
	logger.Debugln(req.RequestURI)
	logger.Debugln(req.Method)

	logger.Debugln(req.Header.Get("token"))

	vars := mux.Vars(req)

	logger.Debugln(vars["token"])
	//logger.Debugln("wechat 1")
	logger.Debugln(vars)

	fmt.Fprint(rw, `xuzhiping test`)
}
