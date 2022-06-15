// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2021
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
	"context"
	"runtime"
	"sync"

	"hz.tools/sdr"
	"hz.tools/sdr/fft"
)

type threadsafePlannerRequest struct {
	IQ        sdr.SamplesC64
	Frequency []complex64
	Direction fft.Direction
	Callback  func(fft.Plan, error)
}

func threadsafePlannerSidecar(
	ctx context.Context,
	requests chan threadsafePlannerRequest,
) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for {
		select {
		case req := <-requests:
			req.Callback(Plan(req.IQ, req.Frequency, req.Direction))
		case <-ctx.Done():
			return
		}
	}
}

// ThreadsafePlanner will span a goroutine, lock to a single thread,
// and conduct all FFT planning on that thread. This must be used in
// places where concurrency is a requirement.
func ThreadsafePlanner(ctx context.Context) fft.Planner {
	chn := make(chan threadsafePlannerRequest, 0)
	go threadsafePlannerSidecar(ctx, chn)

	return func(
		iq sdr.SamplesC64,
		freq []complex64,
		direction fft.Direction,
	) (fft.Plan, error) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		var (
			wg   = sync.WaitGroup{}
			plan fft.Plan
			err  error
		)

		wg.Add(1)
		chn <- threadsafePlannerRequest{
			IQ:        iq,
			Frequency: freq,
			Direction: direction,
			Callback: func(cbPlan fft.Plan, cbErr error) {
				defer wg.Done()
				plan = cbPlan
				err = cbErr
			},
		}
		wg.Wait()
		return plan, err
	}
}

// vim: foldmethod=marker
