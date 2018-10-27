package model

type Interface interface {
	Driver() string
	Start()
}
