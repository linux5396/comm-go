package container

import (
	"github.com/linux5396/comm-go/util"
	"testing"
	"time"
)

func TestNewTimerCache(t *testing.T) {
	tester := util.NewTester(t, nil, nil)
	tester.SetUp()
	//----
	tc := NewTimerCache()
	tc.Put("name", "linxu", 2)
	time.Sleep(1 * time.Second)
	v, err := tc.Get("name")
	tester.Expect(err, nil)
	tester.Expect(v, "linxu")
	time.Sleep(2 * time.Second)
	v1, _ := tc.Get("name")
	tester.Expect(v1, nil)
	time.Sleep(10 * time.Second) //util trigger
	tester.ExpectInt(tc.Size(), 0)
	tc.Destroy()
	//----
	tester.ShowTestCases()
	tester.TearDown()
}
