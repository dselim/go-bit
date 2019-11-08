package bit

import (
	"io"
)

type Reader interface {
	ReadBit() (Bit, error)
	ReadBits(nbits int) (uint64, error)
	ReadByte() (byte, error)
	Reset()
}

type BitReader struct {
	r     io.Reader
	b     [1]byte
	count uint8
}

func NewReader(r io.Reader) *BitReader {
	b := new(BitReader)
	b.r = r
	return b
}

func (b *BitReader) ReadBit() (Bit, error) {
	if b.count == 0 {
		if n, err := b.r.Read(b.b[:]); n != 1 || (err != nil && err != io.EOF) {
			return Zero, err
		}
		b.count = 8
	}
	b.count--
	d := (b.b[0] & 0x80)
	b.b[0] <<= 1
	return d != 0, nil
}

func (b *BitReader) ReadByte() (byte, error) {
	if b.count == 0 {
		n, err := b.r.Read(b.b[:])
		if n != 1 || (err != nil && err != io.EOF) {
			b.b[0] = 0
			return b.b[0], err
		}
		// mask io.EOF for the last byte
		if err == io.EOF {
			err = nil
		}
		return b.b[0], err
	}

	byt := b.b[0]

	var n int
	var err error
	n, err = b.r.Read(b.b[:])
	if n != 1 || (err != nil && err != io.EOF) {
		return 0, err
	}

	byt |= b.b[0] >> b.count

	b.b[0] <<= (8 - b.count)

	return byt, err
}

func (b *BitReader) ReadBits(nbits int) (uint64, error) {
	var u uint64
	for nbits >= 8 {
		byt, err := b.ReadByte()
		if err != nil {
			return 0, err
		}

		u = (u << 8) | uint64(byt)
		nbits -= 8
	}
	var err error
	for nbits > 0 && err != io.EOF {
		byt, err := b.ReadBit()
		if err != nil {
			return 0, err
		}
		u <<= 1
		if byt {
			u |= 1
		}
		nbits--
	}
	return u, nil
}

func (b *BitReader) Reset() {
	b.b[0] = 0x00
	b.count = 0
}
