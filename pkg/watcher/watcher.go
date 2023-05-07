package watcher

type Watcher struct {
	closeChan  chan struct{}
	closedChan chan struct{}
}

func NewWatcher() *Watcher {
	return &Watcher{
		closeChan:  make(chan struct{}),
		closedChan: make(chan struct{}),
	}
}

func (w *Watcher) CloseChan() chan struct{} {
	return w.closeChan
}

func (w *Watcher) ClosedChan() chan struct{} {
	return w.closedChan
}

func (w *Watcher) Shutdown() {
	close(w.closeChan)
	<-w.closedChan
}
