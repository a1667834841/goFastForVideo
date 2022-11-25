package exam

import (
	"bytes"
	"fmt"
	"goFastForVideo/file"
	"goFastForVideo/query"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/piex/transcode"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type QuestionBank struct {
	Question string
	Options  []string
	Answer   string
}

type QuestionForm struct {
	ExamId           string
	MemberId         string
	CourseId         string
	ExamQuestionList []ExamQuestion
}

type ExamQuestion struct {
	Question   string
	QuestionId string
	Answer     string
}

var PASS_SCORE = 60

func Answer(userName, password string, course query.Course) (score string, questionBankList []QuestionBank) {

	// 获取题目
	examUrl := course.ExamUrl
	questionForm := getQuestions(userName, password, examUrl)

	// 填充答案
	examQuestionList := fillAnswer(questionForm.ExamQuestionList, course.Title)
	questionForm.ExamQuestionList = examQuestionList
	// 考试提交
	resultUrl, score := examSubmissions(userName, password, questionForm)
	scoreNumber, err := strconv.Atoi(score)
	Check(err)
	if scoreNumber > PASS_SCORE {
		return score, questionBankList
	}

	// 收集题库
	questionBankList = collectionQuestions(resultUrl, userName, password)

	return score, questionBankList

	// 保存题库
	// saveQuestions(questionBankList,)
	// fmt.Println(questionBankList)
}

// 保存题库
func SaveQuestions(questionBankList []QuestionBank, path string) {

	// 读取题库
	old := file.ReadFileInfo(path)

	// 构造map
	temp := make(map[string]string)
	for _, ques := range questionBankList {
		temp[ques.Question] = strings.Join(ques.Options, "|") + ":" + ques.Answer
	}

	// 写入
	file.WriteFileInfo(temp, old, path)
}

// 获取题目
func getQuestions(userName, password, examUrl string) QuestionForm {
	examQuestionList := []ExamQuestion{}
	questionForm := QuestionForm{}
	res := query.GetRes(userName, password, examUrl)
	// defer res.Body.Close()

	// s, _ := ioutil.ReadAll(res.Body) //把	body 内容读入字符串 s
	// resBody := transcode.FromString(string(s)).Decode("UTF-8").ToString()
	// resBody := ConvertByte2String(s, GB18030)
	// fmt.Println(resBody)
	fmt.Println("开始爬取题目")
	doc, _ := goquery.NewDocumentFromReader(res.Body)

	doc.Find("form>input").Each(func(i int, s *goquery.Selection) {
		name, b := s.Attr("name")
		if b && name == "examId" {
			examId, sb := s.Attr("value")
			if sb && examId != "" {
				questionForm.ExamId = examId
			}
		}

		if b && name == "memberId" {
			memberId, sb := s.Attr("value")
			if sb && memberId != "" {
				questionForm.MemberId = memberId
			}
		}

		if b && name == "courseId" {
			courseId, sb := s.Attr("value")
			if sb && courseId != "" {
				questionForm.CourseId = courseId
			}
		}
	})

	doc.Find("div[class=dContent] > div").Each(func(i int, s *goquery.Selection) {

		examQuestion := ExamQuestion{Answer: "1"}
		titleClass, bool := s.Attr("class")
		if bool && titleClass == "question_title" && s.Text() != "" {
			examQuestion.Question = s.Text()
			next := s.Next()
			contentClass, bc := next.Attr("class")
			// fmt.Println(contentClass)
			if bc && contentClass == "question_content" {
				qAnswer, bq := next.Find("div input").Attr("name")
				if bq {
					examQuestion.QuestionId = qAnswer
				}
			}
			examQuestionList = append(examQuestionList, examQuestion)
		}
	})
	questionForm.ExamQuestionList = examQuestionList
	return questionForm
}

// 收集题库
func collectionQuestions(resultUrl, userName, password string) []QuestionBank {

	questionBankList := []QuestionBank{}

	res := query.GetRes(userName, password, "http://plat.xzjxjy.com/"+resultUrl)
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromResponse(res)
	doc.Find("div[class=dContent] > div").Each(func(i int, s *goquery.Selection) {

		questionBank := QuestionBank{}
		titleClass, bool := s.Attr("class")
		if bool && titleClass == "question_title" && s.Text() != "" {
			question := transcode.FromString(string(s.Text())).Decode("UTF-8").ToString()
			questionBank.Question = strings.Split(question, ".")[1]
			next := s.Next()
			contentClass, bc := next.Attr("class")
			// fmt.Println(contentClass)
			if bc && contentClass == "question_content" {

				qAnswer := ""
				next.Find("div input[check=yes]").Each(func(i int, s *goquery.Selection) {
					value, b := s.Attr("value")
					if b {
						qAnswer += "|" + value
					}

				})
				questionBank.Answer = strings.Replace(qAnswer, "|", "", 1)

				next.Find("div label").Each(func(i int, s *goquery.Selection) {
					value := transcode.FromString(string(s.Text())).Decode("UTF-8").ToString()
					questionBank.Options = append(questionBank.Options, value)

				})

			}
			questionBankList = append(questionBankList, questionBank)
		}
	})
	return questionBankList
}

// 考试提交
func examSubmissions(userName, password string, questionForm QuestionForm) (resultUrl string, score string) {
	var overExamUrl = "https://plat.xzjxjy.com/over_exam.asp"

	client := &http.Client{}

	// 用url.values方式构造form-data参数
	// formValues := url.Values{}
	// formValues.Add("examId", questionForm.ExamId)
	// formValues.Add("memberId", questionForm.MemberId)
	// formValues.Add("courseId", questionForm.CourseId)
	// for _, ques := range questionForm.ExamQuestionList {
	// 	if strings.Contains(ques.Answer, "|") {
	// 		answers := strings.Split(ques.Answer, "|")
	// 		for _, answer := range answers {
	// 			formValues.Add(ques.QuestionId, answer)
	// 		}
	// 	} else {
	// 		if "" == ques.Answer {
	// 			formValues.Add(ques.QuestionId, "0")
	// 		} else {
	// 			formValues.Add(ques.QuestionId, ques.Answer)
	// 		}

	// 	}

	// }
	// formBytesReader := strings.NewReader(formValues.Encode())
	// 考试信息
	// 答案
	params := "examId=" + questionForm.ExamId + "&memberId=" + questionForm.MemberId + "&courseId=" + questionForm.CourseId

	for _, ques := range questionForm.ExamQuestionList {
		if strings.Contains(ques.Answer, "|") {
			answers := strings.Split(ques.Answer, "|")
			for _, answer := range answers {
				params += "&" + ques.QuestionId + "=" + answer
			}
		} else {
			if "" == ques.Answer {
				params += "&" + ques.QuestionId + "=1"
			} else {
				params += "&" + ques.QuestionId + "=" + ques.Answer
			}

		}

	}

	formDataBytes := []byte(params)
	formBytesReader := bytes.NewReader(formDataBytes)

	req, _ := http.NewRequest("POST", overExamUrl, formBytesReader)
	req.Header.Set("Cookie", query.GetCookie(userName, password))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Add("Origin", "https://plat.xzjxjy.com")
	req.Header.Add("Host", "plat.xzjxjy.com")
	// req.Header.Add("Connection", "keep-alive")
	// req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,zh-TW;q=0.8,en-US;q=0.7,en;q=0.6")
	// req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36")
	res, _ := client.Do(req)
	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")
	fmt.Println(contentType)
	s, _ := ioutil.ReadAll(res.Body) //把	body 内容读入字符串 s
	// resBody := transcode.FromString(string(s)).Decode("UTF-8").ToString()
	// resBody := ConvertByte2String(s, GB18030)
	resBody := ConvertByte2String(s, GB2312)
	fmt.Println(resBody) //在返回页面中显示内容。

	if strings.Contains(resBody, "请勿非法提交数据") {
		return
	}

	score = strings.Split(strings.Split(resBody, "考试成绩")[1], "分")[0]
	scoreNumber, err := strconv.Atoi(score)
	Check(err)
	if scoreNumber > PASS_SCORE {
		return "", score
	}
	resultUrl = strings.Split(strings.Split(resBody, "='")[2], "';")[0]

	return resultUrl, score
}

// 填充答案
func fillAnswer(examQuestionList []ExamQuestion, courseName string) []ExamQuestion {
	copyies := []ExamQuestion{}
	for k := range examQuestionList {
		quesForm := examQuestionList[k]
		quesFormCopy := ExamQuestion{}

		ques := transcode.FromString(quesForm.Question).Decode("UTF-8").ToString()
		quesFormCopy.Question = ques
		quesFormCopy.QuestionId = quesForm.QuestionId
		// 查询题库
		res, bool := seartchQuestionBank(ques, courseName)
		if bool {
			quesFormCopy.Answer = res
		} else {
			quesFormCopy.Answer = quesFormCopy.Answer
		}
		copyies = append(copyies, quesFormCopy)
	}

	return copyies
}

func seartchQuestionBank(ques, courseName string) (res string, exists bool) {
	bank := file.ReadFileInfo("../resources/question_bank/" + courseName + ".json")
	question := strings.Split(ques, ".")[1]
	value := bank[question]

	if value != "" {
		answer := strings.Split(value, ":")[1]
		return answer, true
	}
	return "", false
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GB2312  = Charset("GB2312")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case GB2312:
		decodeBytes, _ := simplifiedchinese.GBK.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
