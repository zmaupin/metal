package pipeline

import "strings"
import "testing"

func TestBytesBuilder(t *testing.T) {
	pipeline := NewBytesBuilder(
		func(in []byte) chan []byte {
			ch := make(chan []byte)
			go func() { ch <- in }()
			return ch
		},
		func(in chan []byte) []byte {
			return <-in
		},
	).AddStage(func(in chan []byte) chan []byte {
		out := make(chan []byte)
		go func() {
			for byteArray := range in {
				s := string(byteArray)
				s = strings.Replace(s, "a", "o", -1)
				out <- []byte(s)
			}
		}()
		return out
	}).AddStage(func(in chan []byte) chan []byte {
		out := make(chan []byte)
		go func() {
			for byteArray := range in {
				s := string(byteArray)
				s = strings.Replace(s, "e", "u", -1)
				out <- []byte(s)
			}
		}()
		return out
	}).Build()
	in := []byte("apple")
	expected := "opplu"
	out := pipeline(in)
	if string(out) != expected {
		t.Errorf("expected out to be opplu, got %s", string(out))
	}
}
