package container

//iterator interface
type Foreach interface {
	Foreach(f func(arg interface{}))
}

//provide a unified slicing iterator
type SliceIterator struct {
	slice []interface{}
}

//call ways:
//NewS..(slice).Foreach(f)
func NewSliceIterator(slice []interface{}) SliceIterator {
	return SliceIterator{slice: slice}
}

//impl
func (s SliceIterator) Foreach(f func(arg interface{})) {
	for _, v := range s.slice {
		f(v)
	}
}
