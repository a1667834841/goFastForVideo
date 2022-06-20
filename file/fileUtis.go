package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// 读取文件 返回map
func ReadFileInfo(path string) map[string]string {

	bool, err := PathExists(path)

	if !bool {
		f, err := os.Create(path)
		defer f.Close()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			_, err = f.Write([]byte("{}"))
			Check(err)
		}
	}

	jsonFile, err := os.Open(path)
	Check(err)
	defer jsonFile.Close()

	bytesFile, _ := ioutil.ReadAll(jsonFile)
	// fmt.Println(string(bytesFile))

	pushedMap := make(map[string]string)
	err1 := json.Unmarshal(bytesFile, &pushedMap)
	if err1 != nil {
		panic(err1)
	}
	// fmt.Println("json to map ", pushedMap)
	return pushedMap
}

// 保存已推送文章id 到本地
func WriteFileInfo(temp map[string]string, old map[string]string, path string) {
	for key, value := range temp {
		old[key] = value
	}

	// json 序列化map
	data, _ := json.Marshal(old)

	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}

}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
