package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"time"
)

type Task struct {
	ID          int64     `gorm:"AUTO_INCREMENT" json:"id"`
	YoutubeID   string    `json:"youtube_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	URL         string    `json:"url"`
	Subtitles   string    `json:"subtitles"`
	State       string    `json:"state"`
	State2      string    `json:"state2"`
	Reason      string    `json:"reason"`
	FileName    string    `json:"file_name"`
	ShareLink   string    `json:"share_link"`
	Log         string    `sql:"type:text;" json:"log"`
	Token       string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

func GetTasks(db *gorm.DB, order string, page uint, perPage uint, token string) (tasks []Task, total uint, err error) {
	noLog := "id, youtube_id, title, description, author, url, subtitles, state, state2, reason, file_name, share_link, created_at, updated_at"
	if perPage == 0 {
		perPage = 10
	}
	if token != "" {
		err = db.Where("state <> ?", "Rejected").Where("token = ?", token).
			Order("updated_at " + order).Limit(perPage).Offset((page - 1) * perPage).
			Select(noLog).Find(&tasks).Error
	} else {
		err = db.Where("state <> ?", "Rejected").
			Order("updated_at " + order).Limit(perPage).Offset((page - 1) * perPage).
			Select(noLog).Find(&tasks).Error
	}
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	if token != "" {
		err = db.Model(&Task{}).Where("state <> ?", "Rejected").Where("token = ?", token).Count(&total).Error
	} else {
		err = db.Model(&Task{}).Where("state <> ?", "Rejected").Count(&total).Error
	}
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

func GetRejTasks(db *gorm.DB, order string, page uint, perPage uint, token string) (tasks []Task, total uint, err error) {
	noLog := "id, youtube_id, title, description, author, url, state, reason, file_name, share_link, created_at, updated_at"
	if perPage == 0 {
		perPage = 10
	}
	if token != "" {
		err = db.Where("state = ?", "Rejected").Where("token = ?", token).
			Order("updated_at " + order).Limit(perPage).Offset((page - 1) * perPage).
			Select(noLog).Find(&tasks).Error
	} else {
		err = db.Where("state = ?", "Rejected").
			Order("updated_at " + order).Limit(perPage).Offset((page - 1) * perPage).
			Select(noLog).Find(&tasks).Error
	}
	if err != nil {
		err = errors.Wrap(err, "GetTasks")
		return
	}
	if token != "" {
		err = db.Model(&Task{}).Where("state = ?", "Rejected").Where("token = ?", token).Count(&total).Error
	} else {
		err = db.Model(&Task{}).Where("state = ?", "Rejected").Count(&total).Error
	}
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
