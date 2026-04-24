package audio

/*
#cgo darwin LDFLAGS: -framework CoreAudio -framework AudioUnit -framework CoreFoundation -lm -lpthread
#cgo linux LDFLAGS: -lm -lpthread -ldl
#cgo windows LDFLAGS: -lm
#cgo CFLAGS: -I.
#include <stdlib.h>
#include "audio.h"
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	MP3 = "mp3"
	WAV = "wav"

	LeftPanning  float32 = -0.05
	RightPanning float32 = 0.05
)

// RenderOption defines the configuration for audio rendering.
type RenderOption struct {
	Pan float32 // Stereo balance: -1.0 (Left) to 1.0 (Right)
}

// parseResult converts a miniaudio result code into a Go error.
func parseResult(res C.int) error {
	if res == 0 {
		return nil
	}
	desc := C.GoString(C.audio_error_description(res))
	return errors.Errorf("audio engine (code %d): %s", int(res), desc)
}

// SaveAudio saves the raw audio data to a file.
func SaveAudio(rawAudio []byte, audioPath string) error {
	err := os.WriteFile(audioPath, rawAudio, 0644)

	return errors.WithStack(err)
}

// TranscodeToWav transcodes it directly to a WAV file.
func Transcode(rawAudio []byte, outputPath string) error {
	if len(rawAudio) == 0 {
		return errors.New("audio engine: decoded data is empty")
	}

	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cOut))

	res := C.audio_transcode(
		(*C.uchar)(unsafe.Pointer(&rawAudio[0])),
		C.size_t(len(rawAudio)),
		cOut,
	)

	return parseResult(res)
}

// RenderFile applies panning effects to an audio file and saves the result.
func RenderFile(inputPath, outputPath string, opt RenderOption) error {
	cIn, cOut := C.CString(inputPath), C.CString(outputPath)
	defer C.free(unsafe.Pointer(cIn))
	defer C.free(unsafe.Pointer(cOut))
	return parseResult(C.audio_render(cIn, cOut, C.float(opt.Pan)))
}

// MergeFiles joins multiple audio files into a single WAV file.
func MergeFiles(inputs []string, outputPath string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("audio engine: no input files provided")
	}
	cOut := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cOut))

	cPaths := make([]*C.char, len(inputs))
	for i, p := range inputs {
		cPaths[i] = C.CString(p)
		defer C.free(unsafe.Pointer(cPaths[i]))
	}
	return parseResult(C.audio_merge(&cPaths[0], C.int(len(inputs)), cOut))
}
