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

// #cgo linux LDFLAGS: -lfftw3f -lm
// #cgo linux CFLAGS:
//
// #include <fftw3.h>
import "C"

import (
	"fmt"
	"unsafe"

	"hz.tools/sdr"
	"hz.tools/sdr/fft"
)

type realPlan struct {
	fftwPlan  C.fftwf_plan
	data      []float32
	frequency []complex64
	opts      opt
	backward  bool
}

func realScaleSamples(s []float32, scaler float32) {
	for x := range s {
		s[x] = s[x] / scaler
	}
}

func (p realPlan) Transform() error {
	C.fftwf_execute(p.fftwPlan)
	if p.opts|OptNoScale != 0 {
		realScaleSamples(p.data, float32(len(p.data)))
	}
	return nil
}

func (p realPlan) Close() error {
	C.fftwf_destroy_plan(p.fftwPlan)
	return nil
}

type RealPlan interface {
	Transform() error
	Close() error
}

// PlanReal will create an interface similar to a hz.tools/sdr/fft.Plan to
// except it's used for real-data to frequency conversions.
func PlanReal(
	samples []float32,
	frequency []complex64,
	direction fft.Direction,
	opts interface{},
) (RealPlan, error) {
	switch direction {
	case fft.Forward:
		if len(frequency) < len(samples) {
			return nil, sdr.ErrDstTooSmall
		}
	case fft.Backward:
		if len(samples) < len(frequency) {
			return nil, sdr.ErrDstTooSmall
		}
	}

	var (
		daPtr   *C.float         = (*C.float)(unsafe.Pointer(&samples[0]))
		fqPtr   *C.fftwf_complex = (*C.fftwf_complex)(unsafe.Pointer(&frequency[0]))
		options opt
	)

	switch opts := opts.(type) {
	case opt:
		options = opts
	}

	switch direction {
	case fft.Forward:
		p := C.fftwf_plan_dft_r2c_1d(C.int(len(samples)), daPtr, fqPtr,
			C.FFTW_ESTIMATE)
		return realPlan{
			opts:      options,
			fftwPlan:  p,
			data:      samples,
			frequency: frequency,
			backward:  false,
		}, nil
	case fft.Backward:
		p := C.fftwf_plan_dft_c2r_1d(C.int(len(frequency)), fqPtr, daPtr,
			C.FFTW_ESTIMATE)
		return realPlan{
			opts:      options,
			fftwPlan:  p,
			data:      samples,
			frequency: frequency,
			backward:  true,
		}, nil
	}

	return nil, fmt.Errorf("hz.tools/fftw: unreachable code")
}

// vim: foldmethod=marker
