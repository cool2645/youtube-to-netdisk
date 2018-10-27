package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"time"
)

type Task struct {
	ID          int64 `gorm:"AUTO_INCREMENT"`
	Title       string
	Description string
	Author      string
	URL         string
	State       string
	State2      string
	Reason      string
	FileName    string
	ShareLink   string
	Log         string `sql:"type:text;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func CreateTask(db *gorm.DB, task *Task) (err error) {
	task.ID = 0
	err = db.Create(task).Error
	if err != nil {
		err = errors.Wrap(err, "CreateTask")
		return
	}
	return
}

func GetTasks(db *gorm.DB, order string, page uint, perPage uint) (tasks []Task, total uint, err error) {
	noLog := "id, title, description, author, url, state, reason, file_name, share_link, created_at, updated_at"
	if perPage == 0 {
		perPage = 10
	}
	err = db.Where("state <> ?", "Rejected").
		Order("updated_at " + order).Limit(perPage).Offset((page - 1) * perPage).
		Select(noLog).Find(&tasks).Error
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	err = db.Model(&Task{}).Where("state <> ?", "Rejected").Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	return
}

func GetQueuingTasks(db *gorm.DB) (tasks []Task, err error) {
	err = db.Where("state = ?", "Queuing").Find(&tasks).Error
	if err != nil {
		err = errors.Wrap(err, "GetQueuingTasks")
		return
	}
	return
}

func GetRejTasks(db *gorm.DB, order string, page uint, perPage uint) (tasks []Task, total uint, err error) {
	noLog := "id, title, description, author, url, state, reason, file_name, share_link, created_at, updated_at"
	if perPage == 0 {
		perPage = 10
	}
	err = db.Where("state = ?", "Rejected").
		Order("updated_at " + order).Limit(perPage).Offset((page - 1) * perPage).
		Select(noLog).Find(&tasks).Error
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	err = db.Model(&Task{}).Where("state = ?", "Rejected").Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	return
}

func GetTask(db *gorm.DB, id int64) (task Task, err error) {
	err = db.Where("id = ?", id).Find(&task).Error
	if err != nil {
		err = errors.Wrap(err, "GetTask")
		return
	}
	return
}

func SaveTask(db *gorm.DB, task *Task) (err error) {
	err = db.Model(task).Updates(*task).Error
	if err != nil {
		err = errors.Wrap(err, "SyncTaskStatus: Save task")
		return
	}
	return
}

func CleanTasks(db *gorm.DB) (err error) {
	err = db.Model(Task{}).Where("state IN (?)", []string{
		"Downloading",
		"Uploading",
	}).Updates(Task{
		State: "Exception",
	}).Error
	if err != nil {
		err = errors.Wrap(err, "CleanTasks")
		return
	}
	return
}
