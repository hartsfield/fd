// Copyright 2017 J. Hartsfield

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

// It's recommended by [0] [1] [2] to monitor the number of open file
// descriptors in long running, public facing network applications.
// [0] https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
// [1] https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
// [2] https://burke.libbey.me/conserving-file-descriptors-in-go/

// Package fd provides a mechanism to monitor the number of open file
// descriptors in your program and react if a limit is reached. This package is
// not considered 100% idiomatic because it uses an external program to check
// the number of open file descriptors (lsof -p PID), and may only work on
// Linux and OS X.
package fd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Fdcount lets us configure how often we check the number of open files and
// the maximum amount of files we open before running a fail safe function.
type Fdcount struct {
	// Seconds to wait between running countOpenFiles() to check how many file
	// descriptors are open.
	Interval time.Duration
	// Maximum number of files open before running the fail safe function.
	MaxFiles int
}

// NewFDCount returns a new instance of an FdCount{} struct with reasonable
// default values set.
func NewFDCount() *Fdcount {
	// max is obtained by checking /proc/sys/fs/file-max for the systems max
	// number of open file descriptors, and subtracting 50,000 from that max (you
	// don't want to hit the max do you?). If it can't find a number in
	// /proc/sys/fs/file-max it's set to 100,000.
	max := getMax()
	// The run interval is set to 30 seconds.
	return &Fdcount{Interval: 30, MaxFiles: max}
}

// Start starts the checking process, it uses the Fdcount.Interval property to
// determine how often to run the check, 30 seconds is the default but this can
// be tweaked to your individual needs. Running this process often can heavily
// tax the server, so tight intervals are not recommended.
func (fdc *Fdcount) Start(fn func(int)) {
	ticker := time.NewTicker(time.Second * fdc.Interval)
	go func() {
		for _ = range ticker.C {
			if num := countOpenFiles(); num > fdc.MaxFiles {
				fn(num)
			}
		}
	}()
}

// Code from: https://groups.google.com/forum/#!topic/golang-nuts/c0AnWXjzNIA
// seems to be the current best solution for monitoring the number of open file
// descriptors from within a go program.
// NOTE: not idiomatic
func countOpenFiles() int {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -p %v", os.Getpid())).Output()
	if err != nil {
		log.Fatal(err)
	}
	return bytes.Count(out, []byte("\n"))
}

// getMax gets a Unix systems maximum number of open file descriptors
// NOTE: not idiomatic
func getMax() int {
	out, err := exec.Command("cat", "/proc/sys/fs/file-max").Output()
	if err != nil {
		return 100000
	}
	i, err := strconv.Atoi(string(out[:len(out)-1]))
	if err != nil {
		return 100000
	}
	return (i - 50000)
}
