package model

type SPool struct {
	work chan func()
	sem chan struct{}
}

func New(size int) *SPool {
	return &SPool{
		work: make (chan func()),
		sem: make (chan struct{}, size),
	}
}

func (p *SPool) Schedule (task func()) {
	select {
		case p.work <- task:
		case p.sem <- struct{}{}:
			go p.worker(task)
	}
}

func (p *SPool) worker (task func()) {
	defer func () { <-p.sem }()

	for {
		task()
		task = <- p.work
	}
}
