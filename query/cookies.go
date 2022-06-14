package query

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var loginUrl = "http://plat.xzjxjy.com/member_login.asp"
var cookieMap = make(map[string]string)

var UserId = ""

func GetCookie(userName, password string) string {
	// fmt.Println(cookieMap[userName])
	if "" != cookieMap[userName] {
		return cookieMap[userName]
	}

	reader := strings.NewReader("loginname=" + userName + "&password=" + password)
	res, err := http.Post(loginUrl, " application/x-www-form-urlencoded", reader)
	Check(err)
	if res.Request.Response == nil {
		fmt.Println(userName, "登陆失败")
		panic("登陆失败")
	}
	// fmt.Println(res.Request.Response.Header.Values("Set-Cookie"))
	setCooikes := res.Request.Response.Header.Values("Set-Cookie")

	var cookie = ""
	for _, item := range setCooikes {
		item = strings.Replace(item, "path=/;", "", -1)
		if strings.Contains(item, "&id=") {
			UserId = strings.Split(strings.Split(item, "&id=")[1], "&")[0]
		}
		cookie += item
	}

	// fmt.Println(
	// 	strings.Replace(strings.Replace(cookie, "path=/", "", -1), "Max-Age=1800", "", -1),
	// )
	cookie = strings.Replace(strings.Replace(cookie, "path=/", "", -1), "Max-Age=1800", "", -1)
	cookieMap[userName] = cookie
	return cookie

}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// GetBeforeTime 获取n天前的秒时间戳、日期时间戳
// _day为负则代表取前几天，为正则代表取后几天，0则为今天
func GetBeforeTime(_day int) (int64, string) {
	// 时区
	//timeZone, _ := time.LoadLocation(ServerInfo["timezone"])
	timeZone := time.FixedZone("CST", 8*3600) // 东八区

	// 前n天
	nowTime := time.Now().In(timeZone)
	beforeTime := nowTime.AddDate(0, 0, _day)

	// 时间转换格式
	beforeTimeS := beforeTime.Unix()                                 // 秒时间戳
	beforeDate := time.Unix(beforeTimeS, 0).Format("20060102150405") // 固定格式的日期时间戳

	return beforeTimeS, beforeDate
}
