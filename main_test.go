package main

import (
	"bytes"
	"log"
	"testing"

	"github.com/newde36524/Go-Tool/httptool"
)

func TestGet(t *testing.T) {
	httpClient := &httptool.HttpClient{}
	result, err := httpClient.Get("https://localhost:44316/api/values", nil)
	if err == nil {
		log.Println(result)
	} else {
		log.Println(err)
	}
}
func TestGet2(t *testing.T) {
	httpClient := &httptool.HttpClient{}
	result, err := httpClient.Get("https://localhost:44316/api/values/122222", nil)
	if err == nil {
		log.Println(result)
	} else {
		log.Println(err)
	}
}
func TestPost(t *testing.T) {
	httpClient := httptool.NewHttpClient()
	httpClient.Header["Content-Type"] = "application/json"
	result, err := httpClient.Post("https://localhost:44316/api/values/PostData", `{
		"name":"tom"
	}`)
	if err == nil {
		log.Println(result)
	} else {
		log.Println(err)
	}
}
func TestPost2(t *testing.T) {
	httpClient := httptool.NewHttpClient()
	httpClient.Header["Content-Type"] = "text/plain"
	result, err := httpClient.Post("https://localhost:44316/api/values/PostData2", `{
		"value":"tom"
	}`)
	if err == nil {
		log.Println(result)
	} else {
		log.Println(err)
	}
}
func TestForm(t *testing.T) {
	httpClient := httptool.NewHttpClient()

	forms := make([]httptool.FormItem, 0)
	forms = append(forms, httptool.FormItem{
		Name:  "msg",
		Value: "测试文本",
	})
	forms = append(forms, httptool.FormItem{
		Name:        "file",
		FileData:    bytes.NewBufferString("sss"),
		FileName:    "file.txt",
		ContentType: "text/plain",
	})
	result, err := httpClient.Form("https://localhost:44316/api/values/Form", forms)
	if err == nil {
		log.Println(result)
	} else {
		log.Println(err)
	}
}
func TestForm2(t *testing.T) {
	httpClient := httptool.NewHttpClient()

	forms := make([]httptool.FormItem, 0)
	forms = append(forms, httptool.FormItem{
		Name:  "msg",
		Value: "测试文本",
	})

	result, err := httpClient.Form("https://localhost:44316/api/values/Form2", forms)
	if err == nil {
		log.Println(result)
	} else {
		log.Println(err)
	}
}
