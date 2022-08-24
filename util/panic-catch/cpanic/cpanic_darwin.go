//go:build darwin
// +build darwin

/*
Package wppanic Log the panic under mac darwin to the log file
*/

package cpanic

import (
    "fmt"
    "log"
    "os"
    "syscall"
)

// redirectStderr to the file passed in
func redirectStderr(f *os.File) error {
    // CaptureOutputToFd redirects the current process' stdout and stderr file
    // descriptors to the given file descriptor, using the dup2 syscall.
    err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
    if nil != err {
        log.Fatalf("Failed to redirect stderr to file: %v", err)
        fmt.Printf("Failed to redirect stderr to file: %v\n", err)
        return err
    }

    return nil
}
