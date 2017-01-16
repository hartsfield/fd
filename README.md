<h1>fd: A Go utility for monitoring the number of open file descriptors from
within a program.</h1>

You've probably noticed that some long lived Go programs such as http servers
can end up exceeding the systems maximum number of open file descriptors. 
Package fd provides a mechanism to monitor the number of open file descriptors
in your program and react if a limit is reached. It uses `os/exec.Command` to
run `lsof -p PID`, which is an external command found on Unix and Linux, thus
there are plenty of other ways to monitor the number of open file descriptors
of a program on your server, and it might even be better to just write a bash
script depending on your use case.

EXAMPLE:

```
// This example runs for ever and should always give output every second.
package main

import (
	"fmt"

	"gitlab.com/zfeldt/fd"
)

func main() {
	// Get a custom Fdcount struct (could also use fd.NewFDCount() for an Fdcount
  // struct that uses default values).
  f := &fd.Fdcount{Interval: 1, MaxFiles: 2}
	// Start watching the number of open file descriptors
	f.Start(trigger)

	// this is just here to hold the program open so you can see output
	ch := make(chan bool, 1)
	<-ch
}

// This function will be run when the number of file descriptors goes above
// Fdcount.MaxFiles.
func trigger(num int) {
	fmt.Printf("You have %v file descriptors open.\n", num)
}
```
