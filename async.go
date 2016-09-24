package async

import "sync"

//	doneChan := make(chan bool)
//  errChan := make(chan error)
//
//  go func() {
//  	async.New(op1, op2).Do(doneChan, errChan)
//  }()
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

// Do executes the operation list. This function should be called on within a go routine
func (a *Ops) Do(doneChan chan bool, errChan chan error) {
	var wg sync.WaitGroup
	wg.Add(len(a.ops))

	go func() {
		for _, op := range a.ops {
			go func(op Op) {
				if err := op(); err != nil {
					errChan <- err
					return
				}
				wg.Done()
			}(op)
		}
	}()

	wg.Wait()
	doneChan <- true
}
