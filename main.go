package main

import (
	"fmt"
	"goFastForVideo/file"
	"goFastForVideo/query"
	"goFastForVideo/study"
	"math/rand"
	"strings"
	"sync"
	"time"
)

var query_url = "http://plat.xzjxjy.com"

func main() {
	var Wg sync.WaitGroup
	cons := file.ReadSetting()
	for _, con := range cons {
		yearArr := strings.Split(con.Years, ",")
		myCourses := query.ReadCourses(con.UserName, con.Password, query.GetRes(con.UserName, con.Password, query_url+"/myCourses.asp"))

		for k, v := range myCourses {
			if in(k.Title, yearArr) {
				study.StartStudy(con.UserName, con.Password, v, &Wg)
			}

			// 考试

		}

		time.Sleep(time.Duration(con.DoneTIme) * time.Second)
		fmt.Println("【", con.UserName, "】学习结束")
	}
	Wg.Wait()

	fmt.Println("学习结束")
}

// 测试进度条
func testProgress() {
	// 更新当前的状态
	fmt.Println("开始任务... ")

	var total = 30
	for i := 0; i < total; i++ {
		time.Sleep(time.Duration(rand.Intn(400)) * time.Millisecond)

		// 计算百分比
		percent := int(float32(i+1) * 100.0 / float32(total))

		pro := Progress(percent)
		pro.Show()
	}
	fmt.Printf("\nDone!\n")
}

// Progress 进度
type Progress int

// Show 显示进度
func (x Progress) Show() {
	percent := int(x)
	// fmt.Println("percent: ", percent)

	total := 50 // 这个total是格子数
	middle := int(percent * total / 100.0)
	// fmt.Printf("middle:%d\n", middle)

	arr := make([]string, total)
	for j := 0; j < total; j++ {
		if j < middle-1 {
			arr[j] = "-"
		} else if j == middle-1 {
			arr[j] = ">"
		} else {
			arr[j] = " "
		}
	}
	bar := fmt.Sprintf("[%s]", strings.Join(arr, ""))
	fmt.Printf("\r%s %%%d", bar, percent)
}

func in(target string, yearArr []string) bool {
	for _, year := range yearArr {
		if strings.Contains(target, year) {
			return true
		}
	}
	return false
}
