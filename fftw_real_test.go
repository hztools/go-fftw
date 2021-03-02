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

package fftw_test

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/cmplx"
	"testing"

	"hz.tools/fftw"
	"hz.tools/rf"
	"hz.tools/sdr/fft"
)

var tau = math.Pi * 2

type testFrequencies struct {
	Frequency rf.Hz
	Index     int
}

func generateRealCw(buf []float32, freq rf.Hz, sampleRate int, phase float64) {
	var (
		carrierFreq float64 = float64(freq)
	)

	for i := range buf {
		now := float64(i) / float64(sampleRate)
		buf[i] = float32(math.Sin(tau*carrierFreq*now + phase))
	}
}

func TestForwardRealFFT(t *testing.T) {
	cw := make([]float32, 1024)
	out := make([]complex64, 1024)

	for _, tfreq := range []testFrequencies{
		testFrequencies{Frequency: rf.Hz(10), Index: 0},
		testFrequencies{Frequency: rf.Hz(900000), Index: 512},
		testFrequencies{Frequency: rf.Hz(450000), Index: 256},
		testFrequencies{Frequency: rf.Hz(225000), Index: 128},
	} {
		generateRealCw(cw, tfreq.Frequency, 1.8e6, 0)

		plan, err := fftw.PlanReal(cw, out, fft.Forward, nil)
		assert.NoError(t, err)
		assert.NoError(t, plan.Transform())
		assert.NoError(t, plan.Close())

		var (
			powerMax float64 = 0
			powerI   int     = -1
		)
		power := make([]float64, len(cw))
		for i := range power {
			power[i] = cmplx.Abs(complex128(out[i]))
			if power[i] > powerMax {
				powerMax = power[i]
				powerI = i
			}
		}
		assert.Equal(t, tfreq.Index, powerI)
	}
}

// vim: foldmethod=marker
