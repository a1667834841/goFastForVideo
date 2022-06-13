package query

import (
	"fmt"
	"testing"
)

var userName = "32030519841022041x"
var password = "123456"

func TestDemo(t *testing.T) {
	// Demo()
}

func TestReadCourses(t *testing.T) {
	myCourses := ReadCourses(userName, password, GetRes(userName, password, query_url+"/myCourses.asp"))
	fmt.Println(myCourses)
}

func TestReadDetailCourses(t *testing.T) {
	courses := ReadCourses(userName, password, GetRes(userName, password, query_url+"/myCourses.asp"))
	for k, _ := range courses {
		ReadDetailCourses(userName, password, k)
	}

}

func TestGetStudy(t *testing.T) {
	courses := ReadCourses(userName, password, GetRes(userName, password, query_url+"/myCourses.asp"))
	for _, v := range courses {
		for _, course := range v {
			study := GetStudy(userName, password, course.Url)
			fmt.Println(study)
		}
	}

}

func TestGetChcode(t *testing.T) {
	GetChcode(userName, password, "http://plat.xzjxjy.com/onlineVideo.asp?kcid=283&id=759")
}
