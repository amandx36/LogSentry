# Worker Pool in Go (Complete Theory)

> These notes are written for revision before interviews and for implementing Worker Pools in LogSentry and future backend projects.

---

# Table of Contents

1. What is a Worker Pool?
2. Why do we need it?
3. Sequential vs Concurrent Processing
4. Worker Pool Architecture
5. Components of a Worker Pool
6. Workflow
7. Channels
8. Goroutines
9. WaitGroup
10. Jobs
11. Workers
12. Results
13. Merge
14. Advantages
15. Disadvantages
16. Worker Pool in LogSentry
17. Interview Questions
18. Complexity
19. Production Notes

---

# 1. What is a Worker Pool?

A Worker Pool is a concurrency design pattern in Go where a fixed number of workers
(goroutines) continuously receive tasks from a queue (channel) and process them in parallel.

Instead of creating one goroutine per task, we create a limited number of workers.

Those workers continuously execute incoming jobs.

Think of it as:

```
Tasks
   │
   ▼
Queue
   │
   ▼
Workers
```

---

# Real Life Example

Imagine a restaurant.

Customers

↓

Waiter

↓

Kitchen

↓

Chef 1

Chef 2

Chef 3

Chef 4

Customers don't cook.

The waiter assigns orders to available chefs.

The chefs work simultaneously.

Exactly same happens inside a Worker Pool.

---

# 2. Why do we need Worker Pools?

Suppose we have

```
1000 log files
```

Sequential approach

```
File1

↓

Parse

↓

File2

↓

Parse

↓

File3

↓

Parse
```

Only one CPU core is working.

---

Using Worker Pool

```
Worker1 -> File1

Worker2 -> File2

Worker3 -> File3

Worker4 -> File4
```

Now all CPU cores work together.

---

Benefits

- Faster execution
- Better CPU utilization
- Less idle time
- Controlled concurrency
- Production scalable

---

# 3. Sequential vs Concurrent

Sequential

```
Task1

↓

Task2

↓

Task3

↓

Task4
```

Time

```
1
2
3
4
```

---

Concurrent

```
Worker1 -> Task1

Worker2 -> Task2

Worker3 -> Task3

Worker4 -> Task4
```

Time

```
1
```

All execute simultaneously.

---

# 4. Worker Pool Architecture

```
Jobs

↓

Jobs Channel

↓

Worker Pool

↓

Results Channel

↓

Main Goroutine
```

---

Visual

```
               Jobs

                 │

                 ▼

          Jobs Channel

                 │

     ┌───────────┼───────────┐

     ▼           ▼           ▼

 Worker1     Worker2     Worker3

     │           │           │

     └───────────┼───────────┘

                 ▼

          Results Channel

                 │

                 ▼

          Main Goroutine
```

---

# 5. Components

A Worker Pool consists of

- Jobs
- Workers
- Goroutines
- Channels
- Results
- WaitGroup

---

# 6. Jobs

A Job represents one unit of work.

Example

```
Parse this file

↓

server1.log
```

Job Structure

```go
type Job struct {
    FilePath string
}
```

Every worker receives one Job.

---

# 7. Workers

Workers are goroutines.

Their only responsibility is

Receive Job

↓

Process Job

↓

Return Result

Workers do not know about each other.

---

# 8. Goroutines

A goroutine is a lightweight thread managed by Go.

Instead of

```
Parse File

↓

Wait

↓

Next File
```

We create

```
go Parse(File1)

go Parse(File2)

go Parse(File3)
```

Now everything runs concurrently.

---

# 9. Channels

Channels are communication pipes.

Workers communicate using channels.

Example

```
Job

↓

Channel

↓

Worker
```

Workers never call each other directly.

---

Types

Unbuffered

```go
jobs := make(chan Job)
```

Buffered

```go
jobs := make(chan Job,100)
```

---

# 10. WaitGroup

Problem

Main function exits.

Workers are still running.

Solution

WaitGroup

Workflow

```
Add Workers

↓

Workers Start

↓

Done()

↓

Wait()

↓

Program Ends
```

Without WaitGroup

Workers may terminate unexpectedly.

---

# 11. Results

Workers don't print.

Workers return results.

Example

```go
type Result struct {

    Report LogReport

    Error error

}
```

Results are sent back through another channel.

---

# 12. Merge

Each worker parses one file.

Worker1

↓

LogReport1

Worker2

↓

LogReport2

Worker3

↓

LogReport3

Main Goroutine merges them.

Final

```
LogReport
```

---

# 13. Why Merge?

Suppose

```
server1.log

↓

20 Errors

server2.log

↓

15 Errors
```

Final Dashboard should show

```
35 Errors
```

Therefore reports are merged.

---

# 14. Worker Lifecycle

Worker starts

↓

Wait for Job

↓

Receive Job

↓

Open File

↓

Parse File

↓

Create LogReport

↓

Send Result

↓

Wait for another Job

Workers never stop until channel closes.

---

# 15. Complete Workflow

```
Directory

↓

ReadDir()

↓

Jobs Channel

↓

Worker Pool

↓

ParseSingleFile()

↓

Result Channel

↓

MergeReports()

↓

Final LogReport

↓

Storage

↓

Analytics

↓

Search

↓

Dashboard
```

---

# 16. Worker Pool in LogSentry

Current

```
Directory

↓

for files

↓

ParseSingleFile()

↓

LogReport
```

Future

```
Directory Scanner

↓

Jobs

↓

Worker Pool

↓

Parser

↓

Merge

↓

Final Report

↓

Writer

↓

Database

↓

Analytics

↓

Dashboard
```

---

# 17. Advantages

 Faster

 Better CPU utilization

 Scalable

 Easy to extend

 Production ready


 Industry standard

---

# 18. Disadvantages

- More complex
- Harder debugging
- Synchronization required
- Race conditions possible
- Shared memory problems

---

# 19. Common Mistakes

Creating one goroutine for every file.

Bad

```
100000 files

↓

100000 goroutines
```

Good

```
100000 files

↓

8 workers
```

Workers continuously process files.

---

# 20. Worker Count

Never create unlimited workers.

Usually

```
runtime.NumCPU()
```

or

```
4

8

16
```

depending upon workload.

---

# 21. Interview Questions

### Q1 What is a Worker Pool?

A concurrency pattern where a fixed number of workers continuously process tasks from a shared queue.

---

### Q2 Why not create one goroutine per task?

Because thousands of goroutines increase scheduling overhead and memory usage.

Worker Pools control concurrency.

---

### Q3 What is the role of Channels?

Channels transfer jobs and results safely between goroutines.

---

### Q4 Why WaitGroup?

To wait until all workers finish before exiting the program.

---

### Q5 Why use Worker Pools?

- Better performance
- Controlled concurrency
- Efficient CPU utilization
- Production scalability

---

### Q6 Difference between Goroutine and Worker?

Goroutine

A lightweight thread.

Worker

A goroutine dedicated to repeatedly processing jobs.

---

### Q7 Worker Pool vs Parallelism

Worker Pool is a concurrency pattern.

Parallelism depends on available CPU cores.

---

# 22. Complexity

Sequential

```
Time

O(n)
```

Worker Pool

```
Time

≈ O(n/workers)
```

(actual speed depends on CPU cores, disk I/O, and workload)

---

# 23. Production Architecture (LogSentry)

```
                Log Directory
                      │
                      ▼
              Directory Scanner
                      │
                      ▼
                Worker Pool
                      │
                      ▼
                Parser Engine
                      │
                      ▼
                 Structured Logs
                      │
        ┌─────────────┼─────────────┐
        ▼             ▼             ▼
   Dashboard     Analytics     Storage Layer
                                      │
                  ┌───────────────────┴───────────────────┐
                  ▼                                       ▼
            File Storage                         Database Storage
                                                          │
                                                          ▼
                                                    Search Service
                                                          │
                                                          ▼
                                                       REST API
                                                          │
                                                          ▼
                                                       DevMind
```

---

# Key Takeaways

- Worker Pools limit the number of concurrent goroutines.
- Jobs are sent through channels.
- Workers continuously process jobs.
- Results are returned through another channel.
- WaitGroup waits for all workers.
- Main goroutine merges results.
- Worker Pools are widely used in production systems for file processing, APIs, message queues, and background jobs.