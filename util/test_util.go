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

//serial obj to json string
//if occurred error,will throws panic
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

//compare after serialization
//compare focus on data
func (t *Tester) ExpectInJson(actual interface{}, expect interface{}) bool {
	t.unitNo++
	res := false
	a, e := t.JsonSerialToString(actual), t.JsonSerialToString(expect)
	if a == e {
		t.Logf("%s", "pass")
		res = true
	} else {
		t.Logf("expect is %v but actual is %v", a, e)
		res = false
	}
	t.result[t.unitNo] = res
	return res
}

//show all test cases
//only support expectXX test methods
//not for assert
func (t *Tester) ShowTestCases() {
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

//itf
func (t *Tester) Assert(actual interface{}, expect interface{}) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}

//assert int
func (t *Tester) AssertInt(actual int, expect int) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}
func (t *Tester) AssertInt8(actual int8, expect int8) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}
func (t *Tester) AssertInt16(actual int16, expect int16) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}
func (t *Tester) AssertInt32(actual int32, expect int32) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}
func (t *Tester) AssertInt64(actual int64, expect int64) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}

//float
func (t *Tester) AssertFloat32(actual float32, expect float32) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}
func (t *Tester) AssertFloat64(actual float64, expect float64) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}

//string
func (t *Tester) AssertString(actual string, expect string) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}

//byte
func (t *Tester) AssertByte(actual byte, expect byte) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}

//rune
func (t *Tester) AssertRune(actual rune, expect rune) {
	if actual == expect {
		t.Logf("%s", "pass")
	} else {
		t.Fatalf("expect is %v but actual is %v", expect, actual)
	}
}
