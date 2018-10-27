package model

type Uploader interface {
	Driver() string
	Upload(task Task, broadcaster chan Task) (bool, error)
}
