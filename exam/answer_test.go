package exam

import (
	"fmt"
	"goFastForVideo/query"
	"strconv"
	"testing"
	"time"
)

var userName = "320323196910060022"
var password = "123456"
var query_url = "http://plat.xzjxjy.com"

func TestAnswer(t *testing.T) {
	courses := query.ReadCourses(userName, password, query.GetRes(userName, password, query_url+"/myCourses.asp"))
	for course, _ := range courses {
		for {
			score, questionBankList := Answer(userName, password, course)
			fmt.Println("score", score)
			scoreNumber, err := strconv.Atoi(score)
			Check(err)
			if scoreNumber > PASS_SCORE {
				break
			}
			SaveQuestions(questionBankList, "../resources/question_bank/"+course.Title+".json")
			time.Sleep(2 * time.Second)
		}

	}

}
