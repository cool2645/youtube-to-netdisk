package model

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
)

type Task struct {
	ID          int64  `gorm:"AUTO_INCREMENT"`
	Title       string
	Description string
	Author      string
	URL         string
	State       string
	Reason      string
	FileName    string
	ShareLink   string
	Log         string `sql:"type:text;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(db *gorm.DB, title string, author string, description string, url string, state string, reason string) (newTask Task, err error) {
	var task Task
	task.Title = title
	task.Author = author
	task.Description = description
	task.URL = url
	task.State = state
	task.Reason = reason
	newTask, err = CreateTask(db, task)
	return
}

func CreateTask(db *gorm.DB, task Task) (newTask Task, err error) {
	err = db.Create(&task).Error
	if err != nil {
		err = errors.Wrap(err, "CreateTask")
		return
	}
	newTask = task
	return
}

func GetTasks(db *gorm.DB, state string, order string, page uint, perPage uint) (tasks []Task, total uint, err error) {
	noLog := "id, title, description, author, url, state, reason, file_name, share_link, created_at, updated_at"
	if perPage == 0 {
		perPage = 10
	}
	err = db.Where("state like ?", state).
		Order("updated_at " + order).Limit(perPage).Offset((page - 1) * perPage).
		Select(noLog).Find(&tasks).Error
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	err = db.Model(&Task{}).Where("state like ?", state).Count(&total).Error
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

func UpdateTaskStatus(db *gorm.DB, id int64, state string, fileName string, shareLink string, log string) (err error) {
	var task Task
	err = db.Where("id = ?", id).Find(&task).Error
	if err != nil {
		err = errors.Wrap(err, "SyncTaskStatus: Find task")
		return
	}
	task.State = state
	task.FileName = fileName
	task.ShareLink = shareLink
	task.Log = log
	err = db.Model(&task).Updates(task).Error
	if err != nil {
		err = errors.Wrap(err, "SyncTaskStatus: Update task")
		return
	}
	return
}
