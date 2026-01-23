# Job Scheduler Interview Challenge

This repository contains implementations for a multi-stage job scheduler interview challenge, organized by milestone.

## Structure

Each milestone is in its own directory with a standalone implementation that can be run independently:

- **milestone0/** - Warmup: `processJob` function
- **milestone1/** - Simple FIFO scheduler
- **milestone2/** - Priority-based scheduler (HIGH > NORMAL > LOW)
- **milestone3/** - Async processing with 1-second intervals
- **milestone4/** - Complete implementation with `ProcessNow` functionality

## Running Each Milestone

Navigate to any milestone directory and run:

```bash
cd milestone0
go run main.go

cd milestone1
go run main.go

cd milestone2
go run main.go

cd milestone3
go run main.go

cd milestone4
go run main.go
```

## Complete Implementation

The `main.go` file in the root directory contains the complete, production-ready implementation that includes all features from all milestones. This is the reference implementation that was originally developed.

## Milestone Descriptions

### Milestone 0: Warmup - processJob

Implements a function that groups strings by their length.

**Example:**
```
Input: ["apple", "banana", "kiwi"]
Output: {4: ["kiwi"], 5: ["apple"], 6: ["banana"]}
```

### Milestone 1: Simple FIFO Scheduler

A basic job scheduler that:
- Accepts jobs with a name and payload
- Processes jobs in FIFO order
- Returns results for all processed jobs
- Prevents duplicate job names

### Milestone 2: Priority-based Scheduler

Extends Milestone 1 with:
- Three priority levels: LOW, NORMAL, HIGH
- Processes jobs by priority (HIGH first, then NORMAL, then LOW)
- Maintains FIFO order within each priority level

### Milestone 3: Async Processing

Extends Milestone 2 with:
- Asynchronous job processing
- 1-second interval between jobs
- Thread-safe operations using mutexes
- Start/stop controls for the async runner

### Milestone 4: ProcessNow (Skip the Line)

Complete implementation with:
- All features from Milestone 3
- `ProcessNow(jobName)` - immediately processes a specific job
- Channel-based immediate job processing
- Graceful handling of edge cases

## Key Features

- **Thread-safe**: Uses mutexes to protect shared state
- **Extensible**: Easy to add new features or job types
- **Well-tested**: Includes demonstrations and error handling
- **Production-ready**: Handles edge cases, prevents data races

## Implementation Details

- Uses multiple queues (one per priority) for efficient priority-based scheduling
- Job tracking via map for O(1) existence checks
- Channel-based communication for immediate job processing
- WaitGroup for graceful shutdown
