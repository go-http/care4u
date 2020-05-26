package care4u

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	UrlDailyReport    = "/ess/report/dailyReport"
	UrlGetDailyReport = "/ess/report/getDailyReport"
)

// 日报中的地址信息
type Address struct {
	Prprovince string `json:"prprovince"` //省
	Prcity     string `json:"prcity"`     //市
	Prarea     string `json:"prarea"`     //区
	Prstreet   string `json:"prstreet"`   //街
}

// 上传的日报结构
type ReportUpload struct {
	StudId       string `json:"studId"`       //学生ID
	Type         string `json:"type"`         //学生状态: 5无、1居家隔离、2确诊、3疑似、4留观
	ClassId      string `json:"classId"`      //班级ID
	Contact      string `json:"contact"`      //是否接触，0未接触、1常居、2共餐、3日常接触
	Phycondition string `json:"phycondition"` //身体状态: 0正常
	Temperature  string `json:"temperature"`  //体温:

	Address
}

// 下载的日报结构
type ReportDownload struct {
	Address

	HasInChengdu bool

	StudentReport struct {
		Id           string
		StdId        string
		Remark       string
		Createtime   string
		Reportdate   string
		Type         int
		ClassId      string
		Contact      string
		Phycondition string
		Temperature  string
	}
}

// SetDailyReport 用于获取当日的已填报的健康报告
func (cli *Client) GetDailyReport() (*ReportDownload, error) {
	param := url.Values{
		"classId": {cli.ClassId},
		"id":      {cli.StudentId},
	}

	var info ReportDownload

	err := cli.post(UrlGetDailyReport, param, &info)

	return &info, err
}

// SetDailyReport 用于上报当日的健康报告
// 其中ReportUpload中的StuId、ClassId会自动填写
func (cli *Client) SetDailyReport(report ReportUpload) error {
	//自动填写学生和班级的ID
	if cli.CheckVerifyInfo == nil || cli.ClassInfo == nil {
		return fmt.Errorf("未成功获取班级信息或未检查填报状态")
	}

	report.StudId = cli.CheckVerifyInfo.Id
	report.ClassId = cli.ClassInfo.Id

	data, err := json.Marshal(report)
	if err != nil {
		return err
	}

	param := url.Values{
		"data": {string(data)},
	}

	err = cli.post(UrlDailyReport, param, nil)

	return err
}
