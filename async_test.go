package async

import (
	"errors"
	"sync/atomic"
	"testing"
)

func Test_RunCount(t *testing.T) {

	var count uint64

	op1 := func() error {
		atomic.AddUint64(&count, 1)
		return nil
	}

	op2 := func() error {
		atomic.AddUint64(&count, 2)
		return nil
	}

	op3 := func() error {
		atomic.AddUint64(&count, 4)
		return nil
	}

	doneChan := make(chan bool)
	errChan := make(chan error)

	q := New(op1, op2)
	q.Add(op3)
	q.Run(doneChan, errChan)

	for {
		select {
		case <-doneChan:
			if count != 7 {
				t.Error("All three channels didn't run properly. count = ", count)
			}
			return
		}
	}
}

func Test_Run(t *testing.T) {
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

	New(op1, op2).Run(doneChan, errChan)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				t.Error("no error expected")
			}
			return
		case <-doneChan:
			if a != 13 {
				t.Error("'a' value is not 13, but is", a)
			}

			if b != 37 {
				t.Error("'b' value is not 37, but is", b)
			}
			return
		}
	}
}

func Test_RunWithError(t *testing.T) {
	var err error

	op1 := func() error {
		return nil
	}

	op2 := func() error {
		return errors.New("OMG FAIL")
	}

	doneChan := make(chan bool)
	errChan := make(chan error)

	New(op1, op2).Run(doneChan, errChan)

LOOP:
	for {
		select {
		case err = <-errChan:
			break LOOP
		}
	}

	if err == nil {
		t.Error("error expected")
	}
}
