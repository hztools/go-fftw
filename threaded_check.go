// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2020
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE. }}}

package fftw

import (
	"fmt"
	"syscall"
)

var (
	// ErrPlanOnNonMainThread will be returned if fftw.Plan is invoked from a
	// thread that is not the main OS thread.
	//
	// Plans can not be created in goroutines, due to the way fftw uses those
	// plans.
	//
	ErrPlanOnNonMainThread error = fmt.Errorf("fftw: plans must be created on the main thread")
)

func verifyMainThread() error {
	id := syscall.Gettid()
	pid := syscall.Getpid()

	main := id == pid

	if !main {
		return ErrPlanOnNonMainThread
	}
	return nil
}

// vim: foldmethod=marker
