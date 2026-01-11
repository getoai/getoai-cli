package util

import (
	"fmt"
	"sync"
	"time"
)

type Spinner struct {
	frames   []string
	interval time.Duration
	message  string
	active   bool
	mu       sync.Mutex
	done     chan bool
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:   []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		interval: 80 * time.Millisecond,
		message:  message,
		done:     make(chan bool),
	}
}

func (s *Spinner) Start() {
	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	go func() {
		i := 0
		for {
			select {
			case <-s.done:
				return
			default:
				s.mu.Lock()
				if !s.active {
					s.mu.Unlock()
					return
				}
				fmt.Printf("\r\033[K%s %s", s.frames[i%len(s.frames)], s.message)
				s.mu.Unlock()
				i++
				time.Sleep(s.interval)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		return
	}
	s.active = false
	s.done <- true
	fmt.Print("\r\033[K")
}

func (s *Spinner) Success(message string) {
	s.Stop()
	fmt.Printf("\033[32m✓\033[0m %s\n", message)
}

func (s *Spinner) Error(message string) {
	s.Stop()
	fmt.Printf("\033[31m✗\033[0m %s\n", message)
}

func (s *Spinner) Info(message string) {
	s.Stop()
	fmt.Printf("\033[34mℹ\033[0m %s\n", message)
}

func (s *Spinner) UpdateMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = message
}

// ProgressBar for showing download/install progress
type ProgressBar struct {
	total   int64
	current int64
	width   int
	message string
	mu      sync.Mutex
}

func NewProgressBar(total int64, message string) *ProgressBar {
	return &ProgressBar{
		total:   total,
		width:   40,
		message: message,
	}
}

func (p *ProgressBar) Update(current int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = current
	p.render()
}

func (p *ProgressBar) Increment(delta int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current += delta
	p.render()
}

func (p *ProgressBar) render() {
	percent := float64(p.current) / float64(p.total)
	if percent > 1 {
		percent = 1
	}
	filled := int(percent * float64(p.width))
	bar := ""
	for i := 0; i < p.width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	fmt.Printf("\r%s [%s] %.1f%%", p.message, bar, percent*100)
}

func (p *ProgressBar) Finish() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = p.total
	p.render()
	fmt.Println()
}
