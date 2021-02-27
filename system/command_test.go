package system

import "testing"

func TestExecute(t *testing.T) {
	//1>&2 tests outstream redict
	out, errOut, err := Execute("ls  1>&2")
	if err != nil {
		t.Fatal(err)
	}
	if out != nil {
		t.Logf("%v", string(out))
	}
	if errOut != nil {
		t.Logf("err:\n%v", string(errOut))
	}
}
