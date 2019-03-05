package httptool

import "net/url"

// 获取urlencode后的http请求url参数
func getUrlEncodeParams(params map[string]string) string {
	p := url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return p.Encode()
}
