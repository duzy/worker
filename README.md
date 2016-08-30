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
        return &ContinualJob{}
}

type ContinualJob struct {
}
func (job *ContinualJob) Action() worker.Result {
        // ...
        return &ThirdStep{}
}

type ThirdStep struct {
}
func (job *ThirdStep) Action() worker.Result {
        // ...
        return nil
}

const NumberOfConcurrency = 10

func main() {
        w := worker.SpawnN(NumberOfConcurrency)
        w.Do(&SomeJob{ "anything goes" })
        w.Do(&SomeJob{ "anything goes" })
        w.Do(&SomeJob{ "anything goes" })
        w.Do(&SomeJob{ "anything goes" })
        w.KillAll()

        sentry := w.Sentry()
        sentry.Guard(&SomeJob{ "anything goes" })
        sentry.Guard(&SomeJob{ "anything goes" })
        sentry.Guard(&SomeJob{ "anything goes" })
        for result, _ := range sentry.Wait() {
            // ...
        }
}
```
