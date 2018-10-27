package model

type Broadcaster interface {
	Driver() string
	Listen()
	Broadcast(task Task)
}
