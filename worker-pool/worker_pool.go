package worker_pool

type Task interface {
	Run() interface{}
}

type WorkerPool struct {
	maxWorkers int
	Tasks      chan Task
	Results    chan interface{}
}

func New(maxWorkers int) *WorkerPool {
	pool := &WorkerPool{
		Tasks:      make(chan Task),
		maxWorkers: maxWorkers,
		Results:    make(chan interface{}),
	}

	pool.dispatch()

	return pool
}

func (p *WorkerPool) dispatch() {
	for i := 0; i < p.maxWorkers; i++ {
		go func() {
			select {
			case task := <-p.Tasks:
				p.Results <- task.Run()
			}
		}()
	}
}
