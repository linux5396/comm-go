package container

//iterator interface
type Foreach interface {
	Foreach(f func(arg interface{}))
}
