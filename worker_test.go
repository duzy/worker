package worker

import (
        "testing"
        "sync"
)

var (
        checkCounter = 0
        checkMutex = new(sync.Mutex)
)

type increaseNumber struct {
        mutex *sync.Mutex
        number int
}

func (job *increaseNumber) Go(n int) Result {
        job.mutex.Lock(); defer job.mutex.Unlock()
        job.number++
        return &checkResult{ job.number }
}

type checkResult struct {
        number int
}

func (job *checkResult) Go(n int) Result {
        checkMutex.Lock(); defer checkMutex.Unlock()
        checkCounter++
        return nil
}

func TestWorker(t *testing.T) {
        job := new(increaseNumber)
        job.mutex = new(sync.Mutex)
        w, num := New(), 10000
        w.SpawnN(3)
        for i := 0; i < num; i++ {
                w.Do(job)
        }
        w.Stop()
        if job.number != num { t.Errorf("wrong job number") }
        if checkCounter != num { t.Errorf("wrong number of replies") }
}

func TestSentry(t *testing.T) {
        w := New()
        w.SpawnN(5)

        job := new(increaseNumber)
        job.mutex = new(sync.Mutex)

        sentry, num := w.Sentry(), 100000
        for i := 0; i < num; i++ {
                sentry.Guard(job)
        }
        results := sentry.Wait()
        if n := len(results); n != num { t.Errorf("%v != %v", n, num) } else {
                for i, result := range results {
                        if res, ok := result.(*checkResult); ok {
                                //if res.number != i { t.Errorf("%v != %v", res.number, i) }
                                if res.number <= 0 || num < res.number {
                                        t.Errorf("%v (%v)", res.number, i)
                                }
                        } else {
                                // ...
                        }
                }
        }

        w.Stop()
}

var (
        job0Executed = 0
        job1Executed = 0
        job2Executed = 0
        job0Mutex = new(sync.Mutex)
        job1Mutex = new(sync.Mutex)
        job2Mutex = new(sync.Mutex)
)
type job0 struct {
        tag string
}
func (job *job0) Go(n int) Result {
        job0Mutex.Lock(); defer job0Mutex.Unlock()
        job0Executed++
        return new(job1)
}

type job1 struct {
        tag string
}
func (job *job1) Go(n int) Result {
        job1Mutex.Lock(); defer job1Mutex.Unlock()
        job1Executed++
        return new(job2)
}

type job2 struct {
        tag string
}
func (job *job2) Go(n int) Result {
        job2Mutex.Lock(); defer job2Mutex.Unlock()
        job2Executed++
        return "done"
}

func TestJobChain(t *testing.T) {
        job := new(job0)
        
        w, num := New(), 10000
        w.SpawnN(3)
        for i := 0; i < num; i++ {
                w.Do(job)
        }
        w.Stop()
        
        if job0Executed != num { t.Errorf("job0: %v != %v", job0Executed, num) }
        if job1Executed != num { t.Errorf("job1: %v != %v", job1Executed, num) }
        if job2Executed != num { t.Errorf("job2: %v != %v", job2Executed, num) }
        
        job0Executed = 0
        job1Executed = 0
        job2Executed = 0

        w = SpawnN(1)
        sentry := w.Sentry()
        sentry.Guard(&job0{"job0"})
        sentry.Guard(&job1{"job1"})
        sentry.Guard(&job2{"job2"})
        for _, result := range sentry.Wait() {
                if s, ok := result.(string); !ok || s != "done" {
                        t.Errorf("%v != done", result)
                }
        }
        if job0Executed != 1 { t.Errorf("job0: %v != %v", job0Executed, 1) }
        if job1Executed != 2 { t.Errorf("job1: %v != %v", job1Executed, 2) }
        if job2Executed != 3 { t.Errorf("job2: %v != %v", job2Executed, 3) }
}
