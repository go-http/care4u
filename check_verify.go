package care4u

import (
	"net/url"
)

const (
	urlCheckVerify = "/ess/report/checkVerify"
)

type VerifyInfo struct {
	Contact        string
	HasBasicInfo   bool
	HasDailyReport bool
	HasInChengdu   string
	HasStudent     bool
	Id             string
	Name           string
	Phycondition   string
	Temperature    string
	Token          string
	Type           string
}

// CheckVerify 用于检查当日健康报告的填写状况
func CheckVerify(classId, stuName string) (*VerifyInfo, error) {
	param := url.Values{
		"classId": {classId},
		"code":    {stuName},
	}

	var info VerifyInfo

	err := post(urlCheckVerify, "", param, &info)

	return &info, err
}
