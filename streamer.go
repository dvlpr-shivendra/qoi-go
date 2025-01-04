package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type Streamer struct {
	file *os.File
}

func NewStreamer(filePath string) (*Streamer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return &Streamer{file: file}, nil
}

func (s *Streamer) Stream(byteCount int) ([]byte, error) {
	buffer := make([]byte, byteCount)
	n, err := s.file.Read(buffer)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return buffer[:n], nil
}

func (s *Streamer) StreamUInt32() (uint32, error) {
	data, err := s.Stream(4)
	if err != nil {
		return 0, fmt.Errorf("failed to stream uint32: %w", err)
	}

	if len(data) < 4 {
		return 0, fmt.Errorf("insufficient data: expected 4 bytes, got %d", len(data))
	}

	return binary.BigEndian.Uint32(data), nil
}

func (s *Streamer) StreamUInt8() (uint8, error) {
	data, err := s.Stream(1)
	if err != nil {
		return 0, fmt.Errorf("failed to stream uint8: %w", err)
	}

	if len(data) < 1 {
		return 0, fmt.Errorf("insufficient data: expected 1 byte, got %d", len(data))
	}

	return data[0], nil
}
