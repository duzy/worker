// Copyright 2016 Duzy Chan <code@duzy.info>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
/*
Package worker implements an easy to use concurrency framework for
multiple-job Go program.

Here is a quick example demonstrating the usage.

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

*/
package worker
