package async

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
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

	ch := make(chan error)
	q := New(op1, op2)
	q.Add(op3)
	q.Run(ch)

	for {
		select {
		case <-ch:
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

	ch := make(chan error)
	New(op1, op2).Run(ch)

	for {
		select {
		case err := <-ch:
			if err != nil {
				t.Error("no error expected")
			}
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

	ch := make(chan error)
	New(op1, op2).Run(ch)

LOOP:
	for {
		select {
		case err = <-ch:
			if err != nil {
				break LOOP
			}
		}
	}

	if err == nil {
		t.Error("error expected")
	}
}

func TestRunWithTimeout(t *testing.T) {
	var err error
	longOp := func() error {
		time.Sleep(10 * time.Second)
		return nil
	}
	ch := make(chan error)
	New(longOp).RunWithTimeout(ch, time.Millisecond*100)

LOOP:
	for {
		select {
		case err = <-ch:
			break LOOP
		case <-time.Tick(time.Second):
			break LOOP
		}
	}

	if err != ErrTimeout {
		t.Error("timeout error expected")
	}
}
