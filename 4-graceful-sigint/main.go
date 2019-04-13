//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	var wg sync.WaitGroup

	sigs := make(chan os.Signal, 1)
	proc := MockProcess{}
	signal.Notify(sigs, syscall.SIGINT)
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Create a process
		// Run the process (blocking)
		proc.Run()
	}()

	select {
		case <- sigs:
			proc.Stop()
	}
	// wait until all routines are finished
	wg.Wait()

}
