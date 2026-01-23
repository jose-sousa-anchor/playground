package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// -----------------------------
// Milestone 0: processJob function
// -----------------------------
func processJob(items []string) map[int][]string {
	result := make(map[int][]string)
	for _, s := range items {
		l := len(s)
		result[l] = append(result[l], s)
	}
	return result
}

// -----------------------------
// Milestone 2: Priority types
// -----------------------------

// Priority levels for jobs
type Priority int

const (
	LOW Priority = iota
	NORMAL
	HIGH
)

func (p Priority) String() string {
	switch p {
	case HIGH:
		return "HIGH"
	case NORMAL:
		return "NORMAL"
	default:
		return "LOW"
	}
}

// Job represents a job with name, payload, and priority
type Job struct {
	Name     string
	Payload  []string
	Priority Priority
}

// -----------------------------
// Milestone 4: ProcessNow (Skip the Line)
// -----------------------------

// Scheduler is the complete scheduler with ProcessNow functionality
type Scheduler struct {
	mu          sync.Mutex
	highQueue   []*Job
	normalQueue []*Job
	lowQueue    []*Job
	jobs        map[string]*Job

	// async runner controls
	immediateCh chan *Job // push here to be processed next (ProcessNow)
	stopCh      chan struct{}
	running     bool
	wg          sync.WaitGroup
}

// NewScheduler creates a new scheduler with ProcessNow capability
func NewScheduler() *Scheduler {
	return &Scheduler{
		highQueue:   make([]*Job, 0),
		normalQueue: make([]*Job, 0),
		lowQueue:    make([]*Job, 0),
		jobs:        make(map[string]*Job),
		immediateCh: make(chan *Job, 10),
		stopCh:      make(chan struct{}),
	}
}

// Schedule adds a job with a priority. Returns error if job name already exists.
func (s *Scheduler) Schedule(name string, payload []string, pr Priority) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.jobs[name]; ok {
		return fmt.Errorf("job %s already exists", name)
	}
	job := &Job{
		Name:     name,
		Payload:  payload,
		Priority: pr,
	}

	// Add to appropriate priority queue
	switch pr {
	case HIGH:
		s.highQueue = append(s.highQueue, job)
	case NORMAL:
		s.normalQueue = append(s.normalQueue, job)
	default:
		s.lowQueue = append(s.lowQueue, job)
	}
	s.jobs[name] = job
	return nil
}

// popNextLocked pops the next job according to priority order.
// MUST be called with s.mu held.
func (s *Scheduler) popNextLocked() *Job {
	if len(s.highQueue) > 0 {
		job := s.highQueue[0]
		s.highQueue = s.highQueue[1:]
		delete(s.jobs, job.Name)
		return job
	}
	if len(s.normalQueue) > 0 {
		job := s.normalQueue[0]
		s.normalQueue = s.normalQueue[1:]
		delete(s.jobs, job.Name)
		return job
	}
	if len(s.lowQueue) > 0 {
		job := s.lowQueue[0]
		s.lowQueue = s.lowQueue[1:]
		delete(s.jobs, job.Name)
		return job
	}
	return nil
}

// StartAsync starts the asynchronous runner that processes one job per second and prints the result.
func (s *Scheduler) StartAsync() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-s.stopCh:
				return
			case immediateJob := <-s.immediateCh:
				// Process immediate job now (highest priority)
				printJobProcessing(immediateJob)
			case <-ticker.C:
				// Pop next job by priority
				s.mu.Lock()
				job := s.popNextLocked()
				s.mu.Unlock()

				if job != nil {
					printJobProcessing(job)
				}
			}
		}
	}()
}

// StopAsync stops the async runner and waits for it to exit
func (s *Scheduler) StopAsync() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.stopCh)
	s.wg.Wait()

	// Recreate channels so StartAsync can be called again if needed
	s.immediateCh = make(chan *Job, 10)
	s.stopCh = make(chan struct{})
}

// ProcessNow finds a scheduled job by name, removes it from its queue,
// and schedules it to be processed immediately.
// Returns error if job not found.
func (s *Scheduler) ProcessNow(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("job %s not found", name)
	}

	// Remove job from its queue (linear scan)
	removed := false
	switch job.Priority {
	case HIGH:
		for i := 0; i < len(s.highQueue); i++ {
			if s.highQueue[i].Name == name {
				s.highQueue = append(s.highQueue[:i], s.highQueue[i+1:]...)
				removed = true
				break
			}
		}
	case NORMAL:
		for i := 0; i < len(s.normalQueue); i++ {
			if s.normalQueue[i].Name == name {
				s.normalQueue = append(s.normalQueue[:i], s.normalQueue[i+1:]...)
				removed = true
				break
			}
		}
	case LOW:
		for i := 0; i < len(s.lowQueue); i++ {
			if s.lowQueue[i].Name == name {
				s.lowQueue = append(s.lowQueue[:i], s.lowQueue[i+1:]...)
				removed = true
				break
			}
		}
	}
	if !removed {
		return errors.New("inconsistent state: job found in map but not in queue")
	}

	// Remove from map so scheduler doesn't double-process it
	delete(s.jobs, name)

	// Push to immediate channel (non-blocking if buffer available)
	select {
	case s.immediateCh <- job:
		// Scheduled to run immediately
		return nil
	default:
		// immediateCh full - process synchronously to avoid losing job
		s.mu.Unlock()
		printJobProcessing(job)
		s.mu.Lock()
		return nil
	}
}

// helper to print job processing
func printJobProcessing(job *Job) {
	fmt.Printf("Processing job '%s' (priority=%s) ...\n", job.Name, job.Priority.String())
	res := processJob(job.Payload)
	fmt.Printf("Result for job '%s': %v\n", job.Name, res)
}

func main() {
	fmt.Println("=== Milestone 4: ProcessNow (Skip the Line) ===")
	s := NewScheduler()

	// Schedule jobs
	_ = s.Schedule("async1", []string{"a", "bb", "ccc"}, NORMAL)
	_ = s.Schedule("async2", []string{"alpha", "beta"}, HIGH)
	_ = s.Schedule("async3", []string{"x", "y"}, LOW)

	// Start async runner
	fmt.Println("\nStarting async runner...")
	s.StartAsync()

	// Wait 500ms and then request ProcessNow for async3 (LOW priority job)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("\n>>> Requesting ProcessNow(async3) - skipping the line! <<<")
	if err := s.ProcessNow("async3"); err != nil {
		fmt.Println("ProcessNow error:", err)
	}

	// Schedule another high-priority job while running
	time.Sleep(1 * time.Second)
	_ = s.Schedule("urgent", []string{"urgent1", "urgent22"}, HIGH)

	// Let runner run for a few more seconds
	time.Sleep(4 * time.Second)

	// Stop runner
	fmt.Println("\nStopping async runner.")
	s.StopAsync()

	// Test ProcessNow error case
	fmt.Println("\nTrying ProcessNow on non-existent job:")
	if err := s.ProcessNow("nonexistent"); err != nil {
		fmt.Println("Error:", err)
	}
}
