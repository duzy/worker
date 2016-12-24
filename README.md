# worker

Package worker implements an easy to use concurrency framework for
multiple-job Go program.

[![GoDoc](https://godoc.org/github.com/duzy/worker?status.svg)](http://godoc.org/github.com/duzy/worker)

Here is a quick example demonstrating the usage.

```go
package example

import "github.com/duzy/worker"

type StepOne struct {
        Param string
}
func (job *StepOne) Action() worker.Result {
        // do job for step one
        return &StepTwo{}
}

type StepTwo struct {
}
func (job *StepTwo) Action() worker.Result {
        // do job for step two
        return &StepThree{}
}

type StepThree struct {
}
func (job *StepThree) Action() worker.Result {
        // do job for step three
        return nil
}

const NumberOfConcurrency = 10

func main() {
        w := worker.SpawnN(NumberOfConcurrency)
        w.Do(&StepOne{ "anything goes" })
        w.Do(&StepOne{ "anything goes" })
        w.Do(&StepOne{ "anything goes" })
        w.Do(&StepOne{ "anything goes" })
        w.Kill()

        w = worker.SpawnN(3)
        sentry := w.Sentry()
        sentry.Guard(&StepOne{ "anything goes" })
        sentry.Guard(&StepOne{ "anything goes" })
        sentry.Guard(&StepOne{ "anything goes" })
        for result, _ := range sentry.Wait() {
            // ...
        }
}
```
