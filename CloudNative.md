# Cloud native notes from studies

## Why Golang rules the Cloud Native World?
- of course this is not an absolute true, but yet, go has several advantages that make it very suitable for cloud native applications
- comprehensibility: go is more minimalist, avoiding repetition, bookkeeping and boilerplate
- CSP concurrency: the goroutines and channels are a good solution for concurrency programming and encourages developers to limit sharing memory, by instead allowing communication between processes 
- fast builds: most go builds finish in seconds, this helps new instances of application to be deployed within seconds and allows faster development cicles
- language stability: although Go evolves accross time through new versions, the compatibility between them is ensured
- memory safety: Go does not need and neither allows manual memory management and manipulation such as low-level languages (C++ for instance), also Go has its own garbage-collector that carefully tracks and free up memory, this approach also makes the language secure
- performance: although Go is not the most performatic language, it does have an amazing performance, for example, in seconds:

|Algorithm | Go       | C        | Python |
|----------|----------|----------|--------|
|1         | 8.25     | 7.53     | 285.2  |
|2         | 7.48     | 1.96     | 44.14  |
|3         | 1.42     | 1.52     | 78.36  |
- static linking: go applications are compiled directly into native and statically linked executables, having no external dependencies (the trade off is that this can result is slightly larger files)
- static typing: go variables must have a type strictly defined, this allows validation at compile time and allow better maintenance 

## Contexts 
### What is a context and how to use it
- a context is an interface composed of 4 methods:
    1) deadline: returns the time when the context should be cancelled
    2) done: returns a channel that's closed when the context is cancelled
    3) err: returns an error that indicates why the context was cancelled after done channel is closed, if done channel is not closed, returns nil
    4) value: returns the value associated with this context for key or nil if no value (use with care!!!)
- contexts are send request to request, allowing to communicate to nested requests about a cancellation or timeout, for example
- context allows to coordinate a cancellation signal 
- within the context package, there are additional methods to iterate with the context interface, such as:
    1) `WithDeadline(Context, time.Time) (Context, CancelFunc)`: allow to specify the time when the context will be cancelled and the done channel will be closed
    2) `WithTimeout(Context, time.Duration) (Context, CancelFunc)`: allow to specify a duration after which the context will be canceled and the done channel will be closed
    3) `WithCancel(Context) (Context, CancelFunc)`: returns a function that can be called to explicitely cancel the context and close the done channel
    - NOTE: all above methods return a _derived context_ 
    - IMPORTANT: when a context is cancelled, all contexts that derived from it are also canceled from, but contexts that it was derived from are not
        - example: let's say we use `ctx2 := WithCancel(ctx1)`, in this case, `ctx2` derived from `ctx1`, if `ctx1` is cancelled, `ctx2` will also be cancelled (termination signal is sent), but the opposite is not true, if `ctx2` is cancelled, nothing happens to `ctx1`!
    - NOTE: each of the above methods have a derived form such as `WithDeadlineCause` that allows to specify an error.
### Context use cases
- Timeouts/Deadlines: Set a maximum time for an operation to complete, preventing indefinite waits.
- Cancellation Propagation: Propagate cancellation signals from a parent to all child goroutines, for example, when a request is cancelled.
- Request-Scoped Data: Attach metadata (authentication tokens, trace IDs, etc.) to a context that travels through the call chain.
- Goroutine Lifecycle Management: Gracefully terminate goroutines by signalling them via context cancellation.
- Resource Cleanup: Ensure that resources (like open files or network connections) are released when a context is cancelled or exceeds its deadline.
- Chaining Contexts: Combine multiple contexts (timeout, cancellation, values) to support complex workflows in server applications.

## Stability Patterns
### Circuit breakers
- a circuit breaker is a stability pattern that automatically degrades a service in response to a likely fault, preventing larger or cascading failures by eliminating recurring errors and providing reasonable error response
- most logics apply an automatic close after some time, 
- most logics apply a backoff, meaning that the rate of retries is reduced over time (increase close period between attempts)
- a circuit breaker can assume 3 different states:
    1) closed: when requests are successful (sends the request and receives back)
    2) open: when a number of failures (consecutive or not depending on trigger) happens  (does not even send the request)
    3) half-open: after an waiting interval (sends the request and waits for receiving it back, if received, switch to closed state, if not switch back to open state)

### Retry
- a retry pattern aims to overcome a transient fault in a distributed system by retrying the request/operation
- transient faults often are resolved after a bit of time, in a way that retrying the request after a reasonable delay is likely to be successful
- retry must be only used with idempotent operations (those that produce the same effect when called multiple times)
- most retry patterns implements a backoff algorithm - which increments the waiting interval between failures (usually, a jitter is added to the backoff - a random waiting, in a way to prevent all requests to be retried at the same exact timing and cause warm/high workload to the system)

### Debounce 
- a debouncer is a stability pattern that automatically limits the frequency of a function invocation so that only the first or last of all calls is actually performed
- a serie of similar calls that are tightly clustered in time are restricted to only one call
- implementation is simple: on each call of the outer function a time interval is set, any subsequent calls before the time interval expires, are ignored

### Throttle
- a throttle limits the frequency of a function call to a maximum number of invocations per unit of time
- the most common application is to prevent spikes in the number of requests that could saturate the system 
- debounce and throttle are similar and can both help with the same purpose
- the most used pattern for throttle is the token bucket: 
    - when a request is made, a token is taken from the bucket
    - while there is no token in the bucket, no additional request can be made
    - after some interval, a new token is automatically added to the bucket
    - then, a new request can use this bucket 

#### What is the response from Throttle and Debounce during those intervals that limits the requests?
- Mainly, there are 3 most used patterns:
    1) Return a `http.StatusTooManyRequests`
    2) Replay the last response (usually from a cache)
    3) Enqueue the execution for when there are sufficient tokens available

### Timeout 
- a timeout allows a process to stop waiting for an answer once it is clear that an answer may not be coming
- given a service request or a function call that is running for longer than an expected time, the timeout tells the request to simple stop waiting
- in go, the simpler way of doing so is by using `context`, specifically the `context.WithTimeout(ctx, timeout)`
- most functions from standard libraries that may take a long time to return already have implemented a derived function with context, for example: `http.NewRequestWithContext`
- for the functions that does not already have it, the usual is to wrap the function within your context, another way if to call the function in a goroutine and monitor the time it takes for it to receive a response

## Concurrency Patterns

### FanIn and FanOut
- a fan in pattern multiplexes multiple input channels onto one output channel
    - a new goroutine is triggered for each source channel
    - processing is done in each goroutine
    - all results are sent back through the same channel (single output stream)
- a fan out pattern evenly distributes messages from an input channel to multiple output channels
    - a new goroutine is triggered for each task
    - processing is done in each goroutine
    - each result is sent in a separated channel
- mostly fan in and fan out are used together: first the work is distributed accross multiple workers (fan out) then the results are collected in a single stream (fan in)
- some examples: web scrapping, image/video processing (resize, filtering, video encoding, etc), microservices requests, log processing, batch data processing, etc

### Future // Promise // Delay
- Provides a placeholder for a value that is not known yet (since it is still being generated by an asynchronous process)
- Not largely used in Go, because channels often do the work
- Assume the application is a serie of tasks such as taskA, taskB, taskC, taskD and that tasks B and C does not depend on resulting from A, but D does. In this scenario, user could trigger taskA wrapped into a Future (to start it asynchronously), proceed with taskB and taskC concurrently, and then take the result immediately when needed to use with taskD by using `taskA.future()`, which means that if taskA is completed, the result will be made available right away, and if not it will block waiting for it to be available, to only them proceeding with taskD

### Sharding (vertical sharding) or Concurrent Maps
- (vertical) sharding splits a large data structure into multiple partitions to localize the effects of read/write locks
    - (horizontal) sharing on other hand, partitions data accross multiple instances, providing data redundancy and load balancing, in exchange of latency and complex distributed data handling
- Assume a data structure that is often accessed by concurrent services, causing multiple locks and a bottleneck (more time waiting due to locks than processing), this will cause a problem called _lock contention_. To reduce this problem, an alternative is to split the data structure in multiple, using this strategy only a portion of the structure will be locked at each time, decreasing the overall lock contention time.
- This is achieved by means of a ShardedMap - a map of maps (or something else that will be accessed and locked)
- **When to use? Examples...** 
    1) Local Caching Systems (instead of Redis/Memcached)
    2) Load Balancers / Rate Limiters: e.g. concurrent processes reading/writting to a map where worker IP is the key and value is the number of ongoing tasks per worker (load balancing / rate limiter)
    3) Counting events per user / session : e.g. concurrent processes reading/writting to a map where user/session is the key and value is the count of it

### WorkerPool or ThreadPool 
- a worker pool directly multiple processes to concurrently execute work on a collection of input 
- used to run concurrent tasks but limited to a number of workers (goroutines) 
- each worker is designed by a goroutine, the jobs are input to workers via jobs channel and results are shared by workers via results channel (both channels shared among all workers)

### Chord
- chord (or Join) is a pattern that performs an atomic consumption of messages from each member of a group of channels
- the chord receives inputs from a group of channels and only emits a result after all input channels have received some value, the result is an aggregation of all channels
- different behaviors can be configured when the same channel receives multiple inputs before the other channels, one possibility is to overwrite the already existing value, other possibility is to discard the additional value
- chord pattern is similar to fan in pattern, both consume data from multiple channels and aggregate, the difference is that fan in immediately forwards what it receive in any of input channels, while chord will wait until all channels have received something (thus an atomic consumption of messages)


