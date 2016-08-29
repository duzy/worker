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

func (job *increaseNumber) Action() Result {
        job.mutex.Lock(); defer job.mutex.Unlock()
        job.number++
        return &checkResult{ job.number }
}

type checkResult struct {
        number int
}

func (job *checkResult) Action() {
        checkMutex.Lock(); defer checkMutex.Unlock()
        checkCounter++
}

func TestWorker(t *testing.T) {
        job := new(increaseNumber)
        job.mutex = new(sync.Mutex)
        w, num := New(), 10000
        w.SpawnN(3)
        for i := 0; i < num; i++ {
                w.Do(job)
        }
        w.KillAll()
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

        w.KillAll()
}
