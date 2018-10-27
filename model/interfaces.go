package model

type Base interface {
	Driver() string
}

type Interface interface {
	Base
	Start()
}

type Broadcaster interface {
	Interface
	Broadcast(task Task)
}

type Uploader interface {
	Base
	Upload(task Task, broadcaster chan Task) (bool, error)
}
