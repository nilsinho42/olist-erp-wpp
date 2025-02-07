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

## Circuit breakers
- a circuit breaker is a stability pattern that automatically degrades a service in response to a likely fault, preventing larger or cascading failures by eliminating recurring errors and providing reasonable error response
- most logics apply an automatic close after some time, 
- most logics apply a backoff, meaning that the rate of retries is reduced over time (increase close period between attempts)

## Debounce 
- a debouncer is a stability pattern that automatically limits the frequency of a function invocation so that only the first or last of all calls is actually performed
- a serie of similar calls taht are tightly clustered in time are restricted to only one call
- implementation is simple: on each call of the outer function a time interval is set, any subsequent calls before the time interval expires, are ignored