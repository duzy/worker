package worker

import (
        "testing"
)

var checkCounter = 0

type increaseNumber struct {
        number int
}

func (job *increaseNumber) Action() Result {
        job.number++
        return &checkResult{ job.number }
}

type checkResult struct {
        number int
}

func (job *checkResult) Action() {
        checkCounter++
}

func TestWorker(t *testing.T) {
        job := new(increaseNumber)
        w := New()
        w.SpawnN(3)
        w.Do(job)
        w.Do(job)
        w.Do(job)
        w.KillAll()
        if job.number != 3 { t.Errorf("wrong job number") }
        if checkCounter != 3 { t.Errorf("wrong number of replies") }
}

func TestSentry(t *testing.T) {
        w := New()
        w.SpawnN(3)

        job := new(increaseNumber)

        sentry := w.Sentry()
        w.Do(job)
        w.Do(job)
        w.Do(job)
        results := sentry.Wait()
        if n, x := len(results), 3; n != x { t.Errorf("%v != %v", n, x) } else {
                for i, result := range results {
                        if res, ok := result.(*checkResult); ok {
                                if res.number != i { t.Errorf("%v != %v", res.number, i) }
                        } else {
                                
                        }
                }
        }

        w.KillAll()
}
