package audio

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
	"os"

	"github.com/pkg/errors"
)

const (
	WAV_HEADER_SIZE = 44
	WAV             = "wav"
)

// WAVHeader wav header (44 bytes)
type WAVHeader struct {
	// RIFF chunk
	RiffMark [4]byte // "RIFF"
	FileSize uint32  // file size-8
	WaveMark [4]byte // "WAVE"
	// fmt chunk
	FmtMark    [4]byte // "fmt "
	FmtSize    uint32  // 16 for PCM
	AudioFmt   uint16  // 1 for PCM
	Channels   uint16  // channels
	SampleRate uint32  // sample rate
	ByteRate   uint32  // byte rate
	BlockAlign uint16  // balk align
	BitsPerSmp uint16  // bits per sample
	// data chunk
	DataMark [4]byte // "data"
	DataSize uint32  // data size
}

// EncodeWavs merges WAV data and encodes it to base64 string
func EncodeWavs(wavs ...[]byte) (string, error) {
	header, mergedData, err := mergeWAVData(wavs)
	if err != nil {
		return "", err
	}

	// create buffer to hold complete WAV file
	buf := new(bytes.Buffer)

	// write header
	if err := writeWAVHeader(buf, header); err != nil {
		return "", err
	}

	// write audio data
	if _, err := buf.Write(mergedData); err != nil {
		return "", errors.WithStack(err)
	}

	// encode to base64
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// WriteWavsToFile writes WAV data to a file
func WriteWavsToFile(output string, wavs ...[]byte) error {
	header, mergedData, err := mergeWAVData(wavs)
	if err != nil {
		return err
	}

	file, err := os.Create(output)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() { _ = file.Close() }()

	if err := writeWAVHeader(file, header); err != nil {
		return err
	}

	if _, err := file.Write(mergedData); err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(file.Sync())
}

// mergeWAVData merges WAV data to a file
func mergeWAVData(wavs [][]byte) (*WAVHeader, []byte, error) {
	var audioSegments [][]byte
	var firstHeader *WAVHeader

	for _, wav := range wavs {
		header, err := readWAVHeader(wav)
		if err != nil {
			return nil, nil, err
		}

		if firstHeader == nil {
			firstHeader = header
		} else {
			if header.SampleRate != firstHeader.SampleRate ||
				header.Channels != firstHeader.Channels ||
				header.BitsPerSmp != firstHeader.BitsPerSmp {
				return nil, nil, errors.New("The format of the WAV files does not match.")
			}
		}

		audioData := wav[WAV_HEADER_SIZE:] // skip header
		audioSegments = append(audioSegments, audioData)
	}

	if len(audioSegments) == 0 {
		return nil, nil, errors.New("No audio clips were successfully synthesized.")
	}

	// merge
	var mergedData []byte
	for _, segment := range audioSegments {
		mergedData = append(mergedData, segment...)
	}

	firstHeader.DataSize = uint32(len(mergedData))
	firstHeader.FileSize = uint32(len(mergedData) + 36) // 36 = 44 - 8

	return firstHeader, mergedData, nil
}

// readWAVHeader reads WAV header from data
func readWAVHeader(data []byte) (*WAVHeader, error) {
	if len(data) < WAV_HEADER_SIZE {
		return nil, errors.New("The data is too short and not a valid WAV file.")
	}

	header := &WAVHeader{}
	buf := bytes.NewReader(data)

	binary.Read(buf, binary.LittleEndian, &header.RiffMark)
	binary.Read(buf, binary.LittleEndian, &header.FileSize)
	binary.Read(buf, binary.LittleEndian, &header.WaveMark)
	binary.Read(buf, binary.LittleEndian, &header.FmtMark)
	binary.Read(buf, binary.LittleEndian, &header.FmtSize)
	binary.Read(buf, binary.LittleEndian, &header.AudioFmt)
	binary.Read(buf, binary.LittleEndian, &header.Channels)
	binary.Read(buf, binary.LittleEndian, &header.SampleRate)
	binary.Read(buf, binary.LittleEndian, &header.ByteRate)
	binary.Read(buf, binary.LittleEndian, &header.BlockAlign)
	binary.Read(buf, binary.LittleEndian, &header.BitsPerSmp)
	binary.Read(buf, binary.LittleEndian, &header.DataMark)
	binary.Read(buf, binary.LittleEndian, &header.DataSize)

	// validate
	if string(header.RiffMark[:]) != "RIFF" || string(header.WaveMark[:]) != "WAVE" {
		return nil, errors.New("The data is not a valid WAV file.")
	}

	return header, nil
}

// writeWAVHeader writes WAV header to io.Writer
func writeWAVHeader(w io.Writer, header *WAVHeader) error {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, header.RiffMark)
	binary.Write(buf, binary.LittleEndian, header.FileSize)
	binary.Write(buf, binary.LittleEndian, header.WaveMark)
	binary.Write(buf, binary.LittleEndian, header.FmtMark)
	binary.Write(buf, binary.LittleEndian, header.FmtSize)
	binary.Write(buf, binary.LittleEndian, header.AudioFmt)
	binary.Write(buf, binary.LittleEndian, header.Channels)
	binary.Write(buf, binary.LittleEndian, header.SampleRate)
	binary.Write(buf, binary.LittleEndian, header.ByteRate)
	binary.Write(buf, binary.LittleEndian, header.BlockAlign)
	binary.Write(buf, binary.LittleEndian, header.BitsPerSmp)
	binary.Write(buf, binary.LittleEndian, header.DataMark)
	binary.Write(buf, binary.LittleEndian, header.DataSize)

	_, err := w.Write(buf.Bytes())

	return err
}
