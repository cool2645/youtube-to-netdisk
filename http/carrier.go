package http

import (
	"github.com/cool2645/youtube-to-netdisk/carrier"
	"github.com/yanzay/log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"github.com/cool2645/youtube-to-netdisk/model"
	"strings"
)

func triggerTask(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	log.Debug(req.Form)
	if len(req.Form["url"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing url.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	url := req.Form["url"][0]
	url = strings.TrimSpace(url)
	var title, description, authorName string
	if len(req.Form["title"]) == 1 {
		title = req.Form["title"][0]
		title = strings.TrimSpace(title)
	}
	if len(req.Form["description"]) == 1 {
		description = req.Form["description"][0]
		description = strings.TrimSpace(description)
	}
	if len(req.Form["author_name"]) == 1 {
		authorName = req.Form["author_name"][0]
		authorName = strings.TrimSpace(authorName)
	}
	task := model.Task{
		Title: title,
		Author: authorName,
		Description: description,
		URL: url,
	}
	_, err := carrier.Push(&task)
	if err != nil {
		log.Errorf("failed while pushing task: %v", err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred saving task.",
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   task,
	}
	responseJson(w, res, http.StatusOK)
}

func retryTask(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	taskID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing task id.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	task, err := model.GetTask(model.Db, taskID)
	if err != nil {
		if err.Error() == "GetTask: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred retrieving task: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		log.Fatalf("failed while reading task from db: %v", err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred retrieving task: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	if task.State != "Exception" && task.State != "Finished" && task.State != "Canceled" {
		res := map[string]interface{}{
			"code":   http.StatusOK,
			"result": false,
			"msg":    "Running task cannot be retried, cancel it first.",
		}
		responseJson(w, res, http.StatusNotFound)
		return
	}
	_, err = carrier.Push(&task)
	if err != nil {
		log.Errorf("failed while pushing task: %v", err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred saving task.",
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   task,
	}
	responseJson(w, res, http.StatusOK)
}

func cancelTask(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	taskID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Error occurred parsing task id.",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	task, err := model.GetTask(model.Db, taskID)
	if err != nil {
		log.Error(err)
		if err.Error() == "GetTask: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred querying tasks: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred querying tasks: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	err = carrier.Cancel(task)
	if err != nil {
		log.Errorf("failed while pushing task: %v", err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred saving task.",
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
	}
	responseJson(w, res, http.StatusOK)
}
