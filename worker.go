package worker

import (
        //"os"
        //"sync"
)

// Result represents a delivery of a job done.
type Result interface {
        Action()
}

// Job represents a piece of job to be done.
type Job interface {
        Action() Result
}

type stop struct { waiter chan int }
func (m *stop) Action() Result { return nil }

// Worker represents a worker to dispatch jobs being done.
type Worker struct {
        i chan Job
        o chan Result
}

// New creates a new worker valid to do user-defined jobs.
func New() *Worker {
        return &Worker{
                i: make(chan Job, 1),
                o: make(chan Result, 1), // at least one buffer
        }
}

func (w *Worker) routine(num int) {
        for msg := range w.i {
                if msg == nil { continue }
                w.o <- msg.Action(); go w.reply()
                if stop, ok := msg.(*stop); ok {
                        stop.waiter <- num
                        break
                }
        }
}

func (w *Worker) reply() {
        if res := <-w.o; res != nil {
                res.Action()
        }
}

// Worker.StartN starts a number of `num` threads for jobs.
func (w *Worker) StartN(num int) {
        for i := 0; i < num; i++ {
                go w.routine(i)
        }
}

// Worker.StopN stops a number of `num` threads.
func (w *Worker) StopN(num int) {
        barrier := make(chan int)
        for i := 0; i < num; i++ {
                w.i <- &stop{ barrier }
        }
        for i := 0; i < num; i++ {
                _ = <-barrier
        }
        //close(w.i)
        //close(w.o)
}

// Worker.Do perform a job.
func (w *Worker) Do(m Job) {
        w.i <- m
}
