package assert

import (
	"testing"
)

func TestCallerStr(t *testing.T) {
	s := callerStr(0)
	if len(s) <= 0 {
		t.Errorf("return value of callerStr is empty")
	}
}

func TestCallerStrf(t *testing.T) {
	fmtStr := "%d%d"
	val1 := 1
	val2 := 2
	res := callerStrf(0, fmtStr, val1, val2)
	if len(res) < 3 {
		t.Errorf("return value of callerStrf is not long enough")
	}
}
