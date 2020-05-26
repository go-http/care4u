package care4u

import (
	"fmt"
	"net/url"
)

type Client struct {
	ClassId   string
	StudentId string
	Token     string

	ClassInfo       *ClassInfo
	CheckVerifyInfo *VerifyInfo
}

// New 创建一个新客户端，用于上报健康信息或获取健康信息
// 创建过程中会检查班级号、学生姓名是否有效
func New(classCode, studentName string) (*Client, error) {
	classInfo, err := GetClassInfo(classCode)
	if err != nil {
		return nil, fmt.Errorf("获取班级信息失败: %s", err)
	}

	// 获取当日填报状态
	checkVerifyInfo, err := CheckVerify(classInfo.Id, studentName)
	if err != nil {
		return nil, fmt.Errorf("获取当前填报状态失败: %s", err)
	}

	client := &Client{
		ClassId:         classInfo.Id,
		StudentId:       checkVerifyInfo.Id,
		Token:           checkVerifyInfo.Token,
		ClassInfo:       classInfo,
		CheckVerifyInfo: checkVerifyInfo,
	}

	return client, nil
}

// Post 自动发起请求，并带上Token（如果不为空）
func (cli *Client) Post(urlPath string, param url.Values, respData interface{}) error {
	return Post(urlPath, cli.Token, param, respData)
}

// HasDailyReport 返回当日是否已经上报健康信息
func (cli *Client) HasDailyReport() bool {
	if cli.CheckVerifyInfo == nil {
		return false
	}

	return cli.CheckVerifyInfo.HasDailyReport
}
