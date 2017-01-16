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

// Test file for package fd
package fd

import (
	"os"
	"testing"
	"time"

	"gitlab.com/sigma7/fd"
)

var a bool
var b bool

func TestFdRun(t *testing.T) {
	ticker2 := time.NewTicker(time.Second * 5)
	for _ = range ticker2.C {
		fiveSecondExit()
	}

	ticker := time.NewTicker(time.Second * 3)
	f := &fd.Fdcount{Interval: 1, MaxFiles: 4}
	f.Start(trigger)
	a = false
	for _ = range ticker.C {
		if !a {
			t.Error("Error function failed to run")
			t.FailNow()
		} else {
			// do something?
		}
	}
}

func TestFdDoesNotRun(t *testing.T) {
	ticker := time.NewTicker(time.Second * 3)
	f := &fd.Fdcount{Interval: 1, MaxFiles: 60000}
	f.Start(trigger)
	b = false
	for _ = range ticker.C {
		if !b {
			// ???
		} else {
			t.Error("Error function failed to run")
			t.FailNow()
		}
	}
}

func trigger(num int) {
	a = true
}

func trigger2(num int) {
	b = true
}

func fiveSecondExit() {
	os.Exit(0)
}
