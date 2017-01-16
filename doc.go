/*
Package fd provides a mechanism to monitor the number of open file descriptors
in your program and react if a limit is reached.

EXAMPLE:

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

*/
package fd
