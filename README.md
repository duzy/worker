# worker

Package worker implements an easy to use concurrency framework for
multiple-job Go program.

[![GoDoc](https://godoc.org/github.com/duzy/worker?status.svg)](http://godoc.org/github.com/duzy/worker)

Here is a quick example demonstrating the usage.

```go
package example

import "github.com/duzy/worker"

type SomeJob struct {
        Param string
}

func (job *SomeJob) Action() worker.Result {
        // ...
        return &SomeJobResponder{}
}

type SomeJobResponder struct {
}

func (res *SomeJobResponder) Action() {
        // ...
}

const NumberOfConcurrency = 10

func main() {
        w := worker.New()
        w.StartN(NumberOfConcurrency)

        w.Do(&SomeJob{ "anything goes" })
        w.Do(&SomeJob{ "anything goes" })
        w.Do(&SomeJob{ "anything goes" })
        w.Do(&SomeJob{ "anything goes" })

        w.StopN(NumberOfConcurrency)
}
```
