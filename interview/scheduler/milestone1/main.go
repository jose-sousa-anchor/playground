package main

import "fmt"

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
// Milestone 1: Simple FIFO Scheduler
// -----------------------------

// Job represents a job with a name and payload (no priority yet)
type Job struct {
	Name    string
	Payload []string
}

// Scheduler is a simple FIFO job scheduler
type Scheduler struct {
	queue []*Job
	jobs  map[string]*Job // for quick existence checks
}

// NewScheduler creates a new simple FIFO scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		queue: make([]*Job, 0),
		jobs:  make(map[string]*Job),
	}
}

// Schedule adds a job to the scheduler. Returns error if job name already exists.
func (s *Scheduler) Schedule(name string, payload []string) error {
	if _, ok := s.jobs[name]; ok {
		return fmt.Errorf("job %s already exists", name)
	}
	job := &Job{
		Name:    name,
		Payload: payload,
	}
	s.queue = append(s.queue, job)
	s.jobs[name] = job
	return nil
}

// Process processes all scheduled jobs in FIFO order and returns results.
// Each job is removed immediately after processing.
func (s *Scheduler) Process() map[string]map[int][]string {
	results := make(map[string]map[int][]string)

	for len(s.queue) > 0 {
		// Pop first job (FIFO)
		job := s.queue[0]
		s.queue = s.queue[1:]
		delete(s.jobs, job.Name)

		// Process the job
		res := processJob(job.Payload)
		results[job.Name] = res
	}

	return results
}

func main() {
	fmt.Println("=== Milestone 1: Simple FIFO Scheduler ===")
	s := NewScheduler()

	// Schedule some jobs
	_ = s.Schedule("job1", []string{"apple", "grape"})
	_ = s.Schedule("job2", []string{"one", "two", "six"})
	_ = s.Schedule("job3", []string{"ABCDEFGH", "12345678"})
	_ = s.Schedule("job4", []string{}) // empty payload

	fmt.Println("\nProcessing all jobs in FIFO order:")
	results := s.Process()
	for name, res := range results {
		fmt.Printf("Job '%s' result: %v\n", name, res)
	}

	// Try to schedule duplicate
	fmt.Println("\nTrying to add duplicate job:")
	if err := s.Schedule("job1", []string{"test"}); err != nil {
		fmt.Println("Error:", err)
	}
}
