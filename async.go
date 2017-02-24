package async

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrTimeout = errors.New("operation timeout")
)

//	doneChan := make(chan bool)
//  errChan := make(chan error)
//
//  async.New(op1, op2).Run(doneChan, errChan)
//
//  for {
//  	select {
//  	case err := <-errChan:
//  		// handle error
//  	case <-doneChan:
// 			// peform final logic
//  		return
//  	}
//  }

// // Op performs one of the async operations
// type Op func() error

type Ops struct {
	ops []func() error
}

// New accepts a list of operations to be run
func New(fns ...func() error) *Ops {
	aops := Ops{}
	aops.ops = fns

	return &aops
}

// Add operation to be run
func (a *Ops) Add(fn func() error) {
	a.ops = append(a.ops, fn)
}

// Run executes the operation list within a go routine
func (a *Ops) Run(ch chan error) {
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(a.ops))
		for _, op := range a.ops {
			go func(op func() error) {
				if err := op(); err != nil {
					ch <- err
					return
				}
				wg.Done()
			}(op)
		}
		wg.Wait()
		ch <- nil
	}()
}

func (a *Ops) RunWithTimeout(ch chan error, d time.Duration) {
	tch := make(chan error)
	a.Run(tch)
	go func() {
		for {
			select {
			case err := <-tch:
				ch <- err
				return
			case <-time.Tick(d):
				ch <- ErrTimeout
				return
			}
		}
	}()
}
