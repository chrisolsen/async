package async

import "sync"

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

// Op performs one of the async operations
type Op func() error

// Ops collection of async operations that needs to be performed
type Ops struct {
	ops []Op
}

// New accepts a list of operations to be run
func New(fns ...Op) *Ops {
	aops := Ops{}
	aops.ops = fns

	return &aops
}

// Add operation to be run
func (a *Ops) Add(fn Op) {
	a.ops = append(a.ops, fn)
}

// Run executes the operation list within a go routine
func (a *Ops) Run(doneChan chan bool, errChan chan error) {
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(a.ops))
		for _, op := range a.ops {
			go func(op Op) {
				if err := op(); err != nil {
					errChan <- err
					return
				}
				wg.Done()
			}(op)
		}
		wg.Wait()
		doneChan <- true
	}()
}
