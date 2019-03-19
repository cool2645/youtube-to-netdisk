package http

import (
	"github.com/cool2645/youtube-to-netdisk/carrier"
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/julienschmidt/httprouter"
	logging "github.com/yanzay/log"
	"net/http"
	"strconv"
)

func getTasks(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	state := "%"
	order := "asc"
	token := ""
	var page uint = 1
	var perPage uint = 10
	if len(req.Form["state"]) == 1 {
		state = req.Form["state"][0]
	}
	if len(req.Form["order"]) == 1 {
		order = req.Form["order"][0]
	}
	if len(req.Form["token"]) == 1 {
		token = req.Form["token"][0]
	}
	if len(req.Form["page"]) == 1 {
		page64, err := strconv.ParseUint(req.Form["page"][0], 10, 32)
		if err != nil {
			logging.Error(err)
		}
		page = uint(page64)
	}
	if len(req.Form["perPage"]) == 1 {
		perPage64, err := strconv.ParseUint(req.Form["perPage"][0], 10, 32)
		if err != nil {
			logging.Error(err)
		}
		perPage = uint(perPage64)
	}
	var tasks []model.Task
	var total uint
	var err error
	if state == "Rejected" {
		tasks, total, err = model.GetRejTasks(model.Db, order, page, perPage, token)
	} else {
		tasks, total, err = model.GetTasks(model.Db, order, page, perPage, token)
	}
	if err != nil {
		logging.Error(err)
		if err.Error() == "GetTasks: sql: no rows in result set" {
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
	data := map[string]interface{}{
		"total": total,
		"data":  tasks,
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   data,
	}
	responseJson(w, res, http.StatusOK)
}

func getTask(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	taskID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		logging.Error(err)
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
		logging.Error(err)
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
	log, err := carrier.ReadLog(taskID)
	if err != nil {
		logging.Error(err)
	} else {
		task.Log = log
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   task,
	}
	responseJson(w, res, http.StatusOK)
}

func getTaskLog(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	taskID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		logging.Error(err)
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
		logging.Error(err)
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
	log, err := carrier.ReadLog(taskID)
	if err != nil {
		logging.Error(err)
	} else {
		task.Log = log
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   task.Log,
	}
	responseJson(w, res, http.StatusOK)
}
