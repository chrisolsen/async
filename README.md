# Async

Async allows one to eliminates some of clutter resulting when branching operations to multiple channels then merging the results.

### Before

```go
package main

import (
	"fmt"

	"bitbucket.org/chrisolsen/async"
)

func main() {
	var a, b int

	op1 := func() error {
		a = 13
		return nil
	}

	op2 := func() error {
		b = 37
		return nil
	}

    var wg sync.WaitGroup
	doneChan := make(chan bool)
	errChan := make(chan error)

    wg.Add(2)

    go func() {
        wg.Wait()
        doneChan <- true
    }()

    go func() {
        if err := op1(); err != nil {
           errChan <- err 
           return
        }
        wg.Done()
    }()

    go func() {
        if err := op2(); err != nil {
           errChan <- err
           return
        }
        wg.Done()
    }()

	for {
		select {
		case err := <-errChan:
            // handle error
		case <-doneChan:
			fmt.Println(a, b)
			return
		}
	}
}
```

### After

```go
package main

import (
	"fmt"

	"bitbucket.org/chrisolsen/async"
)

func main() {
	var a, b int

	op1 := func() error {
		a = 13
		return nil
	}

	op2 := func() error {
		b = 37
		return nil
	}

	doneChan := make(chan bool)
	errChan := make(chan error)

	async.New(op1, op2).Run(doneChan, errChan)

	for {
		select {
		case err := <-errChan:
			fmt.Println("Error", err.Error())
		case <-doneChan:
			fmt.Println(a, b)
			return
		}
	}
}
```