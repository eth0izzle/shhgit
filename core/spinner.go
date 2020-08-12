package core

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const SpinnerChars = `|/-\`

type Spinner struct {
	mu     sync.Mutex
	frames []rune
	length int
	pos    int
}

func New() *Spinner {
	s := &Spinner{}
	s.frames = []rune(SpinnerChars)
	s.length = len(s.frames)

	return s
}

func (s *Spinner) Next() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	r := s.frames[s.pos%s.length]
	s.pos++
	return string(r)
}

func ShowSpinner() func() {
	spinner := New()
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
			default:
				fmt.Fprintf(os.Stderr, "\r%s", spinner.Next())
				time.Sleep(150 * time.Millisecond)
			}
		}
	}()

	return func() {
		done <- true
		fmt.Fprintf(os.Stderr, "\033[%dD", 1)
	}
}
