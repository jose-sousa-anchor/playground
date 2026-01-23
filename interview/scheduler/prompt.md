# Description

We’re going to implement a job scheduler and runner by hand, and get to multiple pieces of it.

# Prompt

We’ll be working through a multi stage problem. There’s an infinite amount of steps, you are not expected to finish them all, but we encourage you to think out loud and explain your reasoning as you work through the problems, more than finish all the steps. It’s important to have a well-explained and solid solution. Feel free to make reasonable assumptions, but please state them clearly. Feel free to use the internet, but we do prohibit the use of AI—we’re here to test your skills.

# Progression

## Warm up: Implement the Job Consumer

### Prompt

1. **Write a function named** `process_job` that takes a list of strings and returns a dictionary (or map) where the keys are the unique lengths of the strings in the input list, and the values are lists of strings with that corresponding length. Print the results.

**Example**

```go
Job Input: ["input1", "input12", "input123", "input4", "input1234"]
Expected Output (order of lists doesn't matter):
{6: ["input1", "input4"], 7: ["input12"], 8: ["input123"], 9: ["input1234"]}
```

### Expectations

- They should be able to get through this with few issues.
- They should test for common edge cases (empty, duplicate entries).
- Order of lists doesn't matter.

### Unblocker questions

Should not be expected for this step.

## Milestone 1: Write a simple job scheduler

### Prompt

1. Instead of using a single function, we’re going to create a job management scheduler to schedule and process jobs. This scheduler will manage multiple jobs. Each job will have a name and job contents. 
    
    ### Requirements
    
    - Adding jobs: Implement a `schedule` function to insert a new job with a name and the list of strings to be scheduled.
    - Process jobs: Calling the `process` function on the scheduler should process all the scheduled jobs in a FIFO manner and *return* a list of job results. Every job should be processed individually and once a job has been processed, it should immediately be removed from the scheduler.
    - Print the results.

### Expectations

- For this milestone we’re avoiding on purpose race conditions and language specific locking mechanisms.
- No concurrency for this milestone
- This should take around 20 minutes.
- They should be clear what they use to represent the job. Class, struct, dictionary are all fine.
- They should clearly be able to describe the state machine at this point.
- They should be clear about separating the job definition and processing logic.
- Is the code clear?
- Can they explain their design trade-offs?
- Can they write testable code?
- Should be able to describe runtime and space use.

### Unblocker questions

- How are you storing the jobs?
- What information does each job need to hold?
- How will you know which job to process next?

## Milestone 2: Add Priority

### Prompt

We want to improve the solution by adding priorities: `LOW`, `NORMAL` and `HIGH`. 

Update scheduling to receive a priority and the `process` function to process jobs in order of priority. `HIGH` priority jobs need to be processed first, then `NORMAL`, then `LOW`.

### Gotchas

- They are NOT expected to write a priority queue. If they name it, let them know that they can use the library function.

### Follow-up questions

- Ask why they reached for, and what other options are (multiple lists | single list with sorting | priority queue).
    - What is the runtime?
    - What are the hotspots?

### Expectations

- The task management system should have a way to process values.
- Is the code clear?
- Can they describe the state machines?
- Did they write tests?
- They should be mindful of potential performance implications as they introduce priority.

### Unblocker questions

- How can you efficiently keep track of jobs based on their priority?
- What are some ways to order the jobs before processing?

## Milestone 3: Add async processing of tasks on a 1s delay

### Prompt

Processing of the tasks should not be blocking anymore. We want the scheduler to run asynchronously and process these jobs at a 1s interval. Assume that the `process` function remains a single-threaded operation that returns the processed job result. The asynchronous scheduler should then print the result of each processed job to the terminal.

### Expectations

- Define requirements
- Able to write non-racy code.
- Be able to describe the state machine well, with all the locks or channels.
- Should be able to continue to test the code.

### Unblocker questions

- What are some common ways to handle concurrent access to shared data?
- How can you ensure that the `process` function doesn't interfere with the asynchronous execution?

## Milestone 4: Skip the line endpoint

### Prompt

What if I need to get a task DONE NOW, instead of waiting for the 1s delay.

Add a function call `processNow`.

### Expectations

This is a reach problem (to stretch the goals). If they capture the following, the signs are very good:

- What if they have to process multiple high priority requests?
- Do they stop processing the current request?
- Do they requeue the other requests?
- Can write the logic to handle prioritizing the processing of requests as one.
- Can write a test for this?

## Follow up questions / reach problems:

These can be asked if you think the next step will take too long.

These evaluate for things that can be missed in the technical interview.

1. How do you add deletion?
2. How do you scale this? 
3. What else can this be applied to?
4. How do you get this to a production ready state?
5. What type of monitoring would you add here, and why? 
6. Is there anything else you can optimize?
7. If you had extra time, how would you clean up your implementation?
8. How would you abstract this so it can process many different job types? How would you abstract the dependencies?
9. We’ve also decided that some lists need to be processed together, but come in as separate tasks, how would your structure have to adapt to this? 

# Rubric

### Details

5-Must must hire.

All of the below in 3 and 4, and excellence in the Reach question.

4-Must hire.

- Finished main 2.
- Finished main 1 without much help.
- Code is well formatted.
- TC talks about the state clearly in an organized way.
- Clearly enumerates through edge cases.
- Able to reason about the most important thing, able to disregard less important parts of the code.

3- hire.

- Finished up to main 1 with some help.
- Code is well formatted.
- TC can reason about the state, and can describe state in an organized way for the simpler steps.
- Identifies most edge cases.
- They’re a nice person

2-No hire.

- Needs a lot of help.
- Unclear code

1-Definitely no hire.

- Struggles in the warmups.

|  | Problem Solving | Technical Knowledge/System Design |
| --- | --- | --- |
| Junior | Warmup & M1 |  |
| Mid | M1 & M2 |  |
| Senior | M3 / M4 |  |