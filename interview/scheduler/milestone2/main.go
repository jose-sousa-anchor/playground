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
// Milestone 2: Priority-based Scheduler
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

// Scheduler is a priority-based job scheduler
type Scheduler struct {
	highQueue   []*Job
	normalQueue []*Job
	lowQueue    []*Job
	jobs        map[string]*Job // for quick existence checks
}

// NewScheduler creates a new priority-based scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{
		highQueue:   make([]*Job, 0),
		normalQueue: make([]*Job, 0),
		lowQueue:    make([]*Job, 0),
		jobs:        make(map[string]*Job),
	}
}

// Schedule adds a job with a priority. Returns error if job name already exists.
func (s *Scheduler) Schedule(name string, payload []string, pr Priority) error {
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

// popNext pops the next job according to priority order (HIGH > NORMAL > LOW).
// Returns nil if no jobs available.
func (s *Scheduler) popNext() *Job {
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

// Process processes all scheduled jobs in priority order and returns results.
// Each job is removed immediately after processing.
func (s *Scheduler) Process() map[string]map[int][]string {
	results := make(map[string]map[int][]string)

	for {
		job := s.popNext()
		if job == nil {
			break
		}

		// Process the job
		res := processJob(job.Payload)
		results[job.Name] = res
	}

	return results
}

func main() {
	fmt.Println("=== Milestone 2: Priority-based Scheduler ===")
	s := NewScheduler()

	// Schedule jobs with different priorities
	_ = s.Schedule("normal1", []string{"apple", "grape"}, NORMAL)
	_ = s.Schedule("high1", []string{"one", "two", "six"}, HIGH)
	_ = s.Schedule("low1", []string{"ABCDEFGH", "12345678"}, LOW)
	_ = s.Schedule("normal2", []string{"test"}, NORMAL)
	_ = s.Schedule("high2", []string{"urgent", "task"}, HIGH)

	fmt.Println("\nProcessing all jobs by priority (HIGH > NORMAL > LOW):")
	results := s.Process()

	// Print results in order they were processed
	fmt.Println("\nResults (HIGH priority jobs processed first):")
	for name, res := range results {
		fmt.Printf("Job '%s' result: %v\n", name, res)
	}
}
