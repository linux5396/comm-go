package container

import (
	"errors"
	"log"
	"runtime"
	"sync"
	"time"
)

var entryNotExist = errors.New("this entry does not exist")
var entryExpired = errors.New("entry expired")

//timerCache can provide a feature that kv storage in expiration.
//default evict routine frequency is t/10S
//it is safety for multi threads rw
//use cases:
//1. verify code 、web token、session , etc...
//2. cache data,like some simple articles.
//...
type TimerCache struct {
	//read write lock to get fast
	sync.RWMutex
	//hash map as storage struct
	cache map[interface{}]interface{}
	//record every kv 's deadline
	deadLine map[interface{}]time.Time
	//recent timerEvent
	trigger chan time.Time
	//recent event
	latestTimePoint time.Time
	//close
	close    chan int8
	stopLoop chan int8
}

func NewTimerCache() *TimerCache {
	tc := &TimerCache{
		RWMutex:         sync.RWMutex{},
		cache:           make(map[interface{}]interface{}),
		deadLine:        make(map[interface{}]time.Time),
		trigger:         make(chan time.Time),
		latestTimePoint: time.Time{},
		close:           make(chan int8),
		stopLoop:        make(chan int8),
	}
	//register gc hook
	runtime.SetFinalizer(tc, finalizer)
	//evict events
	tc.registerTrigger()
	tc.timeEventLoop()
	return tc
}

//only stop the world
//the mem collection is gave to GC
func (tc *TimerCache) Destroy() {
	tc.stopLoop <- 1 //signal to close all channels and stop go routines
}

//expire is second
//if expire lt 1, default is ten years
func (tc *TimerCache) Put(key interface{}, val interface{}, expireSec int64) {
	formatDate := time.Now()
	if expireSec < 1 {
		//default permanent
		formatDate = formatDate.Add(time.Hour * 24 * 3600)
	} else {
		formatDate = formatDate.Add(time.Second * time.Duration(expireSec))
	}
	tc.Lock()
	defer tc.Unlock()
	tc.cache[key] = val
	tc.deadLine[key] = formatDate
	if tc.latestTimePoint.IsZero() {
		tc.latestTimePoint = formatDate
	}
	if !tc.latestTimePoint.IsZero() && tc.latestTimePoint.After(formatDate) {
		tc.latestTimePoint = formatDate
	}
}

func (tc *TimerCache) Delete(key interface{}) {
	tc.Lock()
	defer tc.Unlock()
	delete(tc.cache, key)
	delete(tc.deadLine, key)
}

func (tc *TimerCache) Size() int {
	tc.Lock()
	defer tc.Unlock()
	return len(tc.cache)
}

func (tc *TimerCache) Get(key interface{}) (interface{}, error) {
	tc.RLock()
	defer tc.RUnlock()
	v, ok := tc.cache[key]
	if !ok {
		return nil, entryNotExist
	}
	if tc.deadLine[key].Before(time.Now()) || tc.deadLine[key].Equal(time.Now()) {
		//because Get ops get the rlock , so it can not rm the expired entry
		return nil, entryExpired
	}
	return v, nil
}

func (tc *TimerCache) Keys() []interface{} {
	tc.RLock()
	defer tc.RUnlock()
	ks := make([]interface{}, 0)
	for k := range tc.cache {
		ks = append(ks, k)
	}
	return ks
}

func (tc *TimerCache) Values() []interface{} {
	tc.RLock()
	defer tc.RUnlock()
	vs := make([]interface{}, 0)
	for _, v := range tc.cache {
		vs = append(vs, v)
	}
	return vs
}

//gc hook
func finalizer(r *TimerCache) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	if r != nil {
		r.close <- 1 // close the time event
	}
}

//time event loop
func (tc *TimerCache) timeEventLoop() {
	go func() {
		for {
			select {
			case <-tc.stopLoop:
				close(tc.stopLoop)
				return
			case t := <-tc.trigger:
				//should get lock before the judgment
				tc.Lock()
				//exec evict
				if t.Equal(time.Now()) || t.Before(time.Now()) {
					//evict job
					for key, dl := range tc.deadLine {
						//can evict, delete ops is safety when iterating.
						if !dl.IsZero() && dl.Before(time.Now()) {
							delete(tc.cache, key)
							delete(tc.deadLine, key)
						}
					}
				}
				tc.Unlock()
			}
		}
	}()
}

//register trigger for time event
func (tc *TimerCache) registerTrigger() {
	go func() {
		for {
			select {
			case <-tc.close:
				close(tc.trigger)
				close(tc.close)
				tc.stopLoop <- 1
				return
			default:
				time.Sleep(10 * time.Second)
				if !tc.latestTimePoint.IsZero() {
					tc.trigger <- tc.latestTimePoint //trigger
				}
			}
		}
	}()
}
