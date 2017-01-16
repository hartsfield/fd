// Copyright 2017 J. Hartzfeldt

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

// TODO: comment tests

// Test file for package fd. fd works by triggering a function when we exceed a
// pre-defined number of open file descriptors. Because fd needs to run for a
// few seconds before the function is run, the test will take that into
// account. Code from user helmbert on stackoverflow.
package fd

import (
	"testing"
	"time"

	"gitlab.com/sigma7/fd"
)

// Test to make sure the function runs when we exceed the maximum number of
// file descriptors.
func TestFd(t *testing.T) {
	functionCalled := make(chan bool)
	timeoutSeconds := 2 * time.Second
	trigger := func(i int) {
		functionCalled <- true
	}

	timeout := time.After(timeoutSeconds)

	// Change MaxFiles to a high number (more than 9 usually) to make this fail.
	c := &fd.Fdcount{Interval: 1, MaxFiles: 1}
	c.Start(trigger)

	select {
	case <-functionCalled:
		t.Logf("function was called")
	case <-timeout:
		t.Fatalf("function was not called within timeout")
	}
}

// Test to make sure the function doesn't run for no reason.
func TestFdNoRun(t *testing.T) {
	functionCalled := make(chan bool)
	timeoutSeconds := 2 * time.Second
	trigger := func(i int) {
		functionCalled <- true
	}

	timeout := time.After(timeoutSeconds)

	// Change MaxFiles to a low number (less than 9 usually) to make this fail.
	c := &fd.Fdcount{Interval: 1, MaxFiles: 5000}
	c.Start(trigger)

	select {
	case <-functionCalled:
		t.Fatalf("function was called when it wasn't supposed to be.")
	case <-timeout:
		t.Logf("function was called")
	}
}
