package pipeline

// BytesIn is the start of a Bytes Pipeline
type BytesIn func([]byte) chan []byte

// BytesStage is a stage in a Bytes Pipeline
type BytesStage func(chan []byte) chan []byte

// BytesOut is the end of a Bytes Pipeline
type BytesOut func(chan []byte) []byte

// Bytes is the interface for a Bytes Pipeline
type Bytes func([]byte) []byte

// BytesBuilder constructs a Bytes
type BytesBuilder struct {
	bytesIn  BytesIn
	bytes    []BytesStage
	bytesOut BytesOut
}

// NewBytesBuilder returns a pointer to a new BytesBuilder
func NewBytesBuilder(in BytesIn, out BytesOut) *BytesBuilder {
	return &BytesBuilder{bytesIn: in, bytesOut: out}
}

// AddStage adds a BytesStage to the BytesBuilder
func (b *BytesBuilder) AddStage(s BytesStage) *BytesBuilder {
	b.bytes = append(b.bytes, s)
	return b
}

// Build bytes Bytes
func (b *BytesBuilder) Build() Bytes {
	return func(in []byte) []byte {
		ch := b.bytesIn(in)
		for _, fn := range b.bytes {
			ch = fn(ch)
		}
		return b.bytesOut(ch)
	}
}

// BytesNoOp pipeline does nothing
func BytesNoOp(in []byte) []byte {
	return in
}
