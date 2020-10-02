package util

import "testing"

func TestTester_TestCases(t *testing.T) {
	tester := NewTester(t, nil, nil)
	tester.SetUp()
	tester.Expect("1", "1")
	tester.Expect("1", "12")
	tester.Expect("1", "3")
	tester.TestCases()
	tester.TearDown()
}
