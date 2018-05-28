package handler

import (
	"github.com/cool2645/youtube-to-netdisk/model"
	"github.com/julienschmidt/httprouter"
	"github.com/yanzay/log"
	"net/http"
)

func GetKeywords(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	keywords, err := model.GetKeywords(model.Db)
	if err != nil {
		log.Error(err)
		if err.Error() == "GetTasks: sql: no rows in result set" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred querying keywords: " + err.Error(),
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
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   keywords,
	}
	responseJson(w, res, http.StatusOK)
}

func NewKeyword(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if len(req.Form["keyword"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid keyword",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	keyword := req.Form["keyword"][0]
	node, err := model.SaveKeyword(model.Db, keyword)
	if err != nil {
		log.Error(err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred creating node: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   node,
	}
	responseJson(w, res, http.StatusOK)
}

func DeleteKeyword(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	req.ParseForm()
	if len(req.Form["keyword"]) != 1 {
		res := map[string]interface{}{
			"code":   http.StatusBadRequest,
			"result": false,
			"msg":    "Invalid keyword",
		}
		responseJson(w, res, http.StatusBadRequest)
		return
	}
	keyword := req.Form["keyword"][0]
	err := model.RemoveKeyword(model.Db, keyword)
	if err != nil {
		log.Error(err)
		if err.Error() == "DeleteKeyword: record not found" {
			res := map[string]interface{}{
				"code":   http.StatusNotFound,
				"result": false,
				"msg":    "Error occurred deleting keyword: " + err.Error(),
			}
			responseJson(w, res, http.StatusNotFound)
			return
		}
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred deleting keyword: " + err.Error(),
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"msg":    "success",
	}
	responseJson(w, res, http.StatusOK)
}
