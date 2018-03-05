package handler

import (
	"os/exec"
	"io"
	"os"
	"github.com/yanzay/log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"strconv"
	"github.com/cool2645/youtube-to-netdisk/model"
	"strings"
	"fmt"
	. "github.com/cool2645/youtube-to-netdisk/config"
	"regexp"
	"io/ioutil"
	"github.com/pkg/errors"
)

type Carrier struct {
	task model.Task
	kill chan bool
	log  string
}

var runningCarriers = make(map[int64]Carrier, 0)

func TriggerTask(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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
	var title, description, authorName string
	if len(req.Form["title"]) == 1 {
		title = req.Form["title"][0]
	}
	if len(req.Form["description"]) == 1 {
		description = req.Form["description"][0]
	}
	if len(req.Form["author_name"]) == 1 {
		authorName = req.Form["author_name"][0]
	}

	keywords, err := model.GetKeywords(model.Db)
	if err != nil {
		log.Fatalf("failed while getting keywords from db: %v", err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred retrieving keywords.",
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}
	var hit = make([]string, 0)
	for _, v := range keywords {
		if strings.Contains(title, v.Keyword) {
			hit = append(hit, v.Keyword)
		}
	}
	var task model.Task
	if len(hit) == 0 {
		task, err = model.NewTask(model.Db, title, authorName, description, url, "Rejected", "No keywords hit")
	} else {
		reason := fmt.Sprintf("Keywords %v hit", hit)
		task, err = model.NewTask(model.Db, title, authorName, description, url, "Downloading", reason)
	}
	if err != nil {
		log.Fatalf("failed while writing task to db: %v", err)
		res := map[string]interface{}{
			"code":   http.StatusInternalServerError,
			"result": false,
			"msg":    "Error occurred saving task.",
		}
		responseJson(w, res, http.StatusInternalServerError)
		return
	}

	if task.State == "Downloading" {
		runningCarriers[task.ID] = Carrier{task: task, kill: make(chan bool)}
		go runCarrier(task.ID, runningCarriers[task.ID].kill, task.URL, GlobCfg.ND_FOLDER)
	}

	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   task,
	}
	responseJson(w, res, http.StatusOK)
}

func GetRunningTaskStatus(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var data = make([]model.Task, 0)
	for k, v := range runningCarriers {
		l, err := readLog(k)
		if err != nil {
			log.Error(err)
		}
		v.task.Log = l
		data = append(data, v.task)
	}

	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
		"data":   data,
	}
	responseJson(w, res, http.StatusOK)
}

func readLog(k int64) (log string, err error) {
	fo := GlobCfg.TEMP_PATH + "/" + strconv.FormatInt(k, 10) + "/" + strconv.FormatInt(k, 10) + ".log"
	b, err := ioutil.ReadFile(fo)
	if err != nil {
		err = errors.Wrapf(err, "Fail to read from file %s", fo)
		return
	}
	log += string(b)
	log += "\n"
	fe := GlobCfg.TEMP_PATH + "/" + strconv.FormatInt(k, 10) + "/" + strconv.FormatInt(k, 10) + ".err.log"
	b, err = ioutil.ReadFile(fe)
	if err != nil {
		err = errors.Wrapf(err, "Fail to read from file %s", fe)
		return
	}
	log += string(b)
	return
}

func KillTask(w http.ResponseWriter, req *http.Request, ps httprouter.Params)  {
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
	runningCarriers[taskID].kill <- true
	res := map[string]interface{}{
		"code":   http.StatusOK,
		"result": true,
	}
	responseJson(w, res, http.StatusOK)
}

func runCmd(id int64, kill chan bool, tempPath string, c string, a ...string) (state string) {
	tempPath = tempPath + "/" + strconv.FormatInt(id, 10)
	os.MkdirAll(tempPath, os.ModePerm)
	fo, err := os.Create(tempPath + "/" + strconv.FormatInt(id, 10) + ".log")
	if err != nil {
		log.Error(err)
	}
	fe, err := os.Create(tempPath + "/" + strconv.FormatInt(id, 10) + ".err.log")
	if err != nil {
		log.Error(err)
	}
	cmd := exec.Command(c, a...)
	cmd.Dir = "static"
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go io.Copy(fo, stdout)
	go io.Copy(fe, stderr)

	cmd.Start()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-kill:
		if err := cmd.Process.Kill(); err != nil {
			log.Fatalf("failed to kill: %v", err)
			state = "Running"
		}
		log.Info("process killed")
		state = "Killed"
	case err := <-done:
		if err != nil {
			log.Errorf("process done with error = %v", err)
			state = "Error"
		} else {
			log.Info("process done gracefully without error")
			state = "Finished"
		}
	}
	return
}

func runCarrier(id int64, kill chan bool, url string, ndFolder string) {
	defer delete(runningCarriers, id)
	state := runCmd(id, kill, GlobCfg.TEMP_PATH, GlobCfg.PYTHON_CMD, "-u", "../download.py", url)
	l, err := readLog(id)
	if err != nil {
		log.Error(err)
	}
	r := regexp.MustCompile(`fn:"(.*?)"`)
	p := r.FindStringSubmatch(l)
	var fn string
	if len(p) >= 2 {
		fn = p[1]
	} else {
		return
	}
	if state != "Finished" {
		return
	}
	model.UpdateTaskStatus(model.Db, id, "Uploading", fn, "", l)
	state = runCmd(id, kill, GlobCfg.TEMP_PATH, GlobCfg.PYTHON_CMD, "-u", "../syncBaidu.py", fn, ndFolder)
	l2, err := readLog(id)
	if err != nil {
		log.Error(err)
	}
	r = regexp.MustCompile(`fid:"(.*?)"`)
	p = r.FindStringSubmatch(l2)
	var fid string
	if len(p) >= 2 {
		fid = p[1]
	}
	shareLink := fmt.Sprintf("链接：%s?fid=%s 密码：%s", GlobCfg.ND_SHARELINK, fid, GlobCfg.ND_SHAREPASS)
	model.UpdateTaskStatus(model.Db, id, state, fn, shareLink, l + l2)
	return
}
