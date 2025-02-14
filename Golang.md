# Golang notes from studies

## features
- Compiled language: converts code into machine language
- Natively Unicode handling allows to process text in all world languages
- 25 keywords (can't be used as variable names): break, case, chan, default, defer, func, go, interface, map, struct, ...
- predeclared names (can be redeclared, but will replace original one): true, false, int, float, make, len, append, ... 

## packages, imports and dependencies
- Code is organized into packages (similar to libraries/modules)
- Each file in the package must begin with same "package ...package_name..." and be located in the same directory
- The main package defines the standalone program 
- Code must import exactly the packages it uses (no more, no less)
- To export a function/variable, it must start with a capital character 

## input values
### command-line arguments
- available through os.Args, where os.Args[0] is the name of the command, and the remaining are the values

## init functions
- `func init() { ... }` defines a initialization function
- any file may contain any number of init functions
- these functions are executed when the program starts, in the order in which they are declared

## scopes 
- there are 3 types of scopes in Go:
    1) Package scope: higher level scope that contains objects accessible throughout all the package
    2) Function scope: intermediate level scope that contains objects accessible only within the function itself
    3) Block scope: lower level scope that contains objects accessible only within the block itself (is, for, switch, etc)
- NOTE: the same object name can be re-used in different scopes

## variable types
- blank identifier (_): wildcard to use whenever the syntax requires a variable but the program logic does not
- types are split accross 4 categories: 1) basic types, 2) aggregate types, 3) reference types and 4) interfaces
- named types: are used to define different use cases for variables of same underlying type, preventing from incorrect usage
- variable conversion: can happen accross different named types of same underlying type, can happen accross numeric types (but int(float_value) will discard the fraction part), and in some additional scenarios
    - NOTE: conversion never fails at run time, it is checked at compile time, if accepted, it will run, but this does not necessairily mean that only because it runs, it produces the expected result, must be careful

### basic types
- integers: 
    - `int8, int16, int32 and int64`: signed int
    - `uint8, uint16, uint32 and uint64`: unsigned int
    - `int` and `uint` are general types which assign either 32 or 64 depending on platform
    - IMPORTANT: if the result of an arithmetic operation has more bits than the type (e.g. result has 40 bits, but type is 32 bits), the operation WILL NOT FAIL, instead it will silently discard the high-order bits - this is known as _overflow_.
- floats:
    - `float32` and `float64` are the only options, there is no generic `float` type, and `float64` is the overall recomendation
- complex:
    - `complex64` and `complex128` are the only options, there is no generic `complex` type
- infinites:
    - are handled through Go math package
- boolean:
    - either `true` or `false`, mostly used in comparisons
    - the unary operator `!` is used for logical negation
    - in Go, combined comparisons follow the _short circuit_ behavior, where if the answer is already defined by the left operand, the right operand is not even evaluated (e.g.: `if left_operand and right_operand`, when left operand is false, it does not even evaluate the right operand)
- strings: an immutable sequence of bytes, meaning every time a string is modified, Go:
    1) creates a new string, 
    2) allocates a new memory block for the new string, and 
    3) copies the old string with modifications to the new block
    - NOTE1: this is a possible costly process.
    - NOTE2: when concatening strings, it might be better to use strings.Join() 
        - 1) checks memory block size required upfront, 2) creates the block with a single memory allocation, 3) copies each modification into the block (about 50% faster in a local test)
    - `s[i]` can be used to take the i-th byte from the string sequence, or even `s[m:n]` can be used to substring
        - NOTE1: it is not possible to ensure that the i-th byte from slice equals to the i-th character from string, as due to encoding (e.g. UTF-8 of non-ASCII characters), some characters may have up to 4 bytes.
            - If all characters are ASCII, then it is possible to assume that i-th byte = i-th character
        - NOTE2: the substring operation `s[m:n]` yields a new string (new memory allocation)
    - there are 2 ways of using strings in Go:
        1) string literal - defined by usage of quotes `"..."`, allows the usage of _escape sequences_
        2) raw string literal - defined by usage of backquotes ``...``, don't allow the usage of _escape sequences_, mostly **applied for regular expressions, HTML templates, JSON, etc**
        - escape sequences: \n (newline), \t (tab), etc

### composite or aggregated types
- there are 4 main composite types: arrays, slices, maps and structs
- arrays: 
    - a **fixed-length** sequence of elements of a particular type
    - are rarely directly used in Go
    - are accessed through index notation `e[i]`
    - the size of an array is part of its type, so `[3]int` is different from `[4]int`!!!
    - some programming languages pass arrays _by reference_ to functions, meaning the function receives a pointer to the array, and any change in the array inside the function will be reflected everywhere, this **is not** how it works in Go, where the arrays are passed to functions as copies (by default), which can led to inneficience
- slice:
    - a **variable-length** sequence of elements of a particular type
    - actually, the slice is a data structure that access a subset of an array, so for every slice there is intimately an array connected
    - can be accessed using s[m:n] (half open intervals - m is included, but n is excluded)
    - made of 3 components: 
        1) pointer: points to first element of the array that is referred by the slice, not necessairily the first index of the array itself, 
        2) length: the number of elements (can not exceed capacity) and
        3) capacity: the number of elements between the first element referred by the slice and the last element in the underlying array
    - multiple slices can share the same array!
    - slices **are not comparable**, only can check if it is nil (empty), but as there are non-nil empty slices, to check if a slice is empty it is recommended to use `len(s) == 0`
    - `make` function is used to create slices of specific length and capacity: `make([]T, len, cap)`
    - how `append` works? `append` is a method of slices, if there is room in the underlying array for grow, append will use it, if not, append will create a new array (with double capacity) and move the entire object to there (possible costly if the array size keeps increasing - possible huge number of memory allocations)
- map: a key/value pair, where key can be any **comparable** type and the value can be any type at all (even interfaces!) 
    - the order of map iteration should be considered random - different in each run
    - map1 is a **reference** for the data structure, meaning when passed as arg to a function, any changes done within the function will be perceived anywhere else
    - `map1 := make(map[string]int)`
    - trying to access a non-existing value in a map does not cause any failure (`res := map1[non_existing]` => `res` will be the zero value for the type, either 0, "", or something else), but it is possible to check if the value exists by the optional return from map (`res, ok := map1[non_existing]` => `res` will be the zero value, but `ok` will be false)
    - maps **are not comparable**, only against nil
    - maps can host nested maps/slices
- structs:
    - is an aggregated type that groups together named values of arbitrary types (similar to python classes)
    - fields are accessed using dot notation (`struct.field`)
    - a field is only exported if begins with capital letter
    - allow composition: use another struct in a struct
    - don't have inheritance!
    - can be used to unmarshal JSON by using `field tags` such as `type struct1 struct {field1 string json:"json_column_name"}`

### interfaces
- interfaces types express generalizations or abstractions about the behaviors of other types, allowing to write more flexible functions 
- in Golang, interfaces are _satisfied implicitely_, meaning there is no need to declare the interfaces which the concrete type satisfies, as long as the concrete type has the required methods, it will be entitled to the interface
- a _concrete type_ (basic types and aggregated types) specifies the exact representation of its values, it is known exactly what it is and what you can do with it
- on other hand, interfaces are _abstract types_, meaning it is not known what it is, only what it can do
- an empty interface is often used to host any value, avoiding the type constraints
- the interface value is made of the type of the value (which different from other types can not be known at compile time) and the value itself
- interfaces can be compared using == and !=, values will be said equal if both are nil, or if dynamic type and dynamic values are identical
- interfaces can be used as keys for map / switch
    - both cases need to be careful if the interface becomes a slice, in this case can cause a panic, as slices are not comparable
- an interface containing a nil value IS NOT a nil interface, and will return false when compared to nil
- a _type assertion_ is an operation applied to an interface value to check its dynamic type and allow to extract the concrete value from the interface or the expand the methods of the interface
    - a typical application is a _type switch_

## variable lifetime
- is the interval of time during it exists as the program executes
- package level variables: lifetime is the entire execution time
- local variables (within functions): lifetime is dynamic, and lives until it becomes unreachable
- compiler may decide to allocate a variable on the 1) heap or 2) stack
    - 1) heap: slower and dynamically managed memory area that stores variables with longer lifetime or unknown sizes at compile time and requires the garbage collector to reclaim memory
    - 2) stack: faster and structured memory area that stores short lived variables and automatically frees memory when a function exits (LIFO - last in, first out)
    - NOTE: when a variable _escapes_ the function scope, it is allocated on heap... this can happen with 1) package declared variables, 2) functions that return pointers, 3) interface values where concrete value is not known at compile time, etc
- NOTE: the decision about variable allocation is made at **compile time**, by go compiler!
- NOTE: each variable that escapes the function scope and is allocated on heap requires extra memory allocation and it can impact performance
- memory leak: a problem that occurs when a program allocates memory for a variable, but fails to release it when no longer needed, causing memory usage to grow unnecessairily and eventually slow down the system or cause it to run out of memory
    - if a long lived object holds for example a pointer to a short lived object, the garbage collector will not be able to release memory from the short lived object
- tools such as `go tool pprof` and `go build -gcflags="m"` may help with investigation and optimization

## pointers 
- a pointer is the address of a variable, that allows to indirectly update values without knowing the variable name
- `pointer := &x` reads as "pointer is assigned the address of x", and now the variable _pointer_, points to x
- `*p = 1` means that the memory block indicated by p address value will be updated with value of 1, as p contains address for x, this will update x = 1
- NOTE: the zero value for a pointer (`*p*`) is **nil**

## goroutines and channels
- golang supports 2 types of concurrency programming: 1) communicating sequential process (CSP) that can be used through goroutines and channels and 2) shared memory multithreading 
### goroutines
- a goroutine is a concurrent function execution
- every go application has at least 1 goroutine, the main goroutine
- new goroutines are created by means of `go` statement
- one goroutine **can not** directly stop others, but **can** send a signal to other goroutines in a way they stop themselves
- when the main finishes or application exits, all goroutines are terminated
### channels
- a channel is a communication mechanism that allows passing values between goroutines
- a channel is created by `ch := make(chan int, capacity int)` and `ch` will be a reference (points to) the data structure created by `make`
- a channel has 3 operations: `ch <- x` (send), `x <- ch` (receive) and `close(ch)`
    - close sends a flag indicating that no more values will ever be sent on this channel, and if that happens (if any value is sent to a closed channel) it will cause a panic
- there are 2 types of channels: 1) unbuffered and 2) buffered, the main difference is that buffered channels have capacity > 1
    - unbuffered channels: when a value is sent in one channel, blocks the channel until the another goroutine executes the receive (one-to-one communication) (forces synced communication)
    - buffered channels: can sent the capacity number of values to the channel without getting blocked
        - the sent is attached to the back of the queue, and the receive always process from top
- a pipeline is a design pattern where the output of one goroutine is the input of another
- a slightly different version of receive is possible, using `x, ok := <- ch`, where the `ok` value can be used to check if a channel stopped sending (closed) and this result can be used to trigger the receiver close
- it is possible to create unidirectional channels 
### multiplexing with select
- select statements are used to choose “which of a set of possible send or receive operations will proceed”
- is similar to a `switch case`, but specific for communication operations
- it is needed to define the cases and the optional default
- the first non-blocking case will be chosen
- if 2 or more cases are not blocking a single one is chosen via an “uniform pseudo-random” selection (almost a load balancing)
- if all cases are blocking, then the default case is chosen
### wait groups and mutexes
- wait groups: a wait group is a synchronization tool provided by standard library 
### race conditions and concurrency problems
- a race condition is a situation in which the program does not give the correct result due to concurrent goroutines attempting to access the same data
- a data race is a type of race condition in which 2 (or more) goroutines access the same variable and at least one of the accesses is a write
- there are 3 ways to avoid a data race:
    1) avoid writting to the variable (when that is possible according to application design)
    2) avoid accessing the variable from multiple goroutines (when that is possible according to application design)
        - a pattern for this is called _serial confination_, in which a variable can be accessed by multiple goroutines, by passing its address from one stage to the next through a channel, creating a pipeline, and ensuring that the variable will not be concurrently accessed
    3) using mutual exclusion to ensure that each goroutine access the variable one at a time 
- one way to deal with race conditions is to use mutex (Lock and Unlock)
    - but, when there is a high number of access attempts, this solution can cause the variable to be locked most part of the time, not allowing any write
    - to overcome this, there is the `sync.RWMutex`, that is allow concurrent access safely and write
    - there is a trade off, RWMutex is slower than regular mutex, so should be useful carefully
- a stale value occurs when a goroutine reads a value that was already updated by another goroutine (and thus is deprecated) but was not yet updated into memory (and thus could not be seem as updated by first goroutine)
    - an way to overcome it is to do a synchronization, using `sync.Once` for example
- go has a tool to help within race detection: `go build -race -o myProgramName main.go`
### examples for goroutines
- usually, an errors channel is also used to communicate errors between goroutines
- it is also a good practice to limit the number of concurrent goroutines, one example of design pattern for this is a _counting semaphore_
- another good practice is to use a cancellation channel, that broadcasts a cancellation signal (by closing the channel) to all goroutines

## testing and benchmarking
- the `go test` command is a test driver for Go
- files must be named using `_test.go` to be considered by the `go test`
- the test files allow 3 use cases: tests, benchmarks and examples
    - tests: identified by a function name that starts with `Test...` and reports a result as either `pass` or `fail`
    - benchmark: identified by a function name that starts with `Benchmark...` and measures the perforamnce of some operation
    - examples: identified by a function name that starts with `Example...` and provides documentation
### tests
- each test file must import the testing package and refer to it in every function
- a table-driven testing is very common, where a struct `tests` made of `input string` and `want`, and then populated with pairs of inputs and expected results
- test messages are usually of the form: `f(x) = y, want z` 
- there are 2 main types of tests, which complement each other
    - black-box test: assumes nothing about the package other than what is specified by its API and documentation
    - white-box test: has privileged access to internal functions and data structures (usually used to verify tricker parts of the implementation)
- tests can be contained in external packages, and are mostly used in integration tests 
- Go also offers the `go cover` tool which allows to verify which sections of the code are covered by tests and which sections are not
### benchmarks
- benchmarks on other hand, are used to measure the application performance
- `go benchmark` is part of `go tests` and measures the number of memory allocations, the number of operations by unit of time and the memory used (`-benchmem` flag)
- benchmarks are often used to:
    - understand the process time accross different request sizes (1, 10, 100, 1000, 10000, ...)
    - understand the best size for an I/O buffer 
    - compare 2 algorithms that does the same job
- profiling is a helper to identify where to begin a code optimization and it is suited when there is no idea for where to begin the investigation
    - profiling is an automated approach to performance measurement
    - Go offers different types of profiles, such as `CPU profile`, `heap profile` and `blocking profile`
        1) `CPU profile`: identifies the functions whose execution requires the most CPU time
        2) `heap profile`: identifies the statements responsible for allocating most memory
        3) `blocking profile`: identifies the operations responsible for blocking goroutines the longest
    - after gathering a profile, it must be analyzed by the pprof tool (`go tool pprof`), that requires the executable that produced the profile and the profile log
### examples
- example functions serves 3 purporses:
    1) documentation
    2) complementary testing by using `Output:` comments at the end
    3) hands on experimentation
## conventions/programming standards
### common 
- Describe each package with a comment immediatelly after declaration
- Use short (implicit) variable declaration (s := "") only within functions and explicit variable declaration (var s string) elsewhere
- **verbs** are string formatters/converters such as: %s (for strings), %d (for integers), %f (for floats), %t (for boolean) and %v (for any value)
- Describe each function with a comment before declaration
- camelCase for function naming
- structs are passed or returned using pointers
- when using mutex, ensure that the variables it guards are not exported
- constructor functions (New) usually returns a pointer
- encapsulation is used to allow internal changes that does not break user code 

### error handling
- if a function has only 1 possible failure cause, it is recommended to use `ok` as checker (boolean type), but if the function has a variety of possible failure causes, it is recommended to use `err` as checker (error type)
    - `error` type is an interface, when non-nil, must have an error message string
- it is the caller responsability to check and take appropriate action on an error
- 5 patterns of error handling:
    1) propagate the error, then the subroutine error becomes an error in the calling routine
    2) retry the failed operation, possible with a delay and limited number of attempts (recommended for transient or unpredictable problems)
    3) log the error and stop the program gracefully (recommended when progress is impossible, usually applied in the main package)
    4) log the error and continue with reduced functionality
    5) log the error and continue without doing anything else (ignore the error)

###  anonymous functions and closure
- an anonymous function is similar to a traditional function except that the programmers did not give that function a name 
- to define the function: `func(){..}` and to define and run it: `func(){..}()`
- an anonymous function defined into another function F can use elements that are not defined in its scope but in the scope of F
- the **closure** happens when the anonymous function references to values that are outer of its own scope
```
func printer() func() {                 // printer scope 
    k := 1                              // printer scope 
    return func() {                     // anonymous function scope
        fmt.Printf("Print n. %d\n", k)  // anonymous function scope
        k++                             // anonymous function scope
    }                                   // anonymous function scope
}

p := printer()  // p will receive the func(){...}, which can see the value k from the printer scope
p() // prints 1, which is the value of k, then increments the value of k to 2  
p() // prints 2, still can see the k from outer scope
```
- if using recursive calls to the anonymous function, the variable that will receive the function must be explicitely defined in the outer scope, example:
```
func funcName (...) {
    var varName func(name1 type1)
    varName := func(name1 type1){
        ...
        varName(...)
        ...
    }
}
```

### deferred function calls
- a defer is an ordinary function or method call prefixed by the keyworld `defer`
- IMPORTANT: the expression is evaluated when the statement is executed, but the actual call is _deferred_ until the function that contains it is finished, example
- IMPORTANT: if there is a os.Exit() or other termination, the defer will never run, which can possible lead to contexts that are not cancelled and still running!
```
func funcName2 (...) func {
    ...
    
    return f // returns a function
}
func funcName1 (...) {
    defer funcName2(...)() // at execution time, this will be evaluated, meaning it will become defer f() 
    ...
}
// but it will only be called after the function finishes!
```

### encapsulation
- encapsulation is used to 1) prevent clients from directly modifying object's variables, 2) hide implementation details and 3) prevent object's variables arbitrarily sets
- this is achieved by getters and setters, which in Golang are methods from structs
- for getters, the usual "Get" prefix is ommited, for setters the "Set" prefix is used

## performance tips
- check memory allocation (heap/stack) and consider if possible to improve to reduce memory allocations
- passing arrays directly to functions results in copying them (memory allocation), if the array is too big or if the array is constantly used, this can impact performance, on other hand, passing the array as pointer, can led to a long living object, in a way need to decide which makes more sense in each case

## useful external packages
### fmt
- contains the print functions 

### os
- allows to read CLI arguments using os.Args
- allows to open/read files using os.Open

### http

### strings

### bytes
- as strings are immutable, building up strings incrementally can involve a lot of allocation and copying, being more efficient to use bytes.Buffer type

### strconv
- provides functions to convert boolean, integer, floating-point to string (and vice-versa)

### unicode
- provides functions such as IsDigit, IsLetter, IsUpper, IsLower to classify runes

### json
- `json.Marshal(struct1)` is used to create a JSON from a structure
- `json.MarshalIndent(struct1, "", " ")` is used to pretty print a JSON 
    - `field tags` such as `omitempty` can be used to trigger special behaviors
- `json.Unmarshal(data, &struct1)` is used to parse a JSON into a struct, the  
- NOTE: only exported struct1 field can be marshaled/unmarshaled, but when Unmarshaling the association between json fields and struct fields is case insensitive

### sync

### context

### mux