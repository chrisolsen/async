# Async

Async allows one to eliminates some of clutter resulting when branching operations to multiple channels then merging the results.

### Before

```go
package main

import (
	"fmt"

	"bitbucket.org/chrisolsen/test/async"
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

	last := func() {
		fmt.Println(a, b)
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
			last()
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

	"bitbucket.org/chrisolsen/test/async"
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

	last := func() {
		fmt.Println(a, b)
	}

	doneChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		async.New(op1, op2).Do(doneChan, errChan)
	}()

	for {
		select {
		case err := <-errChan:
			fmt.Println("Error", err.Error())
		case <-doneChan:
			last()
			fmt.Println("DONE")
			return
		}
	}
}
```