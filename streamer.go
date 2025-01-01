package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type streamer struct {
	file *os.File
}

func NewStreamer(filePath string) (*streamer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return &streamer{file: file}, nil
}

func (s *streamer) Stream(byteCount int) ([]byte, error) {
	buffer := make([]byte, byteCount)
	n, err := s.file.Read(buffer)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return buffer[:n], nil
}

func (s *streamer) StreamUInt32() (uint32, error) {
	data, err := s.Stream(4)
	if err != nil {
		return 0, fmt.Errorf("failed to stream uint32: %w", err)
	}

	if len(data) < 4 {
		return 0, fmt.Errorf("insufficient data: expected 4 bytes, got %d", len(data))
	}

	return binary.BigEndian.Uint32(data), nil
}

func (s *streamer) StreamUInt8() (uint8, error) {
	data, err := s.Stream(1)
	if err != nil {
		return 0, fmt.Errorf("failed to stream uint8: %w", err)
	}

	if len(data) < 1 {
		return 0, fmt.Errorf("insufficient data: expected 1 byte, got %d", len(data))
	}

	return data[0], nil
}
