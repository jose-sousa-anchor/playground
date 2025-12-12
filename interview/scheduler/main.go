package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// -----------------------------
// Warmup: processJob
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
// Scheduler types & priorities
// -----------------------------
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

type Job struct {
	Name     string
	Payload  []string
	Priority Priority
}

// -----------------------------
// Scheduler
// -----------------------------
type Scheduler struct {
	mu sync.Mutex
	// three queues for priorities (front = index 0)
	highQueue   []*Job
	normalQueue []*Job
	lowQueue    []*Job

	// map for quick existence checks (optional)
	jobs map[string]*Job

	// async runner controls
	immediateCh chan *Job // push here to be processed next (ProcessNow)
	stopCh      chan struct{}
	running     bool
	wg          sync.WaitGroup
}

// NewScheduler creates the scheduler
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

// Schedule inserts a new job. Returns error if job name already exists.
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

// helper to pop next job according to priority order. Returns nil if no jobs.
func (s *Scheduler) popNextJobLocked() *Job {
	// MUST be called with s.mu held
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

// ProcessAll synchronously processes all scheduled jobs and returns map[name] -> result.
// It drains current queues at the time of call.
func (s *Scheduler) ProcessAll() (map[string]map[int][]string, error) {
	results := make(map[string]map[int][]string)

	for {
		s.mu.Lock()
		job := s.popNextJobLocked()
		s.mu.Unlock()

		if job == nil {
			break
		}
		res := processJob(job.Payload)
		results[job.Name] = res
	}

	return results, nil
}

// StartAsync starts the asynchronous runner that processes one job per second and prints the result.
// It's safe to call StartAsync multiple times: only first call starts the runner.
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
				// process immediate job now (highest priority)
				printJobProcessing(immediateJob)
			case <-ticker.C:
				// pop next job by priority
				s.mu.Lock()
				job := s.popNextJobLocked()
				s.mu.Unlock()
				if job != nil {
					printJobProcessing(job)
				}
				// else: nothing to do this tick
			}
		}
	}()
}

// StopAsync stops the async runner (blocks until runner goroutine exits)
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

	// recreate channels so StartAsync can be called again if desired
	s.immediateCh = make(chan *Job, 10)
	s.stopCh = make(chan struct{})
}

// ProcessNow finds a scheduled job by name, removes it from its queue and schedules it to be processed immediately.
// Returns error if job not found.
func (s *Scheduler) ProcessNow(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("job %s not found", name)
	}

	// Remove job from its queue (linear scan). After removal, delete from jobs map.
	removed := false
	switch job.Priority {
	case HIGH:
		for i := 0; i < len(s.highQueue); i++ {
			if s.highQueue[i].Name == name {
				// remove element i
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
	// remove from map so scheduler doesn't double-process it
	delete(s.jobs, name)

	// push to immediate channel (non-blocking if buffer available)
	select {
	case s.immediateCh <- job:
		// scheduled to run immediately (or next select if currently processing)
		return nil
	default:
		// immediateCh full - processed synchronously in caller to avoid losing job
		s.mu.Unlock()
		// process outside lock
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

// -----------------------------
// Example usage / tests in main()
// -----------------------------
func main() {
	fmt.Println("=== Warmup: processJob ===")
	ex1 := []string{"apple", "banana", "kiwi", "grape", "fig", "pear", "peach"}
	fmt.Println("Input:", ex1)
	fmt.Println("Grouped by length:", processJob(ex1))

	fmt.Println("\n=== Scheduler demo ===")
	s := NewScheduler()
	// Schedule three jobs
	_ = s.Schedule("job1", []string{"apple", "grape"}, NORMAL)
	_ = s.Schedule("job2", []string{"one", "two", "six"}, HIGH)
	_ = s.Schedule("job3", []string{"ABCDEFGH", "12345678"}, LOW)
	_ = s.Schedule("job4", []string{}, NORMAL) // empty payload handled

	// ProcessAll synchronously
	fmt.Println("\n-- ProcessAll (synchronous) --")
	results, _ := s.ProcessAll()
	for name, res := range results {
		fmt.Printf("Sync result: %s -> %v\n", name, res)
	}

	// Schedule more jobs for async demo
	_ = s.Schedule("async1", []string{"a", "bb", "ccc"}, NORMAL)
	_ = s.Schedule("async2", []string{"alpha", "beta"}, HIGH)
	_ = s.Schedule("async3", []string{"x", "y"}, LOW)

	// Start async runner
	fmt.Println("\n-- Start async runner (1s tick) --")
	s.StartAsync()

	// Wait 500ms and then request processNow for async3
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Requesting ProcessNow(async3)")
	if err := s.ProcessNow("async3"); err != nil {
		fmt.Println("ProcessNow error:", err)
	}

	// Also schedule another high-priority job while running
	_ = s.Schedule("urgent", []string{"urgent1", "urgent22"}, HIGH)

	// Let runner run for a few seconds
	time.Sleep(5 * time.Second)

	// stop runner and exit
	fmt.Println("\nStopping async runner.")
	s.StopAsync()
}
