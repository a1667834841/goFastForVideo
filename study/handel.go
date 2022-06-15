package study

import (
	"fmt"
	"goFastForVideo/query"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fastjson"
)

var chcode = ""

func init() {

}

func StartStudy(userName string, password string, videos []query.Video, Wg *sync.WaitGroup) {
	fmt.Println("开始学习。。。")
	for _, video := range videos {

		if video.Progress == "100%" {
			continue
		}
		Wg.Add(1)
		go learnServer(userName, password, video, Wg)

	}

}

func learnServer(userName, password string, video query.Video, Wg *sync.WaitGroup) {
	client := &http.Client{}

	for {

		study := "Action=" + video.Study.Action + "&Platid=" + video.Study.Platid + "&Userid=" + video.Study.Userid + "&Courseid=" + video.Study.Courseid + "&Coursewareid=" + video.Study.Coursewareid + "&St=" + video.Study.St
		req, _ := http.NewRequest(http.MethodPost, video.Url, strings.NewReader(study))
		req.Header.Set("Cookie", query.GetCookie(userName, password))
		req.Header.Set("Content-Length", strconv.Itoa(len(study)))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		state := fastjson.GetInt(body, "result")
		percent := fastjson.GetString(body, "progress")
		msg := fastjson.GetString(body, "msg")

		if state != 1 {
			fmt.Println("study:", study, "title:", video.Title, ",msg:", msg)
			continue
		}

		fmt.Println("title:", video.Title, ",percent:", percent)

		if percent == "100%" {
			fmt.Println("title:", video.Title, " 观看结束，percent:", percent)
			for i := 0; i < 10; i++ {
				time.Sleep(100 * time.Millisecond)
				client.Do(req)
			}
			video.Progress = "100%"
			Wg.Done()
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

}
