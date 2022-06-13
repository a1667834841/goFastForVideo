package pass

import "testing"

func TestPasst(t *testing.T) {
	passes := ReadLine()
	for i := 0; i < len(passes); i++ {
		Login(passes[i])
	}
}
