// +build darwin freebsd openbsd netbsd dragonfly
// +build !appengine,!gopherjs

package log

import "golang.org/x/sys/unix"

// Termios contains the unix Termios value.
type Termios unix.Termios
