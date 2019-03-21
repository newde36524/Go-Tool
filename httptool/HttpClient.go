package httptool

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type HttpClient struct {
	Header map[string]string
}

func NewHttpClient() *HttpClient {
	// text/plain
	return &HttpClient{
		Header: make(map[string]string, 0),
	}
}

//发起Get请求
func (this *HttpClient) Get(url string, params map[string]string) (string, error) {
	targetUrl := url
	if params != nil {
		targetUrl = targetUrl + "?" + getUrlEncodeParams(params)
	}
	content, err := this.httpHandle(Get(), targetUrl, strings.NewReader(""))
	return content, err
}

//发起Post请求
func (this *HttpClient) Post(url, body string) (string, error) {
	content, err := this.httpHandle(Post(), url, strings.NewReader(body))
	return content, err
}

type FormItem struct {
	Name        string
	ContentType string
	Value       string
	FileName    string
	FileData    io.Reader
}

// FormData提交
func (this *HttpClient) Form(url string, formItemList []FormItem) (string, error) {
	boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
	spliteStr := boundary
	this.Header["Content-Type"] = "multipart/form-data; boundary=" + boundary
	buffers := make([]io.Reader, 0)
	slice := make([]string, 0)
	slice = append(slice, "")
	slice = append(slice, spliteStr)

	for _, v := range formItemList {
		ContentDisposition := "Content-Disposition: form-data; name=\"" + v.Name + "\""
		if v.FileName != "" && v.FileData != nil {
			ContentDisposition = ContentDisposition + ";filename=" + v.FileName
			buffers = append(buffers, v.FileData)
		}
		slice = append(slice, ContentDisposition)
		if v.ContentType != "" {
			slice = append(slice, "Content-Type:"+v.ContentType)
		}
		slice = append(slice, "")
		slice = append(slice, v.Value)
		slice = append(slice, spliteStr)
		slice = append(slice, "")
	}
	buffer := bytes.NewBufferString(strings.Join(slice, "\r\n"))
	// buffers = append(buffers[:0], buffer)
	endBuffers := make([]io.Reader, 0)
	endBuffers = append(endBuffers, buffer)
	endBuffers = append(endBuffers, buffers...)
	request_reader := io.MultiReader(endBuffers...)

	content, err := this.httpHandle(Post(), url, request_reader)
	return content, err
}

//http请求处理程序
func (this *HttpClient) httpHandle(method *HttpMethod, url string, body io.Reader) (string, error) {
	req, _ := http.NewRequest(method.method, url, body)
	for k, v := range this.Header {
		req.Header.Add(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

/*
POST https://api.example.com/user/upload
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="text"

title
------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="image"; filename="1.png"
Content-Type: image/png

< ./1.png
------WebKitFormBoundary7MA4YWxkTrZu0gW--


*/
func postFile(filename string, target_url string) (*http.Response, error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("userfile", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", filename)
		return nil, err
	}
	req, err := http.NewRequest("POST", target_url, request_reader)
	if err != nil {
		return nil, err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

	return http.DefaultClient.Do(req)
}
