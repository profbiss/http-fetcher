package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func CreateTask(writer http.ResponseWriter, request *http.Request) {
	var task Task

	dec := json.NewDecoder(request.Body)
	err := dec.Decode(&task.job)
	if err != nil {
		log.Printf("Decode request body error %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	task.Id = uuid.New()

	tasks.Store(task.Id, task)
	pool.Tasks <- task

	b, err := json.Marshal(task)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
	if _, err := writer.Write(b); err != nil {
		log.Printf("Write response error %v", err)
	}

}

func GetTasks(writer http.ResponseWriter, request *http.Request) {
	var resp []Task

	tasks.Range(func(key, value interface{}) bool {

		resp = append(resp, value.(Task))

		return true
	})

	b, err := json.Marshal(resp)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(b)
	if err != nil {
		log.Printf("Write response error %v", err)
	}
}

func GetTask(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ouuid, err := uuid.Parse(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	task, ok := tasks.Load(ouuid)
	if !ok {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(task)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(b)
	if err != nil {
		log.Printf("Write response error %v", err)
	}

}

func DeleteTask(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, ok := vars["id"]
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ouuid, err := uuid.Parse(id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks.Delete(ouuid)
	writer.WriteHeader(http.StatusOK)
}
