package api

import (
	"bytes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

func GetAnswer() {
	url := Host + "/api/answer.json"
	var bufReader bytes.Buffer
	mpWriter := multipart.NewWriter(&bufReader)
	mpWriter.WriteField("schoolId", "6")
	mpWriter.WriteField("questionId", "1")
	mpWriter.WriteField("token", Token)
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, &bufReader)
	request.Header.Add("Content-Type", "multipart/form-data; boundary=" + mpWriter.Boundary())
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
