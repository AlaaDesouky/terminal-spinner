package spinner

import (
	"context"
	"io"
	"os"
	"sync"
	"time"
)

type Spinner struct {
	writer io.Writer
	rate time.Duration
	frames []rune
	lock sync.RWMutex

	cancelFunc context.CancelFunc

	doneCh chan struct{}
}

type Config struct {
	Writer io.Writer
	Rate time.Duration
}

func New(c Config) *Spinner {
	s := &Spinner{
		writer: os.Stderr,
		rate: time.Millisecond *250,
		frames: []rune{'-', '\\', '|', '/'},
	}

	if c.Writer != nil {
		s.writer = c.Writer
	}

	if c.Rate != 0 {
		s.rate = c.Rate
	}

	return s
}

func (s *Spinner) Start() {
	if s.isRunning(){
		return 
	}

	s.lock.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	doneCh := make(chan struct{})
	s.doneCh = doneCh
	s.lock.Unlock()

	ticker := time.NewTicker(s.rate)

	go func(){
		defer ticker.Stop()
		for {
			for _, frame := range s.frames {
				s.writer.Write([]byte{byte(frame)})

				select {
				case <- ctx.Done():
					s.writer.Write([]byte("\b"))
					close(doneCh)
					return

				case <-ticker.C:
				}

				s.writer.Write([]byte("\b"))
			}
		}
	}()

}

func (s *Spinner) Stop() {
	if !s.isRunning() {
		return
	}

	s.cancelFunc()
	<-s.doneCh

	s.lock.Lock()
	defer s.lock.Unlock()

	s.doneCh = nil
}

func (s *Spinner) isRunning() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.doneCh != nil
}