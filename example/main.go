package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-http/care4u"
)

func main() {
	classCode := os.Getenv("CLASS_CODE")
	studentName := os.Getenv("STUDENT_NAME")
	address := os.Getenv("ADDRESS")
	if classCode == "" || studentName == "" || address == "" {
		fmt.Println("请通过CLASS_CODE、STUDENT_NAME环境变量分别设置班级码和学生姓名")
		return
	}

	addrs := strings.Split(address, " ")
	if len(addrs) != 4 {
		fmt.Println("请按照「省 市 区 街道」的格式，在ADDRESS环境变量中设置地址")
		return
	}

	err := AutoReport(classCode, studentName, addrs)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func AutoReport(classCode, studentName string, addrs []string) error {
	client, err := care4u.New(classCode, studentName)
	if err != nil {
		return fmt.Errorf("获取班级信息失败: %s", err)
	}

	// 如果当日已填报，则输出填报信息
	if client.HasDailyReport() {
		dailyReport, err := client.GetDailyReport()
		if err != nil {
			return fmt.Errorf("获取已填信息失败: %s", err)
		}
		fmt.Printf("已填日报，内容是: %s", dailyReport)

		return nil
	}

	// 如果当日未填报，则填报
	report := care4u.ReportUpload{
		Temperature:  "36.5", //体温正常
		Phycondition: "0",    //身体健康
		Contact:      "0",    //无新冠病人接触史
		Type:         "5",    //学生状态正常
	}

	report.Prprovince = addrs[0]
	report.Prcity = addrs[1]
	report.Prarea = addrs[2]
	report.Prstreet = addrs[3]

	err = client.SetDailyReport(report)
	if err != nil {
		return fmt.Errorf("填报失败", err)
	}

	fmt.Printf("填报成功")

	return nil
}
