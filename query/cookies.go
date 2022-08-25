package query

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var loginUrl = "http://plat.xzjxjy.com/member_login.asp?platID=FE84A5103525778469ADC0D72967BFD40F6E385C179A7CE3F4E3C5D1555CE33E"
var cookieMap = make(map[string]string)

var UserId = ""

func RemoveCookie(userName string) {
	delete(cookieMap, userName)
}

func GetCookie(userName, password string) string {
	// fmt.Println(cookieMap[userName])
	if "" != cookieMap[userName] {
		return cookieMap[userName]
	}

	client := &http.Client{}

	u, _ := url.ParseRequestURI(loginUrl)
	urlStr := u.String() // "http://127.0.0.1/tpost"

	body, _ := Utf8ToGbk("loginname=" + userName + "&password=" + password)

	reader := strings.NewReader(body)
	r, _ := http.NewRequest("POST", urlStr, reader) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	r.Header.Add("Cache-Control", "max-age=0")
	r.Header.Add("Upgrade-Insecure-Requests", "1")
	r.Header.Add("Origin", "http://plat.xzjxjy.com")
	r.Header.Add("Accept-Encoding", "gzip, deflate")
	r.Header.Add("Connection", "keep-alive")
	r.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
	res, err := client.Do(r)
	// s, _ := ioutil.ReadAll(res.Body) //把	body 内容读入字符串 s
	// fmt.Printf("%s", s)              //在返回页面中显示内容。
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

// UTF-8 转 GBK
func Utf8ToGbk(s string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return "nil", e
	}
	// fmt.Printf("%s", d)
	return string(d), nil
}
