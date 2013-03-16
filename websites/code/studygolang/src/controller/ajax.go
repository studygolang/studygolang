package controller

import (
	"encoding/json"
	"fmt"
	"github.com/studygolang/mux"
	"logger"
	"model"
	"net/http"
	"service"
	"strconv"
)

// 侧边栏的内容通过异步请求获取

// 某节点下其他帖子
func OtherTopicsHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	topics := service.FindTopicsByNid(vars["nid"], vars["tid"])
	data, err := json.Marshal(topics)
	if err != nil {
		logger.Errorln("[OtherTopicsHandler] json.marshal error:", err)
		fmt.Fprint(rw, `{"errno": 1, "error":"解析json出错"}`)
		return
	}
	fmt.Fprint(rw, `{"errno": 0, "data":`+string(data)+`}`)
}

func StatHandler(rw http.ResponseWriter, req *http.Request) {
	topicTotal := service.TopicsTotal()
	replyTotal := service.CommentsTotal(model.TYPE_TOPIC)
	userTotal := service.CountUsers()
	fmt.Fprint(rw, `{"errno": 0, "data":{"topic":`+strconv.Itoa(topicTotal)+`,"reply":`+strconv.Itoa(replyTotal)+`,"user":`+strconv.Itoa(userTotal)+`}}`)
}
