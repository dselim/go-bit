package bit

import (
	"io"
)

type Writer interface {
	WriteBit(bit Bit) error
	WriteBits(u uint64, nbits int) error
	WriteByte(byt byte) error
	Flush() error
}

type BitWriter struct {
	w     io.Writer
	b     [1]byte
	count uint8
}

func NewWriter(w io.Writer) *BitWriter {
	b := new(BitWriter)
	b.w = w
	b.count = 8
	return b
}

func (b *BitWriter) WriteBit(bit Bit) error {

	if bit {
		b.b[0] |= 1 << (b.count - 1)
	}

	b.count--

	if b.count == 0 {
		if n, err := b.w.Write(b.b[:]); n != 1 || err != nil {
			return err
		}
		b.b[0] = 0
		b.count = 8
	}

	return nil
}

func (b *BitWriter) WriteByte(byt byte) error {
	b.b[0] |= byt >> (8 - b.count)
	if n, err := b.w.Write(b.b[:]); n != 1 || err != nil {
		return err
	}
	b.b[0] = byt << b.count
	return nil
}

func (b *BitWriter) Flush() error {
	for b.count != 8 {
		err := b.WriteBit(Zero)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BitWriter) WriteBits(u uint64, nbits int) error {
	u <<= (64 - uint(nbits))
	for nbits >= 8 {
		byt := byte(u >> 56)
		err := b.WriteByte(byt)
		if err != nil {
			return err
		}
		u <<= 8
		nbits -= 8
	}
	for nbits > 0 {
		err := b.WriteBit((u >> 63) == 1)
		if err != nil {
			return err
		}
		u <<= 1
		nbits--
	}
	return nil
}
