package query

import "testing"

func TestGetCookie(t *testing.T) {
	var userName = "32030519841022041x"
	var password = "123456"
	GetCookie(userName, password)
}
