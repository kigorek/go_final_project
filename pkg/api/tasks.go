package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kigorek/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50) // в параметре максимальное количество записей
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	ids := strings.TrimSpace(r.URL.Query().Get("id"))
	if ids == "" {
		writeJson(w, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		writeJson(w, map[string]string{
			"error": "не верный формат id",
		})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	if task == nil {
		writeJson(w, map[string]string{
			"error": "задача не найдена",
		})
		return
	}

	writeJson(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &task); err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if strings.TrimSpace(task.ID) == "" {
		writeJson(w, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeJson(w, map[string]string{
			"error": "заголовок задачи пустой",
		})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, map[string]string{
		"status": "OK",
	})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJson(w, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	writeJson(w, map[string]string{})
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	ids := strings.TrimSpace(r.URL.Query().Get("id"))
	if ids == "" {
		writeJson(w, map[string]string{
			"error": "не указан идентификатор",
		})
		return
	}
	id, _ := strconv.ParseInt(ids, 10, 64)
	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}
	if task == nil {
		writeJson(w, map[string]string{
			"error": "задача не найдена",
		})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := db.DeleteTask(ids); err != nil {
			writeJson(w, map[string]string{
				"error": err.Error(),
			})
			return
		}
		writeJson(w, map[string]any{})
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := db.UpdateDate(next, ids); err != nil {
		writeJson(w, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, map[string]any{})
}
