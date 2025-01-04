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

type Pixel struct {
	r uint8
	g uint8
	b uint8
	a uint8
}

func (p Pixel) String() string {
	return fmt.Sprintf("Pixel(R: %d, G: %d, B: %d, A: %d)", p.r, p.g, p.b, p.a)
}

type Scanline []Pixel

type Qoif struct {
	filePath  string
	header    *qoifHeader
	scanlines []Scanline
}

const (
	QOI_OP_INDEX = 0x00 // 00xxxxxx
	QOI_OP_DIFF  = 0x40 // 01xxxxxx
	QOI_OP_LUMA  = 0x80 // 10xxxxxx
	QOI_OP_RUN   = 0xc0 // 11xxxxxx
	QOI_OP_RGB   = 0xfe // 11111110
	QOI_OP_RGBA  = 0xff // 11111111

	QOI_MASK_2 = 0xc0 // 11000000
)

func NewQoif(filePath string) *Qoif {
	return &Qoif{
		filePath: filePath,
		header: &qoifHeader{
			height:     0,
			width:      0,
			channels:   0,
			colorspace: 0,
		},
		scanlines: make([]Scanline, 0),
	}
}

func (q *Qoif) Process() error {
	streamer, err := NewStreamer(q.filePath)

	if err != nil {
		return err
	}

	err = q.DecodeHeader(streamer)

	if err != nil {
		return err
	}

	return q.PopulateScanlines(streamer)
}

func colorHash(p Pixel) int {
	return int(p.r)*3 + int(p.g)*5 + int(p.b)*7 + int(p.a)*11
}

func (q *Qoif) PopulateScanlines(streamer *Streamer) error {
	q.scanlines = make([]Scanline, q.header.height)
	index := make([]Pixel, 64)
	pixel := &Pixel{r: 0, b: 0, g: 0, a: 255}
	run := 0

	for i := 0; i < len(q.scanlines); i++ {

		scaneline := make([]Pixel, q.header.width)
		for j := 0; j < len(scaneline); j++ {
			fmt.Printf("Reading data for pixel (%d, %d) \n", i, j)
			if run > 0 {
				run--
			}

			chunk, err := streamer.StreamUInt8()
			fmt.Printf("Chunk binary: %08b\n", chunk)
			if err != nil {
				return err
			}

			if chunk == QOI_OP_RGB {
				fmt.Println("Processing QOI_OP_RGB")
				pixel.r, err = streamer.StreamUInt8()
				if err != nil {
					return err
				}
				pixel.g, err = streamer.StreamUInt8()
				if err != nil {
					return err
				}
				pixel.b, err = streamer.StreamUInt8()
				if err != nil {
					return err
				}
				pixel.a = 255
			} else if chunk == QOI_OP_RGBA {
				fmt.Println("Processing QOI_OP_RGBA")
				pixel.r, err = streamer.StreamUInt8()
				if err != nil {
					return err
				}
				pixel.g, err = streamer.StreamUInt8()
				if err != nil {
					return err
				}
				pixel.b, err = streamer.StreamUInt8()
				if err != nil {
					return err
				}
				pixel.a, err = streamer.StreamUInt8()
				if err != nil {
					return err
				}
			} else if (chunk & QOI_MASK_2) == QOI_OP_INDEX {
				fmt.Println("Processing QOI_OP_INDEX")
				pixel = &index[chunk]
			} else if (chunk & QOI_MASK_2) == QOI_OP_DIFF {
				fmt.Println("Processing QOI_OP_DIFF")
				pixel.r += ((chunk >> 4) & 0x03) - 2
				pixel.g += ((chunk >> 2) & 0x03) - 2
				pixel.b += (chunk & 0x03) - 2
			} else if (chunk & QOI_MASK_2) == QOI_OP_LUMA {
				fmt.Println("Processing QOI_OP_LUMA")
				b2, err := streamer.StreamUInt8()
				if err != nil {
					return err
				}
				vg := (chunk & 0x3f) - 32
				pixel.r += vg - 8 + ((b2 >> 4) & 0x0f)
				pixel.g += vg
				pixel.b += vg - 8 + (b2 & 0x0f)
			} else if (chunk & QOI_MASK_2) == QOI_OP_RUN {
				fmt.Println("Processing QOI_OP_RUN")
				run = int((chunk & 0x3f))
			}
			index[colorHash(*pixel)%64] = *pixel
			scaneline[j] = *pixel
		}

		q.scanlines[i] = scaneline
	}

	return nil
}

func (q *Qoif) DecodeHeader(streamer *Streamer) error {
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

	return nil
}
