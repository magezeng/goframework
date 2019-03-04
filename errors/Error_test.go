package errors

import (
	"fmt"
	"io"
	"testing"
)

func TestNewError(t *testing.T) {

	e := func() error {
		return New("hi")
	}()

	if e.Error() != "hi" {
		t.Errorf("Constructor with a string failed")
	}

	if New(fmt.Errorf("yo")).Error() != "yo" {
		t.Errorf("Constructor with an error failed")
	}

	if New(e) != e {
		t.Errorf("Constructor with an Error failed")
	}

	if New(nil) != nil {
		t.Errorf("Constructor with nil failed")
	}
}

func TestIs(t *testing.T) {

	if Is(nil, io.EOF) {
		t.Errorf("nil is an error")
	}

	if !Is(io.EOF, io.EOF) {
		t.Errorf("io.EOF is not io.EOF")
	}

	if Is(io.EOF, fmt.Errorf("io.EOF")) {
		t.Errorf("io.EOF is fmt.Errorf")
	}
}

func TestPrintStack(t *testing.T) {
	defer func() {
		err := recover()
		e := New(err)
		t.Log(e.frames)
	}()

	a()
}

func a() error {
	go b(5)
	return nil
}

func b(i int) {
	go c()
}

func c() {
	d()
}

func d() {
	panic("出错了！")
}
