// Package files contains various useful functions for handling
// files, or input that may be passed into a program via a named
// pipe.

package files

import (
	"fmt"
	"io"
	"os"
)

// isNamedPipe takes the supplied file (os.Stdin or os.Stdout, typically)
// and returns whether it is a named pipe.
func isNamedPipe(file *os.File) (bool, error) {
	pipe := false
	f, err := file.Stat()
	if err != nil {
		return pipe, err
	}
	if f.Mode()&os.ModeNamedPipe != 0 {
		pipe = true
	}
	return pipe, nil
}

// GetInput returns an io.Reader to capture input - this will either be a file
// from the supplied filename, or os.Stdin.
//
// Input is expected from a filename or a pipe - if no filename is given,
// and os.Stdin is not a named pipe, an error is returned.
func GetInput(filename string, in *os.File) (io.Reader, error) {
	var input io.Reader
	if filename != "" {
		f, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("could not open file: %s", err)
		}
		input = f
	} else {
		pipe, err := isNamedPipe(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("could not determine stdin mode: %s", err)
		}
		if !pipe {
			return nil, fmt.Errorf("nothing passed into stdin")
		}

		input = os.Stdin
	}
	return input, nil
}
