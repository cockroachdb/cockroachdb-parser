// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package debugutil

import (
	"runtime/debug"
	"sync/atomic"
)

// IsLaunchedByDebugger returns true in cases where the delve debugger
// was used to launch the cockroach process.
func IsLaunchedByDebugger() bool {
	return isLaunchedByDebugger.Load()
}

var isLaunchedByDebugger atomic.Bool

// SafeStack is an alias for []byte that handles redaction. Use this type
// instead of []byte when you are sure that the stack trace does not contain
// sensitive information.
type SafeStack []byte

func (s SafeStack) SafeValue() {}

// Stack wraps the output of debug.Stack() with redact.Safe() to avoid
// unnecessary redaction.
//
// WARNING: Calling this function grabs system-level locks and could cause high
// system CPU usage resulting in the Go runtime to lock up if called too
// frequently, even if called only in error-handling pathways. Use sporadically
// and only when necessary.
func Stack() SafeStack {
	return debug.Stack()
}
