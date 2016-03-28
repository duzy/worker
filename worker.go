package worker

import (
        //"os"
        "sync"
        "errors"
)

// Result represents a delivery of a job done.
type Result interface {}

// ResultJob represents a delivery of a job done and the Action() will
// be called when done.
type ResultJob interface {
        Action()
}

// Job represents a piece of job to be done.
type Job interface {
        Action() Result
}

// Sentry is used to ensure a sequence of jobs are done at a point.
type Sentry struct {
}

type stop struct { waiter *sync.WaitGroup }
func (m *stop) Action() Result { return nil }

// Worker represents a worker to dispatch jobs being done.
type Worker struct {
        routines int
        i chan Job
        o chan Result
}

// New creates a new worker valid to do user-defined jobs.
func New() *Worker {
        return &Worker{}
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
                if j, ok := res.(ResultJob); ok && j != nil {
                        j.Action()
                }
        }
        if wg != nil {
                wg.Done()
        }
}

// Worker.SpawnN starts a number of `num` threads for jobs.
func (w *Worker) SpawnN(num int) error {
        if w.i != nil {
                return errors.New("worker is busy on jobs")
        }
        if w.o != nil {
                return errors.New("worker is busy on job results")
        }
        w.i, w.o, w.routines = make(chan Job, 1), make(chan Result, 1), num
        for i := 0; i < num; i++ {
                go w.routine(i)
        }
        return nil
}

// Worker.KillN stops a number of `num` threads.
func (w *Worker) KillAll() error {
        if w.i == nil {
                return errors.New("worker free")
        }
        if w.o == nil {
                return errors.New("worker don't have job results")
        }

        c := &stop{ new(sync.WaitGroup) }
        c.waiter.Add(w.routines)
        for i := 0; i < w.routines; i++ {
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

func (w *Worker) Sentry() (s *Sentry) {
        s = new(Sentry)
        // TODO: ...
        return
}

func (s *Sentry) Wait() (results []Result) {
        s.WaitFunc(func(result Result){
                results = append(results, result)
        })
        return
}

func (s *Sentry) WaitFunc(f func(result Result)) {
        // TODO: ...
}
