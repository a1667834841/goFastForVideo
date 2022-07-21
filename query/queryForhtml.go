package query

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var query_url = "http://plat.xzjxjy.com"

type Course struct {
	Title   string
	Url     string
	ExamUrl string
}
type Video struct {
	Title    string
	Url      string
	Progress string
	Study    Study
}

type Study struct {
	Action       string
	Platid       string
	Userid       string
	Courseid     string
	Coursewareid string
	St           string
	Chcode       string
}

// func Demo() {
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", query_url+"/myCourses.asp", nil)
// 	req.Header.Set("Cookie", GetCookie(userName, password, loginUrl))
// 	res, _ := client.Do(req)

// 	defer res.Body.Close()
// 	body, _ := ioutil.ReadAll(res.Body)
// 	fmt.Println(string(body))
// }

func GetRes(userName, password, url string) *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", GetCookie(userName, password))
	res, _ := client.Do(req)
	// body, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))
	return res
}

func ReadCourses(userName string, password string, res *http.Response) map[Course][]Video {
	fmt.Println("开始爬取【" + userName + "】课程爬取...")
	doc, _ := goquery.NewDocumentFromResponse(res)

	myCourses := make(map[Course][]Video)

	doc.Find("ul[class=pxrw] li a").Each(func(i int, selection *goquery.Selection) {
		// fmt.Println(i)
		href, bool := selection.Attr("href")
		title := selection.Text()
		course := Course{}
		if bool && !strings.Contains(title, "开始考试") && !strings.Contains(title, "进入学习") {

			examUrl := strings.Replace("get_exam.asp?"+strings.Split(href, "?")[1], "id", "courseId", 1)

			course = Course{Title: title, Url: query_url + "/" + href, ExamUrl: query_url + "/" + examUrl}
			myCourses[course] = ReadDetailCourses(userName, password, course)
		}

	})

	// fmt.Println(myCourses)
	fmt.Println("【" + userName + "】课程爬取成功")

	return myCourses
}

func ReadDetailCourses(userName string, password string, course Course) []Video {
	res := GetRes(userName, password, course.Url)

	doc, _ := goquery.NewDocumentFromResponse(res)
	videos := []Video{}

	doc.Find("div[class=ckkc-desc] ul li").Each(func(i int, selection *goquery.Selection) {

		progress := selection.Find("div[class=lanbg] p").Text()
		if progress == "ok" {
			return
		}

		title, bool := selection.Attr("title")
		video := Video{}
		if bool {
			video.Title = title
		}

		selection.Find("p a").Each(func(i int, selection *goquery.Selection) {
			url, b := selection.Attr("href")
			if b {
				video.Url = "http://app.chinahrt.cn/app/courseware/"
				// video.Url = "http://plat.xzjxjy.com/learn_server.asp"
				kcid := strings.Split(strings.Split(url, "kcid=")[1], "&id")[0]
				kjid := strings.Split(url, "&id=")[1]

				video.Study = Study{Action: "studylog", Courseid: kcid, Coursewareid: kjid, St: "80"}
				video.Study = GetUserIdAndPlatId(userName, password, query_url+"/"+url, video.Study)

			}
		})

		selection.Find("div p").Each(func(i int, selection *goquery.Selection) {
			progress := selection.Text()
			video.Progress = progress
		})

		videos = append(videos, video)
	})

	return videos

}

func GetStudy(userName, password, url string) string {
	res := GetRes(userName, password, url)
	doc, _ := goquery.NewDocumentFromResponse(res)

	study := strings.Replace(strings.Replace(doc.Find("script:nth-last-of-type(2)").Last().Text(), "\n", "", -1), "\t", "", -1)
	study = strings.Replace(
		strings.Replace(
			strings.Replace(strings.Replace(strings.Replace(strings.Replace(study, ":", "=", -1), ",", "&", -1), "'", "", -1), ";", "", -1),
			"study={",
			"",
			-1),
		"}",
		"",
		-1)

	return study
}

func GetChcode(userName, password, url string) string {
	res := GetRes(userName, password, url)
	doc, _ := goquery.NewDocumentFromResponse(res)

	study := strings.Replace(strings.Replace(doc.Find("script:nth-last-of-type(2)").Last().Text(), "\n", "", -1), "\t", "", -1)

	study = strings.Split(strings.Split(study, "chcode:'")[1], "',v:")[0]

	return study
}

func GetUserIdAndPlatId(userName, password, url string, study Study) Study {
	// fmt.Println("url:", url)
	res := GetRes(userName, password, url)
	doc, _ := goquery.NewDocumentFromResponse(res)

	studyCode := strings.Replace(strings.Replace(doc.Find("script:nth-last-of-type(3)").Last().Text(), "\n", "", -1), "\t", "", -1)
	if !strings.Contains(studyCode, "study") {
		studyCode = strings.Replace(strings.Replace(doc.Find("script:nth-last-of-type(1)").Last().Text(), "\n", "", -1), "\t", "", -1)
	}

	// fmt.Println(studyCode)
	platid := strings.Split(strings.Split(studyCode, "{platid:")[1], ",gcid")[0]
	chcode := strings.Split(strings.Split(studyCode, ",chcode:")[1], "v:")[0]
	userId := UserId
	study.Platid = platid
	study.Userid = userId
	study.Chcode = chcode
	return study
}
