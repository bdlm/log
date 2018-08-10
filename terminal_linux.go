// Based on ssh/terminal:
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !appengine,!gopherjs

package log

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TCGETS

// Termios contains the unix Termios value.
type Termios unix.Termios
