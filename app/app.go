package app

import (
	"flag"
	"github.com/gorilla/mux"
	"http-fetcher/worker-pool"
	"log"
	"net/http"
	"sync"
)

var (
	tasks      TaskList
	pool       *worker_pool.WorkerPool
	numWorkers = flag.Int("workers", 16, "Num workers")
	useSyncMap = flag.Bool("syncmap", false, "Use sync.Map")
	addr       = flag.String("addr", ":8080", "Bind addr")
)

func Start() {
	flag.Parse()
	pool = worker_pool.New(*numWorkers)

	if *useSyncMap {
		tasks = &sync.Map{}
	} else {
		tasks = NewTaskStore()
	}

	go func() {
		for {
			select {
			case result := <-pool.Results:
				if task, ok := result.(Task); ok {
					tasks.Store(task.Id, task)
				}
			}
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/task", CreateTask).Methods("POST")
	r.HandleFunc("/task", GetTasks).Methods("GET")
	r.HandleFunc("/task/{id}", GetTask).Methods("GET")
	r.HandleFunc("/task/{id}", DeleteTask).Methods("DELETE")

	log.Printf("Bind to addr %v", *addr)
	log.Fatal(http.ListenAndServe(*addr, r))
}
