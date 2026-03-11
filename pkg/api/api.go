package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const dateFormat = "20060102"

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", GetTasksHandler)
	http.HandleFunc("/api/task/done", doneHandler)

}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	nowDate := r.FormValue("now")
	dStart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowDate == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowDate)
		if err != nil {
			http.Error(w, "некоректный формат параметра now", http.StatusBadRequest)
			return
		}
	}

	if dStart == "" {
		http.Error(w, "поле Дата не определено", http.StatusBadRequest)
	}

	if repeat == "" {
		http.Error(w, "поле повтор не определено ", http.StatusBadRequest)
	}

	nextDate, err := NextDate(now, dStart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprint(w, nextDate)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		putTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)

	}

}

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
