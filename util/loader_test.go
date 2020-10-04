package util

import "testing"

func TestConfLoader_LoadConf(t *testing.T) {
	loader := NewConfLoader()
	loader.LoadConf("E:\\golang_projects\\comm-go\\client\\redis_cli.conf")
	m := loader.GetSect("standalone")
	t.Logf("%v", m)
}
