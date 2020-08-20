package container

//Set itf
type Set interface {
	//foreach is call f to process all elements in set.
	Foreach
	//put
	Put(key interface{})
	//Determine if it exists
	Contains(key interface{}) bool
	//Weed out
	Evict(key interface{})
	//dump all elements from set
	Dump() []interface{}
	//set size
	Size() int
	//set is empty or not
	Empty() bool
	//support intersection search
	Union(set Set) []interface{}
}

//Implementation based on hash table
type HashSet struct {
	size     int
	innerMap map[interface{}]interface{}
}

//Default capacity
const defaultCapacity = 1 << 3

//Construct a hash set of specified capacity
func NewHashSet(initCapacity int) *HashSet {
	//not illegal
	if initCapacity < defaultCapacity {
		initCapacity = defaultCapacity
	}
	return &HashSet{
		size:     0,
		innerMap: make(map[interface{}]interface{}, initCapacity),
	}
}

//Construct a hash set
func MakeHashSet() *HashSet {
	return NewHashSet(defaultCapacity)
}

func NewHashSetByHashMap(replicas map[interface{}]interface{}) *HashSet {
	return &HashSet{
		size:     len(replicas),
		innerMap: replicas,
	}
}

func NewHashSetBySlice(slice []interface{}) *HashSet {
	set := NewHashSet(len(slice))
	//If inline optimization is enabled, the go compiler will automatically inline，no need to care about put call loss
	for _, v := range slice {
		set.Put(v)
	}
	return set
}

func (h *HashSet) Put(key interface{}) {
	h.innerMap[key] = 1
	h.size++
}

//determine whether to include
func (h *HashSet) Contains(key interface{}) bool {
	if _, ok := h.innerMap[key]; ok {
		return ok
	}
	return false
}

//weed out
func (h *HashSet) Evict(key interface{}) {
	delete(h.innerMap, key)
	h.size--
}

//foreach support to traverse all keys and execute f function
func (h *HashSet) Foreach(f func(val interface{})) {
	for k := range h.innerMap {
		f(k)
	}
}

//Dump all keys
func (h *HashSet) Dump() []interface{} {
	keySet := make([]interface{}, 0)
	for k := range h.innerMap {
		keySet = append(keySet, k)
	}
	return keySet
}

//the actual amount
func (h *HashSet) Size() int {
	return h.size
}

func (h *HashSet) Empty() bool {
	if h.size == 0 {
		return true
	}
	return false
}

//be mixed
//basic algorithm：use small sets to find intersections in large sets
//assume SetA(N)  SetB(M)  if N < M , the time cost is N
//Use the Set interface as a parameter to realize the intersection calculation of multiple sets
func (h *HashSet) Union(set Set) []interface{} {
	res := make([]interface{}, 0)
	if h.size < set.Size() {
		h.Foreach(func(val interface{}) {
			if set.Contains(val) {
				res = append(res, val)
			}
		})
	}
	set.Foreach(func(val interface{}) {
		if h.Contains(val) {
			res = append(res, val)
		}
	})
	return res
}
