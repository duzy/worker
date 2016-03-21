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
        w.StartN(3)
        w.Do(job)
        w.Do(job)
        w.Do(job)
        w.StopN(3)
        if job.number != 3 { t.Errorf("wrong job number") }
        if checkCounter != 3 { t.Errorf("wrong number of replies") }
}
