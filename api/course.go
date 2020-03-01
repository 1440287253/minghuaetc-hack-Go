package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
)

type courseRes struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Result courseResList `json:"result"`
	Status bool          `json:"status"`
}

type courseResList struct {
	Announcement interface{} `json:"announcement"` //公告的信息
	Banner       interface{} `json:"banner"`       //banner 横幅广告的信息
	Enlist       interface{} `json:"enlist"`       //全校各学院开放的课程
	Finish       interface{} `json:"finish"`       //已经结束的课程
	List         []courseObj `json:"list"`         //当前正在进行的课程
}

type courseObj struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	ClassTeacher string  `json:"classTeacher"`
	Progress     float32 `json:"progress"`
}

type chapterRes struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result chapterList `json:"result"`
	Status bool        `json:"status"`
}

type chapterList struct {
	List []chapterListObj `json:"list"`
}

type chapterListObj struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	NodeList []chapterNodeObj `json:"nodeList"`
}

type chapterNodeObj struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	VideoDuration string `json:"videoDuration"`
	VideoState    int    `json:"videoState"`
	TabVideo      bool   `json:"tabVideo"`
}

type studyRes struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result studyResObj `json:"result"`
	Status bool        `json:"status"`
}

type studyResObj struct {
	Data studyObj `json:data`
}

type studyObj struct {
	StudyId int `json:"studyId"`
}

var CourseResObj = courseRes{}
var reply bool

//获取课程列表
func GetCourseList() {
	url := Host + "/api/course/list.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
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
	err = json.Unmarshal(body, &CourseResObj)
	if err != nil {
		log.Fatal(err)
	}
	if CourseResObj.Status {
		//log.Printf("当前正在学习课程获取成功\n\n")
		fmt.Printf("\n%s %s %s\n", "课程序号", "完成度", "课程名称")
		var courseIndex = 0
		for _, item := range CourseResObj.Result.List {
			fmt.Printf("%03d %.2f%% %s\n", item.Id, item.Progress*100, item.Name)
			courseIndex++
		}
		var selectIndex int
		fmt.Printf("\n请输入课程的序号: ")
		fmt.Scanf("%d", &selectIndex)
		getCourseChapter(selectIndex)
	} else {
		log.Fatalf("课程获取失败 %s\n", CourseResObj.Msg)
	}
}

//获取课程的详情，但如果提交课程，并不需要此请求
func getCourseDetail(courseId int) {
	url := Host + "/api/course/detail.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("courseId", strconv.Itoa(courseId))
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
	log.Println(string(body))
}

//获取课程章节
func getCourseChapter(courseId int) {
	url := Host + "/api/course/chapter.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("courseId", strconv.Itoa(courseId))
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
	m := chapterRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		var replyYes string
		fmt.Printf("是否开启学习评论 格式为 学习 [课程名称] 记录打卡 ? [y/N] ")
		fmt.Scan(&replyYes)
		switch replyYes {
		case "y":
			fmt.Printf("已经开启学习后自动评论 !\n")
			reply = true
		case "n":
			fmt.Printf("关闭学习后自动评论 !\n")
			reply = false
		default:
			fmt.Printf("不开启学习后自动评论 !\n")
			reply = false
		}
		for _, course := range m.Result.List {
			chapterName := course.Name
			for _, chapter := range course.NodeList {
				if chapter.TabVideo {
					if chapter.VideoState == 2 {
						log.Printf("%s[%s] 已经完成学习，自动略过", chapterName, chapter.Name)
					} else {
						doChapter(chapterName, chapter)
					}
				} else {
					log.Printf("%s[%s] 不是视频教学 !", chapterName, chapter.Name)
				}
			}
		}
	} else {
		log.Printf("课程章节获取出错 %s 请确保输入了正确的课程序号", m.Msg)
	}
	GetCourseList()
}

//判断是否需要学习
func doChapter(chapterName string, chapter chapterNodeObj) {
	var answer string
YesOrNo:
	fmt.Printf("%s[%s] 完成学习? [y/n] ", chapterName, chapter.Name)
	fmt.Scan(&answer)
	switch answer {
	case "y":
		studyChapter(chapter)
	case "n":
		//log.Printf("%s[%s] 跳过了学习", chapterName, chapter.Name)
	default:
		goto YesOrNo
	}
}

//完成章节学习
func studyChapter(chapter chapterNodeObj) {
	url := Host + "/api/node/study.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("nodeId", strconv.Itoa(chapter.Id))
	mpWriter.WriteField("studyTime", chapter.VideoDuration)
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
	m := studyRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		if reply {
			replyMsg := AddReply(chapter)
			log.Printf("%s %s %s\n", chapter.Name, m.Msg, replyMsg)
		} else {
			log.Printf("%s %s\n", chapter.Name, m.Msg)
		}
	} else {
		log.Printf("%s 学习失败 ! %s\n", chapter.Name, m.Msg)
	}
}
