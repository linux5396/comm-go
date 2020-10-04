package util

import "testing"

//expect测试
func TestTester_TestCases(t *testing.T) {
	tester := NewTester(t, nil, nil)
	tester.SetUp()
	tester.Expect("1", "1")
	tester.Expect("1", "12")
	tester.Expect("1", "3")
	tester.ExpectFloat64(0.1, 0.2)
	tester.ExpectString("hello", "hello world")
	type TestObj struct {
		Val int `json:"val"`
		Key int `json:"key"`
	}
	tester.JsonSerialToLog(TestObj{
		Val: 1,
		Key: 3,
	})
	tester.ShowTestCases()
	tester.TearDown()
}

//assert测试
func TestTester_Assert(t *testing.T) {
	tester := NewTester(t, nil, nil)
	tester.Assert("1", 2)
	tester.TearDown()
}
