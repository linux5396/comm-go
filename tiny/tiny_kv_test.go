package tiny

import (
	"github.com/linux5396/comm-go/util"
	"testing"
)

var kv *TinyKV

func TestNewTinyKV(t *testing.T) {
	cf := util.NewConfLoader()
	cf.LoadConf("/home/antonlin/go_work_space/comm-go/tiny/tiny.conf")
	p := cf.GetValueOrDefault("Tiny", "DataPath", "1")
	t.Log(p)
	kv = NewTinyKV(3, cf)
	err := kv.Put("hello", 1)
	if err != nil {
		t.Fatal(err)
	}
}
