package utils

import "sync"

type Semaphore struct {
	maxNum  int
	Threads chan int
	Wg      sync.WaitGroup
}

func NewSemaphore (maxNum int) *Semaphore {
	var sem = new(Semaphore)
	sem.maxNum  = maxNum
	sem.Threads = make(chan int, maxNum)
	return sem
}

func (sem *Semaphore) P() {
	sem.Threads <- 1
	sem.Wg.Add(1)
}

func (sem *Semaphore) V() {
	sem.Wg.Done()
	<-sem.Threads
}

func (sem *Semaphore) Wait() {
	sem.Wg.Wait()
}

func (sem *Semaphore) Discard() {
	sem.Threads = make(chan int, sem.maxNum)
}


