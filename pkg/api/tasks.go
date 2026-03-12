package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kigorek/go_final_project/pkg/db"
	"github.com/kigorek/go_final_project/tests"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(tests.Limit)
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
			"error": "ошибка получения задачи из БД",
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

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
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
			"error": "неправильгный формат даты или повторения",
		})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{
			"error": "ошибка обновления задачи в БД",
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
			"error": "ошибка удаления залдачи из БД",
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
			"error": "ошибка получения задачи из БД",
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
				"error": "ошибка удаления залдачи из БД",
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
			"error": "неправильгный формат даты или повторения",
		})
		return
	}

	if err := db.UpdateDate(next, ids); err != nil {
		writeJson(w, map[string]string{
			"error": "ошибка при обновлении даты повтора в БД",
		})
		return
	}

	writeJson(w, map[string]any{})
}
