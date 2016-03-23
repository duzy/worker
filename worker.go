package worker

import (
        //"os"
        "sync"
        "errors"
)

// Result represents a delivery of a job done.
type Result interface {
        Action()
}

// Job represents a piece of job to be done.
type Job interface {
        Action() Result
}

type stop struct { waiter *sync.WaitGroup }
func (m *stop) Action() Result { return nil }

// Worker represents a worker to dispatch jobs being done.
type Worker struct {
        i chan Job
        o chan Result
}

// New creates a new worker valid to do user-defined jobs.
func New() *Worker {
        return &Worker{
                //i: make(chan Job, 1),
                //o: make(chan Result, 1), // at least one buffer
        }
}

func (w *Worker) routine(num int) {
        for msg := range w.i {
                if msg != nil {
                        var sw *sync.WaitGroup
                        if stop, ok := msg.(*stop); ok && stop != nil {
                                sw = stop.waiter
                        }
                        w.o <- msg.Action(); go w.reply(sw)
                        if sw != nil { return }
                }
        }
}

func (w *Worker) reply(wg *sync.WaitGroup) {
        if res := <-w.o; res != nil {
                res.Action()
        }
        if wg != nil {
                wg.Done()
        }
}

// Worker.StartN starts a number of `num` threads for jobs.
func (w *Worker) StartN(num int) error {
        if w.i != nil {
                return errors.New("worker is busy on jobs")
        }
        if w.o != nil {
                return errors.New("worker is busy on job results")
        }
        w.i, w.o = make(chan Job, 1), make(chan Result, 1)
        for i := 0; i < num; i++ {
                go w.routine(i)
        }
        return nil
}

// Worker.StopN stops a number of `num` threads.
func (w *Worker) StopN(num int) error {
        if w.i == nil {
                return errors.New("worker free")
        }
        if w.o == nil {
                return errors.New("worker don't have job results")
        }

        c := &stop{ new(sync.WaitGroup) }
        c.waiter.Add(num)
        for i := 0; i < num; i++ {
                w.i <- c
        }
        c.waiter.Wait()
        close(w.i)
        close(w.o)
        w.i, w.o = nil, nil
        return nil
}

// Worker.Do perform a job.
func (w *Worker) Do(m Job) {
        if w.i != nil {
                w.i <- m
        }
}
