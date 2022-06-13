package file

import (
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Config struct {
	UserName string
	Password string
	Years    string
	DoneTIme int
}

func ReadSetting() []Config {
	f, err := excelize.OpenFile("settings.xlsx")
	Check(err)

	rows, err := f.GetRows("Sheet1")

	configs := []Config{}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		doneTIme, _ := strconv.Atoi(row[3])
		con := Config{
			UserName: row[0],
			Password: row[1],
			Years:    row[2],
			DoneTIme: doneTIme,
		}
		configs = append(configs, con)

	}

	return configs

}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
