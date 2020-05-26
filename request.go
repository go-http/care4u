package care4u

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	ApiVer  = "1.0"
	UrlBase = "https://tcm.care4u.cn:8443"
)

var generalHeaders = map[string]string{
	"Accept-Language": "zh-cn",
	"Accept-Encoding": "br, gzip, deflate",
	"Accept":          "application/json, text/javascript, */*; q=0.01",
	"Referer":         "https://cdntcm.care4u.cn/wechat/index.html?urlCode=3f031153133b",
	"Origin":          "https://cdntcm.care4u.cn",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1 Safari/605.1.15",
	"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
}

// ResponseCommon 代表API返回的通用数据
// 其中的Data字段，当有数据时是实际的数据结构，没有数据时是空字符串
// 这样的设计在Golang的标准库中无法直接解析，因此视作json.RawMessage以懒加载
type ResponseCommon struct {
	Msg   string
	State int
	Data  json.RawMessage
}

func (r *ResponseCommon) Error() error {
	if r.State != 1 || r.Msg != "success" {
		return fmt.Errorf("[%d]%s", r.State, r.Msg)
	}

	return nil
}

func Post(urlPath, token string, param url.Values, respData interface{}) error {
	// 设置公共的API版本号参数
	if param == nil {
		param = url.Values{}
	}
	param.Set("ver", ApiVer)

	req, err := http.NewRequest(http.MethodPost, UrlBase+urlPath, strings.NewReader(param.Encode()))
	if err != nil {
		return err
	}

	// 插入一些常见的HTTP头
	for k, v := range generalHeaders {
		req.Header.Set(k, v)
	}

	if token != "" {
		req.Header.Add("token", token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respInfo ResponseCommon
	err = json.NewDecoder(resp.Body).Decode(&respInfo)
	if err != nil {
		return err
	}

	if err := respInfo.Error(); err != nil {
		return err
	}

	// 如果上层不关心返回数据的业务部分，则直接返回
	if respData == nil {
		return nil
	}

	// 如果返回数据的业务部分为空字符串或者为空，也直接返回
	if string(respInfo.Data) == `""` || string(respInfo.Data) == "" {
		return nil
	}

	// 其他情况下统一解析返回的JSON数据
	err = json.Unmarshal(respInfo.Data, respData)

	return err
}
