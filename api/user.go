package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type schoolList struct {
	Code   int           `json:"code"`
	Msg    string        `json:"msg"`
	Result schoolResList `json:"result"`
	Status bool          `json:"status"`
}

type schoolResList struct {
	List []schoolObj `json:"list"`
}

type schoolObj struct {
	Badge string `json:"badge"`
	Host  string `json:"host"`
	Id    string `json:"id"`
	Ident string `json:"ident"`
	Name  string `json:"name"`
}

type LoginRes struct {
	Code   int       `json:"code"`
	Msg    string    `json:"msg"`
	Result LoginInfo `json:"result"`
	Status bool      `json:"status"`
}

type LoginInfo struct {
	Data LoginObj `json:"data"`
}

type LoginObj struct {
	Id          int    `json:"id"`
	Token       string `json:"token"`
	Name        string `json:"name"`
	ClassId     int    `json:"classId"`
	CollegeId   int    `json:"collegeId"`
	Point       int    `json:"point"`
	Rank        int    `json:"rank"`
	ClassName   string `json:"className"`
	CollegeName string `json:"collegeName"`
}

type SignObj struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
	Status bool        `json:"status"`
}

var Host string
var Token string
var SchoolId string


// 获取学校
func GetSchool() {
	url := "http://www.minghuaetc.com/api/login/school.json"
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := schoolList{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		log.Printf("获取学校信息列表成功\n\n")
		fmt.Printf("%s %s\n", "序号", "学校名称")
		var schoolIndex = 0
		for _, item := range m.Result.List {
			fmt.Printf("%03d %s\n", schoolIndex, item.Name)
			schoolIndex++
		}
		var selectIndex int
		fmt.Printf("\n请输入学校的序号: ")
		//reader := bufio.NewReader(os.Stdin)
		fmt.Scanf("%d\n", &selectIndex)
		for {
			if selectIndex < len(m.Result.List) {
				break
			}
			fmt.Printf("请输入学校的序号: ")
			fmt.Scanf("%d\n", &selectIndex)
		}
		fmt.Printf("选择了 %s\n", m.Result.List[(selectIndex)].Name)
		Host = m.Result.List[selectIndex].Host
		SchoolId = m.Result.List[selectIndex].Id
	} else {
		log.Fatalf("学校列表获取失败 %s\n", m.Msg)
	}
}

// 登录
func Login() {
	var username string
	var password string
	//reader := bufio.NewReader(os.Stdin)
	fmt.Printf("请输入学号: ")
	fmt.Scanf("%s\n", &username)
	//username, _ = reader.ReadString('\n')
	fmt.Printf("请输入密码: ")
	fmt.Scanf("%s\n", &password)
	//password, _ = reader.ReadString('\n')
	url := Host + "/api/login.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("username", username)
	mpWriter.WriteField("password", password)
	mpWriter.WriteField("schoolId", SchoolId)
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
	m := LoginRes{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	if m.Status {
		Token = m.Result.Data.Token
		log.Printf("%s 登录成功 ! 来自 %s-%s 的 %s %s\n", username, m.Result.Data.CollegeName, m.Result.Data.ClassName, m.Result.Data.Name, signIn())
	} else {
		log.Fatalf("%s 登录失败! %s", username, m.Msg)
	}
	GetCourseList()
}


// 签到
func signIn() string {
	url := Host + "/api/user/sign_in.json"
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
	m := SignObj{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		log.Fatal(err)
	}
	return m.Msg
}
