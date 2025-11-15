package audio

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
)

const MP3 = "mp3"

// MP3FrameHeader represents an MP3 frame header
type MP3FrameHeader struct {
	Sync         uint16
	Version      int
	Layer        int
	BitrateIndex int
	SamplingRate int
	Padding      bool
	FrameSize    int
}

// MP3Result contains the merged MP3 data and its duration
type MP3Result struct {
	Data     []byte
	Duration time.Duration
}

// EncodeMp3s merges multiple MP3 files and returns base64 encoded data with duration
func EncodeMp3s(mp3s ...[]byte) (string, time.Duration, error) {
	result, err := mergeMp3Data(mp3s...)
	if err != nil {
		return "", 0, err
	}

	encoded := base64.StdEncoding.EncodeToString(result.Data)
	return encoded, result.Duration, nil
}

// WriteMp3sToFile writes MP3 data to a file and returns the duration
func WriteMp3sToFile(output string, mp3s ...[]byte) (time.Duration, error) {
	result, err := mergeMp3Data(mp3s...)
	if err != nil {
		return 0, err
	}

	file, err := os.Create(output)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if _, err := file.Write(result.Data); err != nil {
		return 0, errors.WithStack(err)
	}

	if err := errors.WithStack(file.Sync()); err != nil {
		return 0, err
	}

	return result.Duration, nil
}

// mergeMp3Data merges multiple MP3 data into a single file and calculates duration
func mergeMp3Data(mp3s ...[]byte) (*MP3Result, error) {
	if len(mp3s) == 0 {
		return nil, fmt.Errorf("no mp3 data provided")
	}

	var mergedData bytes.Buffer
	var totalDuration time.Duration

	// Process first MP3 - keep ID3v2 tag if present
	firstMP3 := mp3s[0]
	id3v2Size := getID3v2Size(firstMP3)

	if id3v2Size > 0 {
		// Write ID3v2 tag from first file
		mergedData.Write(firstMP3[:id3v2Size])
	}

	// Extract and merge audio frames from all MP3s
	for i, mp3Data := range mp3s {
		frames, duration, err := extractMP3Frames(mp3Data)
		if err != nil {
			return nil, fmt.Errorf("failed to extract frames from mp3 #%d: %w", i, err)
		}
		mergedData.Write(frames)
		totalDuration += duration
	}

	return &MP3Result{
		Data:     mergedData.Bytes(),
		Duration: totalDuration,
	}, nil
}

// getID3v2Size returns the size of ID3v2 tag (returns 0 if no ID3v2 tag)
func getID3v2Size(data []byte) int {
	if len(data) < 10 {
		return 0
	}

	// Check for ID3v2 header: "ID3"
	if data[0] != 'I' || data[1] != 'D' || data[2] != '3' {
		return 0
	}

	// ID3v2 size is stored in bytes 6-9 (synchsafe integer)
	// Each byte uses only 7 bits
	size := int(data[6]&0x7F)<<21 | int(data[7]&0x7F)<<14 | int(data[8]&0x7F)<<7 | int(data[9]&0x7F)

	// Add 10 bytes for the header itself
	return size + 10
}

// extractMP3Frames extracts pure MP3 audio frames and calculates duration
func extractMP3Frames(data []byte) ([]byte, time.Duration, error) {
	if len(data) == 0 {
		return nil, 0, fmt.Errorf("empty mp3 data")
	}

	var frames bytes.Buffer
	var totalDuration time.Duration
	offset := 0

	// Skip ID3v2 tag at the beginning
	id3v2Size := getID3v2Size(data)
	if id3v2Size > 0 {
		offset = id3v2Size
	}

	// Extract frames until we hit ID3v1 tag or end of file
	for offset < len(data) {
		// Check for ID3v1 tag (at the end, starts with "TAG")
		if len(data)-offset >= 128 &&
			data[offset] == 'T' && data[offset+1] == 'A' && data[offset+2] == 'G' {
			break // Skip ID3v1 tag
		}

		// Look for frame sync (11 bits set to 1: 0xFFE or 0xFFF)
		if offset+4 > len(data) {
			break
		}

		// Check for MP3 frame sync marker (0xFF 0xFx or 0xFF 0xEx)
		if data[offset] != 0xFF {
			offset++
			continue
		}

		// Validate it's actually a frame header
		if (data[offset+1] & 0xE0) != 0xE0 {
			offset++
			continue
		}

		// Parse frame header to get frame size
		header, err := parseMP3FrameHeader(data[offset:])
		if err != nil || header.FrameSize == 0 {
			// A more robust skip: maybe try to find the next sync word (optional)
			// But for now, just skip 1 byte.
			offset++
			continue
		}

		// Make sure we have enough data for this frame
		if offset+header.FrameSize > len(data) {
			break
		}

		// Write the complete frame
		frames.Write(data[offset : offset+header.FrameSize])

		// Duration = (SamplesPerFrame / SamplingRate) * time.Second
		samplesPerFrame := getSamplesPerFrame(header.Version, header.Layer)
		if samplesPerFrame > 0 && header.SamplingRate > 0 {
			nanoseconds := (int64(samplesPerFrame) * int64(time.Second)) / int64(header.SamplingRate)
			frameDuration := time.Duration(nanoseconds)
			totalDuration += frameDuration
		}

		offset += header.FrameSize
	}

	if frames.Len() == 0 {
		// ID3v1 at the end of the file is not an issue if frames were found.
		if id3v2Size > 0 && offset == id3v2Size && len(data) > id3v2Size {
			// We had a tag but found no frames, maybe it's VBR (Variable Bitrate)
			// and the current sync method is failing, or it's just corrupted/empty.
			return nil, 0, fmt.Errorf("could not find valid MP3 frames after ID3v2 tag")
		}

		// Final check only returns an error if no frames were found at all.
		if offset == 0 || offset == id3v2Size {
			return nil, 0, fmt.Errorf("no valid MP3 frames found")
		}
	}

	return frames.Bytes(), totalDuration, nil
}

// getSamplesPerFrame returns the number of samples per frame for given version and layer
func getSamplesPerFrame(version, layer int) int {
	if layer == 1 {
		// Layer I
		return 384
	} else if layer == 2 {
		// Layer II
		return 1152
	} else if layer == 3 {
		// Layer III
		if version == 1 {
			// MPEG 1
			return 1152
		} else {
			// MPEG 2 and 2.5
			return 576
		}
	}
	return 0
}

// parseMP3FrameHeader parses MP3 frame header and calculates frame size
func parseMP3FrameHeader(data []byte) (*MP3FrameHeader, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("insufficient data for frame header")
	}

	header := &MP3FrameHeader{}

	// Byte 0-1: Sync word (11 bits all set to 1)
	header.Sync = uint16(data[0])<<8 | uint16(data[1])
	if (header.Sync & 0xFFE0) != 0xFFE0 {
		return nil, fmt.Errorf("invalid sync word")
	}

	// Byte 1: MPEG version (bits 3-4)
	versionBits := (data[1] >> 3) & 0x03
	switch versionBits {
	case 0:
		header.Version = 25 // MPEG 2.5
	case 2:
		header.Version = 2 // MPEG 2
	case 3:
		header.Version = 1 // MPEG 1
	default:
		return nil, fmt.Errorf("invalid MPEG version")
	}

	// Byte 1: Layer (bits 1-2)
	layerBits := (data[1] >> 1) & 0x03
	switch layerBits {
	case 1:
		header.Layer = 3
	case 2:
		header.Layer = 2
	case 3:
		header.Layer = 1
	default:
		return nil, fmt.Errorf("invalid layer")
	}

	// Byte 2: Bitrate index (bits 4-7)
	header.BitrateIndex = int(data[2] >> 4)
	if header.BitrateIndex == 0 || header.BitrateIndex == 15 {
		return nil, fmt.Errorf("invalid bitrate index")
	}

	// Byte 2: Sampling rate index (bits 2-3)
	samplingRateIndex := int((data[2] >> 2) & 0x03)
	if samplingRateIndex == 3 {
		return nil, fmt.Errorf("invalid sampling rate")
	}

	// Byte 2: Padding (bit 1)
	header.Padding = (data[2] & 0x02) != 0

	// Calculate actual bitrate (kbps)
	bitrate := getBitrate(header.Version, header.Layer, header.BitrateIndex)
	if bitrate == 0 {
		return nil, fmt.Errorf("invalid bitrate")
	}

	// Calculate actual sampling rate (Hz)
	samplingRate := getSamplingRate(header.Version, samplingRateIndex)
	if samplingRate == 0 {
		return nil, fmt.Errorf("invalid sampling rate")
	}
	header.SamplingRate = samplingRate

	// Calculate frame size
	padding := 0
	if header.Padding {
		if header.Layer == 1 {
			padding = 4
		} else {
			padding = 1
		}
	}

	if header.Layer == 1 {
		// Layer I
		header.FrameSize = ((12*bitrate*1000)/samplingRate + padding/4) * 4
	} else {
		// Layer II & III
		header.FrameSize = (144*bitrate*1000)/samplingRate + padding
	}

	return header, nil
}

// getBitrate returns bitrate in kbps
func getBitrate(version, layer, bitrateIndex int) int {
	// Bitrate table [version][layer][index]
	bitrateTable := [][][]int{
		// MPEG 2.5
		{
			{}, // Layer reserved
			{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0},      // Layer III
			{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0},      // Layer II
			{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0}, // Layer I
		},
		// MPEG 2
		{
			{}, // Layer reserved
			{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0},      // Layer III
			{0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160, 0},      // Layer II
			{0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256, 0}, // Layer I
		},
		// MPEG 1
		{
			{}, // Layer reserved
			{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0},     // Layer III
			{0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384, 0},    // Layer II
			{0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448, 0}, // Layer I
		},
	}

	versionIdx := 0
	switch version {
	case 1:
		versionIdx = 2
	case 2:
		versionIdx = 1
	case 25:
		versionIdx = 0
	default:
		return 0
	}

	if versionIdx >= len(bitrateTable) || layer >= len(bitrateTable[versionIdx]) {
		return 0
	}

	table := bitrateTable[versionIdx][layer]
	if bitrateIndex >= len(table) {
		return 0
	}

	return table[bitrateIndex]
}

// getSamplingRate returns sampling rate in Hz
func getSamplingRate(version, samplingRateIndex int) int {
	samplingRateTable := [][]int{
		{11025, 12000, 8000, 0},  // MPEG 2.5
		{0, 0, 0, 0},             // Reserved
		{22050, 24000, 16000, 0}, // MPEG 2
		{44100, 48000, 32000, 0}, // MPEG 1
	}

	versionIdx := 0
	switch version {
	case 1:
		versionIdx = 3
	case 2:
		versionIdx = 2
	case 25:
		versionIdx = 0
	default:
		return 0
	}

	if versionIdx >= len(samplingRateTable) || samplingRateIndex >= len(samplingRateTable[versionIdx]) {
		return 0
	}

	return samplingRateTable[versionIdx][samplingRateIndex]
}
