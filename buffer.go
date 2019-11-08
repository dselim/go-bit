package bit

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// buffered bit reader/writer

type BufReader struct {
	br *bufio.Reader
	*BitReader
}

func NewBufReader(r io.Reader) *BufReader {
	br := bufio.NewReader(r)
	return &BufReader{
		br:        br,
		BitReader: NewReader(br),
	}
}

type BufWriter struct {
	bw *bufio.Writer
	*BitWriter
}

func NewBufWriter(r io.Writer) *BufWriter {
	bw := bufio.NewWriter(r)
	return &BufWriter{
		bw:        bw,
		BitWriter: NewWriter(bw),
	}
}

func (w *BufWriter) Flush() error {
	err := w.BitWriter.Flush()
	if err != nil {
		return err
	}
	return w.bw.Flush()
}

type BufReadWriter struct {
	*BufReader
	*BufWriter
}

func NewBufReadWriter(r *BufReader, w *BufWriter) *BufReadWriter {
	return &BufReadWriter{r, w}
}

// bit buffer

type Buffer struct {
	*BitReadWriter
	buf *bytes.Buffer
}

func NewBuffer() *Buffer {
	buf := new(bytes.Buffer)
	return &Buffer{
		BitReadWriter: NewReadWriter(NewReader(buf), NewWriter(buf)),
		buf:           buf,
	}
}

func (b *Buffer) Copy() *Buffer {
	data := make([]byte, len(b.buf.Bytes()), b.buf.Cap())
	copy(data, b.buf.Bytes())
	buf := bytes.NewBuffer(data)
	return &Buffer{
		BitReadWriter: NewReadWriter(NewReader(buf), NewWriter(buf)),
		buf:           buf,
	}
}

func (b *Buffer) Reset() {
	b.Flush()
	b.buf.Reset()
	b.BitReader.Reset()
}

func (b *Buffer) ReadFrom(r Reader, nbits int) error {
	for nbits >= 8 {
		v, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("bit buffer error: read from bit reader: %w", err)
		}
		err = b.WriteByte(v)
		if err != nil {
			return fmt.Errorf("bit buffer error: write to bit buffer: %w", err)
		}
		nbits -= 8
	}
	if nbits > 0 {
		v, err := r.ReadBits(nbits)
		if err != nil {
			return fmt.Errorf("bit buffer error: read from bit reader: %w", err)
		}
		err = b.WriteBits(v, nbits)
		if err != nil {
			return fmt.Errorf("bit buffer error: write to bit buffer: %w", err)
		}
	}
	return nil
}

func (b *Buffer) WriteTo(w Writer, nbits int) error {
	for nbits >= 8 {
		v, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("bit buffer error: read from bit buffer: %w", err)
		}
		err = w.WriteByte(v)
		if err != nil {
			return fmt.Errorf("bit buffer error: write to bit writer: %w", err)
		}
		nbits -= 8
	}
	if nbits > 0 {
		v, err := b.ReadBits(nbits)
		if err != nil {
			return fmt.Errorf("bit buffer error: read from bit buffer: %w", err)
		}
		err = w.WriteBits(v, nbits)
		if err != nil {
			return fmt.Errorf("bit buffer error: write to bit writer: %w", err)
		}
	}
	return nil
}
