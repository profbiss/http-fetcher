package app

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Task struct {
	Id     uuid.UUID   `json:"id"`
	Result *TaskResult `json:"result,omitempty"`
	job    TaskJob
}

type TaskJob struct {
	Method  string              `json:"method"`
	Url     string              `json:"url"`
	Body    string              `json:"body"`
	Headers map[string][]string `json:"headers,omitempty"`
}

type TaskResult struct {
	InternalError string              `json:"internal_error,omitempty"`
	Status        int                 `json:"status,omitempty"`
	Headers       map[string][]string `json:"headers,omitempty"`
	Length        int64               `json:"length,omitempty"`
}

func (t Task) Run() interface{} {
	t.Result = &TaskResult{}

	r, err := http.NewRequest(t.job.Method, t.job.Url, strings.NewReader(t.job.Body))
	if err != nil {
		t.Result.InternalError = fmt.Sprintf("Request create error %v", err)
		log.Print(t.Result.InternalError)
		return t
	}

	r.Header = t.job.Headers

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		t.Result.InternalError = fmt.Sprintf("Request execute error %v", err)
		log.Print(t.Result.InternalError)
		return t
	}

	t.Result.Status = resp.StatusCode
	t.Result.Length = resp.ContentLength
	t.Result.Headers = resp.Header

	return t
}

type TaskList interface {
	Load(key interface{}) (value interface{}, ok bool)
	Store(key, value interface{})
	Delete(key interface{})
	Range(f func(key, value interface{}) bool)
}

type TaskStore struct {
	mx    sync.RWMutex
	tasks map[interface{}]interface{}
}

func NewTaskStore() *TaskStore {
	store := TaskStore{
		tasks: make(map[interface{}]interface{}),
	}
	return &store
}

func (s *TaskStore) Load(key interface{}) (value interface{}, ok bool) {
	s.mx.RLock()
	value, ok = s.tasks[key]
	s.mx.RUnlock()
	return
}

func (s *TaskStore) Store(key, value interface{}) {
	s.mx.Lock()
	s.tasks[key] = value
	s.mx.Unlock()
}

func (s *TaskStore) Delete(key interface{}) {
	s.mx.Lock()
	delete(s.tasks, key)
	s.mx.Unlock()
}

func (s *TaskStore) Range(f func(key, value interface{}) bool) {
	for k, v := range s.tasks {
		if !f(k, v) {
			break
		}
	}
}
