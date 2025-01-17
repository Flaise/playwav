// The MIT License (MIT)
//
// Copyright (c) 2017 aerth <aerth@riseup.net>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// playwav plays a .wav file using ALSA
package playwav

import (
	"errors"
	"fmt"
	"io"
	"os"

	alsa "github.com/cocoonlife/goalsa"
	"github.com/cryptix/wav"
)

func FromReader(reader io.ReadSeeker, size int64) error {

	// wavReader
	wavReader, err := wav.NewReader(reader, size)
	if err != nil {
		return errors.New(fmt.Sprint("WAV reader:", err))
	}

	// require wavReader
	if wavReader == nil {
		return errors.New("wav reader is nil")
	}

	fileinfo := wavReader.GetFile()
	// open default ALSA playback device
	samplerate := int(fileinfo.SampleRate)
	if samplerate == 0 {
		samplerate = 44100
	}
	if samplerate > 100000 {
		samplerate = 44100
	}

	// Temporary fix. Now plays at correct pitch but wrong duration.
	samplerate /= 2

	out, err := alsa.NewPlaybackDevice("default", 1, alsa.FormatS16LE, samplerate, alsa.BufferParams{})
	if err != nil {
		return errors.New(fmt.Sprint("ALSA:", err))
	}

	// require ALSA device
	if out == nil {
		return errors.New("nil ALSA device")
	}

	// close device when finished
	defer out.Close()

	for {
		s, err := wavReader.ReadSampleEvery(2, 0)
		var cvert []int16
		for _, b := range s {
			cvert = append(cvert, int16(b))
		}
		if cvert != nil {
			// play!
			out.Write(cvert)
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return errors.New(fmt.Sprint("WAV Decode:", err))
		}
	}

	return nil
}

func FromFile(filename string) error {
	// file exists
	file, err := os.Open(filename)
	if err != nil {
		return errors.New(fmt.Sprint("open:", err))
	}

	// stat for size
	sndfileinfo, err := os.Stat(file.Name())
	if err != nil {
		return errors.New(fmt.Sprint("stat:", err))
	}
	size := sndfileinfo.Size()

	return FromReader(file, size)
}
