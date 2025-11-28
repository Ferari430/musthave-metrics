package pkg

type Semaphore struct {
	SemaChan chan struct{}
}

func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{SemaChan: make(chan struct{}, maxReq)}
}

func (s *Semaphore) Acquire() {
	s.SemaChan <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.SemaChan
}
