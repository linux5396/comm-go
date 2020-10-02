package util

import (
	"encoding/json"
	"testing"
)

//a tester tool
type Tester struct {
	*testing.T
	setup    func()
	tearDown func()
	//for suit
	unitNo int
	result map[int]bool //inner map , so should using the ptr as receiver
}

//new a tester
func NewTester(ts *testing.T, setUp func(), tearDown func()) *Tester {
	t := &Tester{
		T: ts,
	}
	t.setup = func() {
		if setUp != nil {
			setUp()
		}
		t.unitNo = 0
		t.result = make(map[int]bool)
	}
	t.tearDown = func() {
		if tearDown != nil {
			tearDown()
		}
		t.unitNo = 0
		for k := range t.result {
			delete(t.result, k)
		}
	}
	return t
}

func (t *Tester) SetUp() {
	if t.setup != nil {
		t.setup()
	}
}

func (t *Tester) TearDown() {
	if t.tearDown != nil {
		t.tearDown()
	}
}

//
func (t *Tester) ExpectInt(actual int, expect int) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}
func (t *Tester) ExpectInt8(actual int8, expect int8) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}
func (t *Tester) ExpectInt16(actual int16, expect int16) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}
func (t *Tester) ExpectInt32(actual int32, expect int32) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}
func (t *Tester) ExpectInt64(actual int64, expect int64) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}

//
func (t *Tester) ExpectFloat32(actual float32, expect float32) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}
func (t *Tester) ExpectFloat64(actual float64, expect float64) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}

//
func (t *Tester) ExpectByte(actual byte, expect byte) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}

//
func (t *Tester) ExpectRune(actual rune, expect rune) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}

//
func (t *Tester) ExpectString(actual string, expect string) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}

//
func (t *Tester) Expect(actual interface{}, expect interface{}) bool {
	t.unitNo++
	res := false
	if actual == expect {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", expect, actual)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}

func (t *Tester) JsonSerialToString(obj interface{}) string {
	if obj != nil {
		jsonObj, err := json.MarshalIndent(obj, "", "\t")
		if err != nil {
			panic(err)
		}
		return string(jsonObj)
	}
	return ""
}

func (t *Tester) JsonSerialToLog(obj interface{}) {
	t.Logf("%v", t.JsonSerialToString(obj))
}

//show all test cases
func (t *Tester) TestCases() {
	pass := 0
	for k := range t.result {
		if t.result[k] {
			pass++
		}
	}
	t.Log("------test result -----------")
	t.Logf("| count: %d , pass: %d ,fail: %d  , pass_rate: %f%% \n", len(t.result), pass, len(t.result)-pass, float64(pass)*100/float64(len(t.result)))
	for k := range t.result {
		if t.result[k] {
			t.Logf("| CASE : %d  pass  \n", k)
		} else {
			t.Logf("| CASE : %d  fail  \n", k)
		}
	}
	t.Log("------end test result -------")
}

//TODO impl assert
func Assert() {

}
