package docUtil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"baliance.com/gooxml/document"
)

type DocUtil struct {
}

//ConvertTemplateByFilePath 结构体数据填充到doc文档
//@filepath 模板文件地址
//@data 填充的数据
/*
	ssssssss{字段名}sssss
*/
func (d *DocUtil) ConvertTemplateByFilePath(filepath string, data interface{}) (io.Reader, error) {
	doc, err := document.Open(filepath)
	if err != nil {
		return nil, err
	}
	return d.process(doc, data)
}

//ConvertTemplateByReader 结构体数据填充到doc文档
//@r 模板文件
//@data 填充的数据
/*
	ssssssss{字段名}sssss
*/
func (d *DocUtil) ConvertTemplateByReader(r io.Reader, data interface{}) (io.Reader, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	doc, err := document.Read(bytes.NewReader(bs), int64(len(bs)))
	if err != nil {
		return nil, err
	}
	return d.process(doc, data)
}

//process 处理文档
func (d *DocUtil) process(doc *document.Document, data interface{}) (io.Reader, error) {
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			text := run.Text()
			run.Clear()
			text, err := d.replaceP(text, d.structToMap(data))
			if err != nil {
				return nil, err
			}
			run.AddText(text)
		}
	}
	buf := bytes.NewBuffer(nil)
	if err := doc.Save(buf); err != nil {
		return nil, err
	}
	return buf, nil
}
func (d *DocUtil) structToMap(c interface{}) map[string]string {
	t := reflect.TypeOf(c)
	v := reflect.ValueOf(c)
	var data = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).String()
	}
	return data
}

func (d *DocUtil) findPlaceHolder(str string, fn func(placeHolder string)) (err error) {
	defer func() {
		if err == io.EOF {
			err = nil
		}
	}()
	rd := bufio.NewReader(strings.NewReader(str))
	for {
		_, err = rd.ReadString('{')
		if err != nil {
			return err
		}
		placeHolder, err := rd.ReadString('}')
		if err != nil {
			return err
		}
		fn(strings.TrimRight(placeHolder, "}"))
	}
}

func (d *DocUtil) replaceP(p string, mp map[string]string) (result string, err error) {
	return p, d.findPlaceHolder(p, func(placeHolder string) {
		v := mp[strings.ToLower(placeHolder)]
		p = strings.Replace(p, fmt.Sprintf("{%s}", placeHolder), v, 1)
	})
}
