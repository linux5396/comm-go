//+build linux

package tiny

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/linux5396/comm-go/util"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

//tiny kv is a mapping cache with disk persistence
//but it does not achieve strong consistency between memory and disk key-value data
//its application scenarios are below:
//	suitable for fixed key value data caching.This means that there is not much need to perform persistence.
//	suitable for loading the last cached data at every startup.Especially those key-value data with high computational cost.
//why is tiny——extreme scene demand.
//	tiny also support "SegmentedMap"
//	Segmented cache is essentially a multi-instance maps.
//	It can reduce the number of concurrent locks under the premise of ensuring concurrency safety,
//	and can load part of the data according to the index file to improve the speed of loading data.
//Note that as long as the same hash algorithm is implemented,
//the data in the data directory can be deserialized and loaded onto the corresponding segment.Regardless of whether the number of segments is the same.

const Tiny = "[TinyKV]"

type Snapshot interface {
	Write() error
	Read() error
	ClearUp() error
}

type KV interface {
	//get
	Get(key string) (interface{}, error)
	//put
	Put(key string, value interface{}) error
	//evict
	Evict(key string) error
	//size
	Size() int
}

//TinyKV
type TinyKV struct {
	//segmentation lock
	rwLockers []*sync.RWMutex
	//bucketSize
	bucketSize int
	//hash function
	hash func(key string) int
	//inner buckets
	buckets []map[string]interface{}
	//data path: 1.snapshot wal; 2.sync index of buckets;
	//conf loader loads all conf of tiny kv
	conf *util.ConfLoader
}

//full file path init
//bucket size = 1<<bucketSizeShift
func NewTinyKVWithConf(bucketSizeShift int, fp string) *TinyKV {
	c := util.NewConfLoader()
	if fileExists(fp) {
		c.LoadConf(fp)
	} else {
		panic(fmt.Sprintf("%s file{%s} not exist", Tiny, fp))
	}
	return NewTinyKV(bucketSizeShift, c)
}

//bucket size = 1<<bucketSizeShift
func NewTinyKV(bucketSizeShift int, conf *util.ConfLoader) *TinyKV {
	if bucketSizeShift < 0 {
		panic("bucketSize left shift can not lt 0 ")
	}
	if conf == nil {
		panic("conf cannot be nil")
	}
	kv := &TinyKV{
		bucketSize: 1 << bucketSizeShift,
		hash:       crc,
		buckets:    make([]map[string]interface{}, 1<<bucketSizeShift),
		conf:       conf,
	}
	//init
	for i := 0; i < kv.bucketSize; i++ {
		kv.buckets[i] = make(map[string]interface{})
		kv.rwLockers = append(kv.rwLockers, &sync.RWMutex{})
	}
	return kv
}

//bucket size = 1<<bucketSizeShift
func NewTinyKVWithHash(bucketSizeShift int, conf *util.ConfLoader, hash func(k string) int) *TinyKV {
	if hash == nil {
		panic("hash function is nil")
	}
	kv := NewTinyKV(bucketSizeShift, conf)
	kv.hash = hash
	return kv
}

func (tiny *TinyKV) Get(key string) (interface{}, error) {
	idx := tiny.hash(key) % tiny.bucketSize
	tiny.rwLockers[idx].RLock()
	defer func() {
		tiny.rwLockers[idx].RUnlock()
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(os.Stderr, Tiny, r)
		}
	}()
	return tiny.buckets[idx][key], nil
}

func (tiny *TinyKV) Put(key string, value interface{}) error {
	idx := tiny.hash(key) % tiny.bucketSize
	tiny.rwLockers[idx].Lock()
	defer func() {
		tiny.rwLockers[idx].Unlock()
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(os.Stderr, Tiny, r)
		}
	}()
	tiny.buckets[idx][key] = value
	return nil
}

func (tiny *TinyKV) Evict(key string) error {
	idx := tiny.hash(key) % tiny.bucketSize
	tiny.rwLockers[idx].Lock()
	defer func() {
		tiny.rwLockers[idx].Unlock()
		if r := recover(); r != nil {
			_, _ = fmt.Fprintln(os.Stderr, Tiny, r)
		}
	}()
	delete(tiny.buckets[idx], key)
	return nil
}

//lazy
func (tiny *TinyKV) Size() int {
	cnt := 0
	for i := 0; i < tiny.bucketSize; i++ {
		cnt += len(tiny.buckets[i])
	}
	return cnt
}

//default hash function
func Hash(key string) int {
	return crc(key)
}

//crc
func crc(key string) int {
	v := crc32.ChecksumIEEE([]byte(key))
	return int(v)
}

//f exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//write tiny kv data to disk
func (tiny *TinyKV) Write() error {
	ch := make(chan error, tiny.bucketSize) //async chan avoid blocked
	defer close(ch)
	wg := sync.WaitGroup{}
	wg.Add(tiny.bucketSize)
	for i := 0; i < tiny.bucketSize; i++ {
		idx := i
		go func() {
			tiny.rwLockers[idx].RLock()
			defer tiny.rwLockers[idx].RUnlock()
			defer wg.Done()
			if len(tiny.buckets[idx]) == 0 {
				return //quit
			}
			buf, err := json.Marshal(tiny.buckets[idx]) //json 序列化
			if err != nil {
				ch <- err
				return
			}
			f, err := os.OpenFile(fmt.Sprintf(tiny.conf.GetValue("Tiny", "DataPath")+"/tiny_%d.d", idx), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				ch <- err
				return
			}
			n, err := f.Write(buf)
			if n < len(buf) {
				if err != nil {
					ch <- err
					return
				}
				ch <- errors.New("short write err") //this
			}
			err = tiny.writeIndexes(idx) //write idx
			if err != nil {
				ch <- err
				return
			}
			defer f.Close()
			_ = f.Sync() //flush to disk
		}()
	}
	wg.Wait() //wait for all goroutines done
	select {
	case err := <-ch: //check error
		return err
	default:
		return nil
	}
}

//write index file
//index file : record the mapping of index and data'files
func (tiny *TinyKV) writeIndexes(i int) error {
	path := tiny.conf.GetValue("Tiny", "DataPath")
	f, err := os.Create(fmt.Sprintf("%s/tiny_%d.idx", path, i))
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

func (tiny *TinyKV) Read() error {
	path := tiny.conf.GetValue("Tiny", "DataPath")
	for i := 0; i < tiny.bucketSize; i++ {
		if fileExists(fmt.Sprintf("%s/tiny_%d.idx", path, i)) {
			curPath := fmt.Sprintf("%s/tiny_%d.d", path, i)
			buf, err := ioutil.ReadFile(curPath)
			if err != nil {
				return err
			}
			err = json.Unmarshal(buf, &tiny.buckets[i])
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s unmarshal datafile %s err:{%s}\n", Tiny, curPath, err.Error())
				return err
			}
		}
	}
	return nil
}

//snapshot clear
func (tiny *TinyKV) ClearUp() error {
	path := tiny.conf.GetValue("Tiny", "DataPath")
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		e := os.Remove(path)
		if e == nil {
			_, _ = fmt.Fprintf(os.Stdout, "%s removed (%s) in ts(%d)\n", Tiny, path, time.Now().Unix())
		}
		return e
	})
}
func (tiny *TinyKV) ClearUpOverTime() error {
	path := tiny.conf.GetValue("Tiny", "DataPath")
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		//reserve
		if time.Now().Unix()-info.ModTime().Unix()-int64(tiny.conf.GetValueIntOrDefault("Tiny", "ReservedInterval", 0)) > 0 {
			//
			return nil
		}
		e := os.Remove(path)
		if e == nil {
			_, _ = fmt.Fprintf(os.Stdout, "%s removed (%s) in ts(%d)\n", Tiny, path, time.Now().Unix())
		}
		return e
	})
}
