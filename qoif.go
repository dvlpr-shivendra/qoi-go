package main

import (
	"fmt"
)

type qoifHeader struct {
	height     uint32
	width      uint32
	channels   uint8
	colorspace uint8
}

func (q qoifHeader) String() string {
	return fmt.Sprintf(
		"QoIF Header:\n  Width: %d px\n  Height: %d px\n  Channels: %d\n  Colorspace: %d",
		q.width, q.height, q.channels, q.colorspace,
	)
}

type pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

type scanline []pixel

type qoif struct {
	filePath  string
	header    *qoifHeader
	scanlines []scanline
}

func NewQoif(filePath string) *qoif {
	return &qoif{
		filePath: filePath,
		header: &qoifHeader{
			height:     0,
			width:      0,
			channels:   0,
			colorspace: 0,
		},
		scanlines: make([]scanline, 0),
	}
}

func (q *qoif) Process() error {
	streamer, err := NewStreamer(q.filePath)

	if err != nil {
		return err
	}

	data, err := streamer.Stream(4)

	if err != nil {
		return err
	}

	if string(data) != "qoif" {
		return fmt.Errorf("invalid file")
	}

	width, err := streamer.StreamUInt32()

	if err != nil {
		return err
	}

	q.header.width = width

	height, err := streamer.StreamUInt32()

	if err != nil {
		return err
	}

	q.header.height = height

	channels, err := streamer.StreamUInt8()

	if err != nil {
		return err
	}

	q.header.channels = channels

	colorspace, err := streamer.StreamUInt8()

	if err != nil {
		return err
	}

	q.header.colorspace = colorspace

	fmt.Println(q.header)

	return nil
}
