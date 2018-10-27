package rsync

// tasks pool to be sync
type Worker struct {

	// task pool
	TaskPool chan *Task

	// goroutine num
	maxProcesses int

}

// set consumer config
func (self *Rsync) SetWorkers(maxProcesses int, params []string) {
	self.params = params
	rchan := make(chan *Task, 1000)
	self.Workers.TaskPool = rchan
	self.Workers.maxProcesses = maxProcesses
}

// put rsync task into pool
func (self *Rsync) Collect(r *Task) {
	self.Workers.TaskPool <- r
}

// start some goroutine to consume rsync task
func (self *Rsync) StartWork() {
	if self.recycler != nil {
		go func() {
			self.recycler()
		}()
	}
	for i := 0; i < self.Workers.maxProcesses; i++ {
		go func() {
			for{
				select {
				case rt,ok :=  <- self.Workers.TaskPool: // Task
					if ok {
						self.Sync(rt)
					}
				}
			}
		}()
	}
}