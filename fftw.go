// {{{ Copyright (c) Paul R. Tagliamonte <paul@k3xec.com>, 2020-2021
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

// #cgo pkg-config: fftw3f
//
// #include <fftw3.h>
import "C"

import (
	"fmt"
	"unsafe"

	"hz.tools/sdr"
	"hz.tools/sdr/fft"
)

type plan struct {
	fftwPlan  C.fftwf_plan
	iq        sdr.SamplesC64
	frequency []complex64
	backward  bool
}

func (p plan) Transform() error {
	C.fftwf_execute(p.fftwPlan)
	p.iq.Scale(1 / float32(len(p.iq)))
	return nil
}

func (p plan) Close() error {
	C.fftwf_destroy_plan(p.fftwPlan)
	return nil
}

// Plan will create a hz.tools/sdr/fft.Plan to be used to perform
// frequency-to-time or time-to-frequency conversions of complex data.
func Plan(
	iq sdr.SamplesC64,
	frequency []complex64,
	direction fft.Direction,
) (fft.Plan, error) {
	switch direction {
	case fft.Forward:
		if len(frequency) < len(iq) {
			return nil, sdr.ErrDstTooSmall
		}
	case fft.Backward:
		if len(iq) < len(frequency) {
			return nil, sdr.ErrDstTooSmall
		}
	}

	var (
		iqPtr *C.fftwf_complex = (*C.fftwf_complex)(unsafe.Pointer(&iq[0]))
		fqPtr *C.fftwf_complex = (*C.fftwf_complex)(unsafe.Pointer(&frequency[0]))
	)

	switch direction {
	case fft.Forward:
		p := C.fftwf_plan_dft_1d(C.int(iq.Length()), iqPtr, fqPtr,
			C.FFTW_FORWARD, C.FFTW_ESTIMATE)
		return plan{
			fftwPlan:  p,
			iq:        iq,
			frequency: frequency,
			backward:  false,
		}, nil
	case fft.Backward:
		p := C.fftwf_plan_dft_1d(C.int(len(frequency)), fqPtr, iqPtr,
			C.FFTW_BACKWARD, C.FFTW_ESTIMATE)
		return plan{
			fftwPlan:  p,
			iq:        iq,
			frequency: frequency,
			backward:  true,
		}, nil
	}

	return nil, fmt.Errorf("hz.tools/fftw: unreachable code")
}

// vim: foldmethod=marker
