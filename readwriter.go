package bit

type BitReadWriter struct {
	*BitReader
	*BitWriter
}

func NewReadWriter(r *BitReader, w *BitWriter) *BitReadWriter {
	return &BitReadWriter{r, w}
}
