package bit

type Bit bool

const (
	Zero Bit = false
	One  Bit = true
)

type ReadWriter interface {
	Reader
	Writer
}
