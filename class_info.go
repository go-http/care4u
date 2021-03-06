package care4u

import (
	"net/url"
)

const (
	urlGetClassInfo = "/ess/report/getClassInfo"
)

// ClassInfo 表示班级信息
type ClassInfo struct {
	Id         string
	ClassName  string
	SchoolName string
	StudentNum int
}

// GetClassInfo 用于获取指定班级号的班级信息，班级号来自于URL中的urlCode参数
func GetClassInfo(code string) (*ClassInfo, error) {
	param := url.Values{
		"urlCode": {code},
	}

	var info ClassInfo

	err := post(urlGetClassInfo, "", param, &info)

	return &info, err
}
