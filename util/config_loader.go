package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

//configuration loader which read by file stream and read per line to parse and make cache.
//each conf file should new a confLoader to load the configuration for using.
type ConfLoader struct {
	conf map[string]map[string]string
}

func NewConfLoader() *ConfLoader {
	return &ConfLoader{conf: make(map[string]map[string]string)}
}

func (loader *ConfLoader) LoadConf(confPath string) {
	//keep strong
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("load conf occurred panic:%v", r)
		}
	}()
	file, err := os.Open(confPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	sect := ""
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		txtLine := string(line)
		if strings.HasPrefix(txtLine, "[") && strings.HasSuffix(txtLine, "]") {
			//is sect
			sect = txtLine[1 : len(txtLine)-1]
			loader.conf[sect] = make(map[string]string)
		} else if strings.HasPrefix(txtLine, "#") {
			continue
		} else {
			kv := strings.Split(txtLine, "=")
			if kv != nil && len(kv) == 2 {
				loader.conf[sect][kv[0]] = kv[1]
				continue
			}
		}
	}
}

//it may panic
func (loader *ConfLoader) GetSect(sect string) map[string]string {
	return loader.conf[sect]
}

//it may panic
func (loader *ConfLoader) GetValue(sect, key string) string {
	return loader.conf[sect][key]
}

//get or default
func (loader *ConfLoader) GetValueOrDefault(sect, key, def string) string {
	v := loader.conf[sect][key]
	if v == "" {
		return def
	}
	return v
}

//get int val
func (loader *ConfLoader) GetValueIntOrDefault(sect, key string, def int) int {
	val := loader.conf[sect][key]
	res, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return res
}
