package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
)

type replyAddRes struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
	Status bool        `json:"status"`
}

//添加课程的评论
func AddReply(chapter chapterNodeObj) string {
	url := Host + "/api/node/add_reply.json"
	var bufReader bytes.Buffer
	an := regexp.MustCompile("[1-9]\\d*")
	chapterName := an.ReplaceAllString(chapter.Name, "")
	content := "学习 " + chapterName + " 记录打卡 !"
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("nodeId", strconv.Itoa(chapter.Id))
	mpWriter.WriteField("content", content)
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary="+mpWriter.Boundary())
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := replyAddRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m.Msg
}
